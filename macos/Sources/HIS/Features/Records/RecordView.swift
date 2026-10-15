import HISAPI
import SwiftUI

struct RecordView: View {
  let gateway: Gateway
  let worker: Staff
  let patient: Patient
  @Environment(\.dismiss) private var dismiss
  @State private var draft: RecordDraft?
  @State private var isEditing = false
  @State private var confirmingUnlock = false
  @State private var isBusy = false
  @State private var recordIsMissing = false
  @State private var errorText: String?
  @State private var hasVersionConflict = false
  @State private var sheet: EditorSheet?
  @State private var pendingAction: PendingAction?
  @State private var showingDiscardConfirmation = false
  @State private var selectedPrescriptionID: UUID?

  private enum EditorSheet: Identifiable {
    case diagnosis
    case prescription(Prescription?)

    var id: String {
      switch self {
      case .diagnosis: return "diagnosis"
      case .prescription(let prescription): return prescription?.id.uuidString ?? "new-prescription"
      }
    }
  }
  private enum PendingAction { case close, reload, cancel }
  private var hasUnsavedChanges: Bool { draft?.isDirty ?? false }

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
        Button("Обновить", systemImage: "arrow.clockwise") { requestAction(.reload) }
        Button("Закрыть") { requestAction(.close) }.keyboardShortcut(.cancelAction)
      }.padding()
      ErrorNotice(message: errorText)
      if hasVersionConflict {
        HStack {
          Text("Черновик сохранён. Для загрузки актуальной карты отмените изменения.")
          Button("Обновить карту") { requestAction(.reload) }
        }.padding()
      }
      if let draft {
        HStack {
          Label(
            isEditing ? "Редактирование" : "Только просмотр",
            systemImage: isEditing ? "lock.open" : "lock")
          Text("Версия \(draft.base.version)").foregroundStyle(.secondary)
          Spacer()
          if isEditing {
            Text("Автор: \(worker.fullName)")
          } else {
            Button("Разблокировать…") { confirmingUnlock = true }
          }
        }.padding()
        TabView {
          RecordDiagnosesView(
            gateway: gateway,
            draft: draft,
            isEditing: isEditing,
            addDiagnosis: { sheet = .diagnosis },
            removeDiagnosis: { diagnosisID in self.draft?.removeDiagnosis(id: diagnosisID) },
            restoreDiagnosis: { diagnosis in self.draft?.restoreDiagnosis(diagnosis) }
          )
          .tabItem { Text("Диагнозы") }
          prescriptionsTab.tabItem { Text("Назначения") }
          RecordHistoryView(gateway: gateway, recordID: draft.base.id, version: draft.base.version)
            .tabItem { Text("История") }
        }
        .padding(.horizontal)
        if isEditing {
          HStack {
            Text("Несохранённых действий: \(draft.commands.count)").foregroundStyle(.secondary)
            Spacer()
            Button("Отменить изменения") { requestAction(.cancel) }
            Button("Сохранить") { Task { await saveChanges() } }
              .keyboardShortcut("s", modifiers: .command)
              .disabled(!hasUnsavedChanges)
          }.padding()
        }
      } else if recordIsMissing {
        ContentUnavailableView {
          Label("Медицинская карта не создана", systemImage: "heart.text.clipboard")
        } actions: {
          Button("Создать ЭМК") { Task { await createRecord() } }
        }
      } else {
        ContentUnavailableView(
          "Загрузка карты", systemImage: "heart.text.clipboard",
          description: Text("Если сервер недоступен, проверьте подключение и повторите запрос."))
      }
    }
    .frame(minWidth: 800, idealWidth: 960, minHeight: 580, idealHeight: 680)
    .disabled(isBusy)
    .overlay { if isBusy { ProgressView() } }
    .interactiveDismissDisabled()
    .task { await loadRecord() }
    .confirmationDialog(
      "Разблокировать карту?", isPresented: $confirmingUnlock, titleVisibility: .visible
    ) {
      Button("Разблокировать") { isEditing = true }
      Button("Отмена", role: .cancel) {}
    } message: {
      Text("Автор изменений: \(worker.fullName) · \(worker.position)")
    }
    .sheet(item: $sheet) { destination in
      switch destination {
      case .diagnosis:
        DirectoryPicker<Diagnosis>(
          title: "Добавить диагноз", label: { "\($0.code) · \($0.name)" },
          fetch: { offset in
            let page = try await gateway.diagnoses.listDiagnoses(
              limit: 50, offset: offset)
            return (page.items, page.total)
          }, select: { draft?.addDiagnosis($0) })
      case .prescription(let prescription):
        PrescriptionEditorView(prescription: prescription) { text in
          draft?.setPrescription(prescription, text: text)
        }
      }
    }
    .confirmationDialog(
      "Отменить несохранённые изменения?",
      isPresented: $showingDiscardConfirmation, titleVisibility: .visible
    ) {
      Button("Отменить изменения", role: .destructive) {
        if let action = pendingAction { performAction(action) }
        pendingAction = nil
      }
      Button("Продолжить редактирование", role: .cancel) { pendingAction = nil }
    }
  }

  private var prescriptionsTab: some View {
    VStack {
      Table(draft?.prescriptions ?? [], selection: $selectedPrescriptionID) {
        TableColumn("Назначение", value: \.text)
        TableColumn("Статус") { Text($0.status == .active ? "Активно" : "Отменено") }
      }
      if isEditing {
        HStack {
          Button("Добавить…", systemImage: "plus") {
            sheet = .prescription(nil)
          }
          Button("Изменить…") {
            sheet = .prescription(selectedPrescription)
          }.disabled(selectedPrescription == nil)
          Button("Отменить назначение") {
            if let selectedPrescription { draft?.cancelPrescription(selectedPrescription) }
          }
          .disabled(!canCancelPrescription)
          .help("Изменённое назначение сначала нужно сохранить.")
          Button("Удалить") { draft?.prescriptions.removeAll { $0.id == selectedPrescriptionID } }
            .disabled(selectedPrescription == nil)
          Spacer()
        }.padding()
      }
    }
  }

  private var selectedPrescription: Prescription? {
    draft?.prescriptions.first { $0.id == selectedPrescriptionID }
  }

  private var canCancelPrescription: Bool {
    guard let selectedPrescription, let draft else { return false }
    return !draft.isEdited(selectedPrescription)
  }

  private func requestAction(_ action: PendingAction) {
    if hasUnsavedChanges {
      pendingAction = action
      showingDiscardConfirmation = true
    } else {
      performAction(action)
    }
  }

  private func performAction(_ action: PendingAction) {
    switch action {
    case .close: dismiss()
    case .reload: Task { await loadRecord() }
    case .cancel:
      if let draft { self.draft = RecordDraft(draft.base) }
      isEditing = false
      errorText = nil
      hasVersionConflict = false
    }
  }

  private func replaceDraft(_ record: Record) {
    draft = RecordDraft(record)
    isEditing = false
    recordIsMissing = false
    errorText = nil
    hasVersionConflict = false
    selectedPrescriptionID = nil
  }

  private func loadRecord() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let record = try await gateway.patients.getPatientRecord(id: patient.id)
      replaceDraft(record)
    } catch {
      if case ErrorResponse.error(404, _, _, _) = error {
        recordIsMissing = true
      } else {
        errorText = errorMessage(error)
      }
    }
  }

  private func createRecord() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let input = CreateRecordInput(staffId: worker.id)
      let record = try await gateway.records.createRecord(id: patient.id, createRecordInput: input)
      replaceDraft(record)
    } catch {
      errorText = errorMessage(error)
    }
  }

  private func saveChanges() async {
    guard let draft else { return }
    isBusy = true
    defer { isBusy = false }
    do {
      let input = ChangesInput(
        commands: draft.commands, expectedVersion: draft.base.version, staffId: worker.id)
      let record = try await gateway.records.saveRecord(id: draft.base.id, changesInput: input)
      replaceDraft(record)
    } catch {
      errorText = errorMessage(error)
      hasVersionConflict = serverDetail(error)?.code == "VERSION_CONFLICT"
    }
  }
}
