import HISAPI
import SwiftUI

struct PatientsView: View {
  let gateway: Gateway
  let worker: Staff
  var recordsOnly = false
  var editingChanged: (Bool) -> Void = { _ in }
  @State private var patients: [Patient] = []
  @State private var selectedPatientID: UUID?
  @State private var search = ""
  @State private var offset = 0
  @State private var total = 0
  @State private var isBusy = false
  @State private var errorText: String?
  @State private var editor: PatientEditorSelection?
  @State private var deleting = false
  @State private var recordPatient: Patient?
  @State private var viewedPatient: Patient?
  private var selectedPatient: Patient? { patients.first { $0.id == selectedPatientID } }

  var body: some View {
    Group {
      if let patient = recordPatient {
        RecordView(gateway: gateway, worker: worker, patient: patient,
          close: { recordPatient = nil; Task { await load() } }, editingChanged: editingChanged)
      } else { patientsList }
    }
  }

  private var patientsList: some View {
    VStack(spacing: 0) {
      HStack {
        Text(recordsOnly ? "Медицинские карты пациентов" : "Пациенты медицинской организации").foregroundStyle(.secondary)
        Spacer()
        VStack(alignment: .trailing) {
          Text(worker.fullName).font(.headline)
          Text("\(worker.position) · \(worker.role.label)").font(.caption).foregroundStyle(.secondary)
        }
      }.padding()
      Divider()
      ErrorNotice(message: errorText)
      Table(patients, selection: $selectedPatientID) {
        TableColumn("Номер ЭМК") { Text($0.recordLabel).font(.system(.caption, design: .monospaced)).lineLimit(1).truncationMode(.middle).help($0.recordLabel) }
        TableColumn("ФИО", value: \.fullName)
        TableColumn("Дата рождения") {
          Text($0.birthDate.calendarDate, format: .dateTime.day().month().year())
        }
        TableColumn("Пол") { Text($0.administrativeSex.label) }
        TableColumn("Действия") { patient in
          Button("Открыть", systemImage: "arrow.up.right.square") { recordPatient = patient }
            .buttonStyle(.borderless)
        }.width(100)
      }
      .contextMenu(forSelectionType: UUID.self) { _ in
      } primaryAction: { ids in
        if let value = patients.first(where: { ids.contains($0.id) }) { recordPatient = value }
      }
      PageControls(offset: $offset, total: total, count: patients.count)
    }
    .navigationTitle(recordsOnly ? "ЭМК" : "Пациенты")
    .searchable(text: $search, prompt: "ФИО или номер ЭМК")
    .toolbar {
      Button("Медицинская карта", systemImage: "heart.text.clipboard") {
        recordPatient = selectedPatient
      }
      .disabled(selectedPatient == nil)
      Button("Просмотреть", systemImage: "eye") { viewedPatient = selectedPatient }.disabled(
        selectedPatient == nil)
      Button("Создать пациента", systemImage: "plus") { edit(nil) }.disabled(!worker.role.canCreate)
      Button("Изменить", systemImage: "square.and.pencil") { edit(selectedPatient) }.disabled(
        selectedPatient == nil || !worker.role.canCreate)
      Button("Удалить", systemImage: "trash") { deleting = true }.disabled(selectedPatient == nil || !worker.role.canCreate)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .disabled(isBusy)
    .overlay { if isBusy { ProgressView() } }
    .task(id: "\(offset)|\(search)|\(recordsOnly)") {
      do { try await Task.sleep(for: .milliseconds(250)); await load() } catch {}
    }
    .onChange(of: search) { offset = 0; selectedPatientID = nil }
    .sheet(item: $viewedPatient, onDismiss: { Task { await load() } }) { patient in
      PatientCardView(gateway: gateway, worker: worker, patient: patient)
    }
    .sheet(item: $editor) { selection in
      PatientEditorView(gateway: gateway, patient: selection.value) { Task { await load() } }
    }
    .confirmationDialog("Удалить пациента?", isPresented: $deleting, titleVisibility: .visible) {
      Button("Удалить", role: .destructive) { Task { await delete() } }
    } message: {
      Text(selectedPatient?.fullName ?? "")
    }
  }

  private func edit(_ patient: Patient?) {
    editor = PatientEditorSelection(value: patient)
  }

  private func load() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let page = try await gateway.patients.listPatients(limit: 50, offset: offset, search: search, hasRecord: recordsOnly ? true : nil)
      try Task.checkCancellation()
      patients = page.items
      total = page.total
      errorText = nil
    } catch {
      if !Task.isCancelled {
        errorText = errorMessage(error)
      }
    }
  }

  private func delete() async {
    guard let selectedPatient else { return }
    isBusy = true
    do {
      try await gateway.patients.deletePatient(id: selectedPatient.id)
      await load()
    } catch {
      errorText = errorMessage(error)
    }
    isBusy = false
  }
}

private struct PatientEditorSelection: Identifiable {
  let id = UUID()
  let value: Patient?
}
