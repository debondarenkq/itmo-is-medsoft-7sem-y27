import HISAPI
import SwiftUI

struct RecordView: View {
  let gateway: Gateway
  let worker: Staff
  let patient: Patient
  @Environment(\.dismiss) private var dismiss
  @State private var draft: RecordDraft?
  @State private var editing = false
  @State private var confirmingUnlock = false
  @State private var busy = false
  @State private var missing = false
  @State private var message: String?
  @State private var conflict = false
  @State private var sheet: EditorSheet?
  @State private var pendingAction: PendingAction?
  @State private var selectedPrescription: UUID?

  private enum EditorSheet: String, Identifiable {
    case diagnosis, prescription
    var id: Self { self }
  }
  private enum PendingAction { case close, reload, cancel }
  @State private var editedPrescription: Prescription?
  private var dirty: Bool { draft?.isDirty ?? false }

  var body: some View {
    VStack(spacing: 0) {
      HStack {
        VStack(alignment: .leading) {
          Text(patient.fullName).font(.title2)
          Text("Электронная медицинская карта").foregroundStyle(.secondary)
          Text(
            "\(patient.birthDate.calendarDate.formatted(date: .numeric, time: .omitted)) · \(patient.administrativeSex.label)"
          )
          .foregroundStyle(.secondary)
        }
        Spacer()
        Button("Обновить", systemImage: "arrow.clockwise") { request(.reload) }
        Button("Закрыть") { request(.close) }.keyboardShortcut(.cancelAction)
      }.padding()
      ErrorNotice(message: message)
      if conflict {
        HStack {
          Text("Черновик сохранён. Для загрузки актуальной карты отмените изменения.")
          Button("Обновить карту") { request(.reload) }
        }.padding()
      }
      if let draft {
        HStack {
          Label(
            editing ? "Редактирование" : "Только просмотр",
            systemImage: editing ? "lock.open" : "lock")
          Text("Версия \(draft.base.version)").foregroundStyle(.secondary)
          Spacer()
          if editing {
            Text("Автор: \(worker.fullName)")
          } else {
            Button("Разблокировать…") { confirmingUnlock = true }
          }
        }.padding()
        TabView {
          RecordDiagnosesView(
            gateway: gateway,
            draft: Binding(
              get: { self.draft ?? draft }, set: { self.draft = $0 }
            ), editing: editing, add: { sheet = .diagnosis }
          )
          .tabItem { Text("Диагнозы") }
          prescriptionsTab.tabItem { Text("Назначения") }
          RecordHistoryView(gateway: gateway, recordID: draft.base.id, version: draft.base.version)
            .tabItem { Text("История") }
        }
        .padding(.horizontal)
        if editing {
          HStack {
            Text("Несохранённых действий: \(draft.commands.count)").foregroundStyle(.secondary)
            Spacer()
            Button("Отменить изменения") { request(.cancel) }
            Button("Сохранить") { Task { await save() } }
              .keyboardShortcut("s", modifiers: .command)
              .disabled(!dirty)
          }.padding()
        }
      } else if missing {
        ContentUnavailableView {
          Label("Медицинская карта не создана", systemImage: "heart.text.clipboard")
        } actions: {
          Button("Создать ЭМК") { Task { await create() } }
        }
      } else {
        ContentUnavailableView(
          "Загрузка карты", systemImage: "heart.text.clipboard",
          description: Text("Если сервер недоступен, проверьте подключение и повторите запрос."))
      }
    }
    .frame(minWidth: 800, idealWidth: 960, minHeight: 580, idealHeight: 680)
    .disabled(busy)
    .overlay { if busy { ProgressView() } }
    .interactiveDismissDisabled()
    .task { await load() }
    .confirmationDialog(
      "Разблокировать карту?", isPresented: $confirmingUnlock, titleVisibility: .visible
    ) {
      Button("Разблокировать") { editing = true }
      Button("Отмена", role: .cancel) {}
    } message: {
      Text("Автор изменений: \(worker.fullName) · \(worker.position)")
    }
    .sheet(item: $sheet) { destination in
      switch destination {
      case .diagnosis:
        DirectoryPicker<Diagnosis>(
          title: "Добавить диагноз", label: { "\($0.code) · \($0.name)" },
          fetch: { search, offset in
            let page = try await gateway.diagnoses.listDiagnoses(
              limit: 50, offset: offset, search: search)
            return (page.items, page.total)
          }, select: { draft?.addDiagnosis($0) })
      case .prescription:
        PrescriptionEditor(prescription: editedPrescription) { text in
          draft?.setPrescription(editedPrescription, text: text)
        }
      }
    }
    .confirmationDialog(
      "Отменить несохранённые изменения?",
      isPresented: Binding(
        get: { pendingAction != nil }, set: { if !$0 { pendingAction = nil } }
      ), titleVisibility: .visible
    ) {
      Button("Отменить изменения", role: .destructive) {
        if let action = pendingAction { perform(action) }
        pendingAction = nil
      }
      Button("Продолжить редактирование", role: .cancel) { pendingAction = nil }
    }
  }

  private var prescriptionsTab: some View {
    VStack {
      Table(draft?.prescriptions ?? [], selection: $selectedPrescription) {
        TableColumn("Назначение", value: \.text)
        TableColumn("Статус") { Text($0.status == .active ? "Активно" : "Отменено") }
      }
      if editing {
        HStack {
          Button("Добавить…", systemImage: "plus") {
            editedPrescription = nil
            sheet = .prescription
          }
          Button("Изменить…") {
            editedPrescription = prescription
            sheet = .prescription
          }.disabled(prescription == nil)
          Button("Отменить назначение") {
            if let prescription { draft?.cancelPrescription(prescription) }
          }
          .disabled(prescription == nil || prescription.map { draft?.isEdited($0) == true } == true)
          .help("Изменённое назначение сначала нужно сохранить.")
          Button("Удалить") { draft?.prescriptions.removeAll { $0.id == selectedPrescription } }
            .disabled(prescription == nil)
          Spacer()
        }.padding()
      }
    }
  }

  private var prescription: Prescription? {
    draft?.prescriptions.first { $0.id == selectedPrescription }
  }

  private func request(_ action: PendingAction) {
    if dirty { pendingAction = action } else { perform(action) }
  }

  private func perform(_ action: PendingAction) {
    switch action {
    case .close: dismiss()
    case .reload: Task { await load() }
    case .cancel:
      if let draft { self.draft = RecordDraft(draft.base) }
      editing = false
      message = nil
      conflict = false
    }
  }

  private func replace(_ record: Record) {
    draft = RecordDraft(record)
    editing = false
    missing = false
    message = nil
    conflict = false
    selectedPrescription = nil
  }

  private func load() async {
    busy = true
    defer { busy = false }
    do { replace(try await gateway.patients.getPatientRecord(id: patient.id)) } catch {
      if case ErrorResponse.error(404, _, _, _) = error {
        missing = true
      } else {
        message = errorMessage(error)
      }
    }
  }

  private func create() async {
    busy = true
    defer { busy = false }
    do {
      replace(
        try await gateway.records.createRecord(
          id: patient.id, createRecordInput: CreateRecordInput(staffId: worker.id)))
    } catch { message = errorMessage(error) }
  }

  private func save() async {
    guard let draft else { return }
    busy = true
    defer { busy = false }
    do {
      let input = ChangesInput(
        commands: draft.commands, expectedVersion: draft.base.version, staffId: worker.id)
      replace(try await gateway.records.saveRecord(id: draft.base.id, changesInput: input))
    } catch {
      message = errorMessage(error)
      conflict = serverDetail(error)?.code == "VERSION_CONFLICT"
    }
  }
}

private struct PrescriptionEditor: View {
  let prescription: Prescription?
  let save: (String) -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var text = ""

  var body: some View {
    VStack(spacing: 0) {
      Form {
        TextField("Назначение", text: $text, axis: .vertical).lineLimit(3...8)
      }.formStyle(.grouped)
      FormActions(
        busy: false, saveTitle: "Добавить в черновик", cancel: { dismiss() },
        save: {
          save(text)
          dismiss()
        })
    }
    .frame(width: 520)
    .interactiveDismissDisabled()
    .onAppear { text = prescription?.text ?? "" }
  }
}
