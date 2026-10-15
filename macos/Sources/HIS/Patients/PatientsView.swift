import HISAPI
import SwiftUI

struct PatientsView: View {
  let gateway: Gateway
  let worker: Staff
  @State private var items: [Patient] = []
  @State private var selection: UUID?
  @State private var search = ""
  @State private var offset = 0
  @State private var total = 0
  @State private var busy = false
  @State private var message: String?
  @State private var editor: PatientEditorSelection?
  @State private var deleting = false
  @State private var recordPatient: Patient?
  @State private var viewedPatient: Patient?
  private var selected: Patient? { items.first { $0.id == selection } }

  var body: some View {
    VStack(spacing: 0) {
      ErrorNotice(message: message)
      Table(items, selection: $selection) {
        TableColumn("ФИО", value: \.fullName)
        TableColumn("Дата рождения") {
          Text($0.birthDate.calendarDate, format: .dateTime.day().month().year())
        }
        TableColumn("Пол") { Text($0.administrativeSex.label) }
        TableColumn("ЭМК") { Text($0.hasRecord ? "Создана" : "Нет") }
      }
      .contextMenu(forSelectionType: UUID.self) { _ in
      } primaryAction: { ids in
        if let value = items.first(where: { ids.contains($0.id) }) { viewedPatient = value }
      }
      PageControls(offset: $offset, total: total, count: items.count)
    }
    .navigationTitle("Пациенты")
    .searchable(text: $search, prompt: "Фамилия или имя")
    .toolbar {
      Button("Медицинская карта", systemImage: "heart.text.clipboard") { recordPatient = selected }
        .disabled(selected == nil)
      Button("Просмотреть", systemImage: "eye") { viewedPatient = selected }.disabled(
        selected == nil)
      Button("Добавить", systemImage: "plus") { edit(nil) }
      Button("Изменить", systemImage: "square.and.pencil") { edit(selected) }.disabled(
        selected == nil)
      Button("Удалить", systemImage: "trash") { deleting = true }.disabled(selected == nil)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .disabled(busy)
    .overlay { if busy { ProgressView() } }
    .task(id: "\(search)|\(offset)") { await load() }
    .onChange(of: search) { offset = 0 }
    .sheet(item: $viewedPatient, onDismiss: { Task { await load() } }) { patient in
      PatientCard(gateway: gateway, worker: worker, patient: patient)
    }
    .sheet(item: $editor) { selection in
      PatientEditor(gateway: gateway, patient: selection.value) { Task { await load() } }
    }
    .sheet(item: $recordPatient, onDismiss: { Task { await load() } }) { patient in
      RecordView(gateway: gateway, worker: worker, patient: patient)
    }
    .confirmationDialog("Удалить пациента?", isPresented: $deleting, titleVisibility: .visible) {
      Button("Удалить", role: .destructive) { Task { await delete() } }
    } message: {
      Text(selected?.fullName ?? "")
    }
  }

  private func edit(_ patient: Patient?) {
    editor = PatientEditorSelection(value: patient)
  }

  private func load() async {
    busy = true
    defer { busy = false }
    do {
      let page = try await gateway.patients.listPatients(limit: 50, offset: offset, search: search)
      try Task.checkCancellation()
      items = page.items
      total = page.total
      message = nil
    } catch { if !Task.isCancelled { message = errorMessage(error) } }
  }

  private func delete() async {
    guard let selected else { return }
    busy = true
    do {
      try await gateway.patients.deletePatient(id: selected.id)
      await load()
    } catch { message = errorMessage(error) }
    busy = false
  }
}

private struct PatientEditorSelection: Identifiable {
  let id = UUID()
  let value: Patient?
}

struct PatientEditor: View {
  let gateway: Gateway
  let patient: Patient?
  let saved: () -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var firstName = ""
  @State private var lastName = ""
  @State private var middleName = ""
  @State private var birthDate = Date()
  @State private var sex: AdministrativeSex = .unknown
  @State private var comment = ""
  @State private var busy = false
  @State private var message: String?

  init(gateway: Gateway, patient: Patient?, saved: @escaping () -> Void) {
    self.gateway = gateway
    self.patient = patient
    self.saved = saved
    _firstName = State(initialValue: patient?.firstName ?? "")
    _lastName = State(initialValue: patient?.lastName ?? "")
    _middleName = State(initialValue: patient?.middleName ?? "")
    _birthDate = State(initialValue: patient?.birthDate.calendarDate ?? Date())
    _sex = State(initialValue: patient?.administrativeSex ?? .unknown)
    _comment = State(initialValue: patient?.comment ?? "")
  }

  var body: some View {
    VStack(spacing: 0) {
      Form {
        Section(patient == nil ? "Новый пациент" : "Пациент") {
          TextField("Фамилия", text: $lastName)
          TextField("Имя", text: $firstName)
          TextField("Отчество", text: $middleName)
          DatePicker("Дата рождения", selection: $birthDate, displayedComponents: .date)
          Picker("Пол", selection: $sex) {
            ForEach(AdministrativeSex.allCases, id: \.self) { Text($0.label).tag($0) }
          }
          TextField("Комментарий", text: $comment, axis: .vertical)
        }
      }.formStyle(.grouped).disabled(busy)
      ErrorNotice(message: message)
      FormActions(busy: busy, cancel: { dismiss() }, save: { Task { await save() } })
    }
    .frame(width: 520)
    .interactiveDismissDisabled()
  }

  private func save() async {
    busy = true
    defer { busy = false }
    do {
      let input = PatientInput(
        administrativeSex: sex,
        birthDate: OpenAPIDateWithoutTime(wrappedDate: birthDate, timezone: .current),
        comment: comment, firstName: firstName, lastName: lastName,
        middleName: middleName.isEmpty ? nil : middleName)
      if let patient {
        _ = try await gateway.patients.updatePatient(id: patient.id, patientInput: input)
      } else {
        _ = try await gateway.patients.createPatient(patientInput: input)
      }
      saved()
      dismiss()
    } catch { message = errorMessage(error) }
  }
}
