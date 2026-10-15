import HISAPI
import SwiftUI

struct PatientsView: View {
  let gateway: Gateway
  let worker: Staff
  @State private var patients: [Patient] = []
  @State private var selectedPatientID: UUID?
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
    VStack(spacing: 0) {
      ErrorNotice(message: errorText)
      Table(patients, selection: $selectedPatientID) {
        TableColumn("ФИО", value: \.fullName)
        TableColumn("Дата рождения") {
          Text($0.birthDate.calendarDate, format: .dateTime.day().month().year())
        }
        TableColumn("Пол") { Text($0.administrativeSex.label) }
        TableColumn("ЭМК") { Text($0.hasRecord ? "Создана" : "Нет") }
      }
      .contextMenu(forSelectionType: UUID.self) { _ in
      } primaryAction: { ids in
        if let value = patients.first(where: { ids.contains($0.id) }) { viewedPatient = value }
      }
      PageControls(offset: $offset, total: total, count: patients.count)
    }
    .navigationTitle("Пациенты")
    .toolbar {
      Button("Медицинская карта", systemImage: "heart.text.clipboard") {
        recordPatient = selectedPatient
      }
      .disabled(selectedPatient == nil)
      Button("Просмотреть", systemImage: "eye") { viewedPatient = selectedPatient }.disabled(
        selectedPatient == nil)
      Button("Добавить", systemImage: "plus") { edit(nil) }
      Button("Изменить", systemImage: "square.and.pencil") { edit(selectedPatient) }.disabled(
        selectedPatient == nil)
      Button("Удалить", systemImage: "trash") { deleting = true }.disabled(selectedPatient == nil)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .disabled(isBusy)
    .overlay { if isBusy { ProgressView() } }
    .task(id: offset) { await load() }
    .sheet(item: $viewedPatient, onDismiss: { Task { await load() } }) { patient in
      PatientCardView(gateway: gateway, worker: worker, patient: patient)
    }
    .sheet(item: $editor) { selection in
      PatientEditorView(gateway: gateway, patient: selection.value) { Task { await load() } }
    }
    .sheet(item: $recordPatient, onDismiss: { Task { await load() } }) { patient in
      RecordView(gateway: gateway, worker: worker, patient: patient)
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
      let page = try await gateway.patients.listPatients(limit: 50, offset: offset)
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
