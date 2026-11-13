import HISAPI
import SwiftUI

struct RecordLabTestsView: View {
  let gateway: Gateway
  let patient: Patient
  let canOrder: Bool
  @State private var requests: [FhirFHIRServiceRequest] = []
  @State private var results: [FhirFHIRObservation] = []
  @State private var errorText: String?
  @State private var ordering = false
  @State private var isBusy = false

  var body: some View {
    VStack(alignment: .leading) {
      ErrorNotice(message: errorText)
      Table(requests) {
        TableColumn("Назначен") { Text($0.authoredOn, format: .dateTime.day().month().year().hour().minute()) }
        TableColumn("Анализ") { Text("\($0.testName) · \($0.loincCode)") }
        TableColumn("Назначил") { Text($0.requester.display ?? "—") }
        TableColumn("Статус") { Text($0.status == .completed ? "Выполнен" : "Назначен") }
        TableColumn("Результат") { request in Text(result(for: request)?.valueString ?? "—").textSelection(.enabled) }
        TableColumn("Дата результата") { request in
          if let result = result(for: request) { Text(result.effectiveDateTime, format: .dateTime.day().month().year().hour().minute()) }
          else { Text("—") }
        }
      }
      HStack {
        Button("Назначить анализы…", systemImage: "plus") { ordering = true }.disabled(!canOrder)
        Button("Обновить результаты", systemImage: "arrow.clockwise") { Task { await load() } }
        Spacer()
        Text("Результат — строка; единицы измерения не используются.").font(.caption).foregroundStyle(.secondary)
      }.padding(.vertical, 8)
    }
    .disabled(isBusy)
    .task { await load() }
    .sheet(isPresented: $ordering, onDismiss: { Task { await load() } }) {
      LabOrderView(gateway: gateway, patient: patient)
    }
  }

  private func result(for request: FhirFHIRServiceRequest) -> FhirFHIRObservation? {
    results.first { $0.basedOn.contains { $0.reference == "ServiceRequest/\(request.id.uuidString.lowercased())" } }
  }

  private func load() async {
    isBusy = true
    defer { isBusy = false }
    do {
      requests = try await gateway.fhir.searchFHIRServiceRequest(patient: patient.id.uuidString.lowercased()).resources(FhirFHIRServiceRequest.self)
      results = try await gateway.fhir.searchFHIRObservation(patient: patient.id.uuidString.lowercased()).resources(FhirFHIRObservation.self)
      errorText = nil
    } catch { errorText = errorMessage(error) }
  }
}

private struct LabOrderView: View {
  let gateway: Gateway
  let patient: Patient
  @Environment(\.dismiss) private var dismiss
  @State private var tests: [LabTest] = []
  @State private var selected: Set<UUID> = []
  @State private var comment = ""
  @State private var errorText: String?
  @State private var isBusy = false

  var body: some View {
    VStack(alignment: .leading, spacing: 16) {
      Text("Назначить анализы").font(.title2)
      Text(patient.fullName).foregroundStyle(.secondary)
      ErrorNotice(message: errorText)
      List(tests) { test in
        Toggle("\(test.name) · LOINC \(test.code)", isOn: Binding(
          get: { selected.contains(test.id) },
          set: { if $0 { selected.insert(test.id) } else { selected.remove(test.id) } }
        ))
      }.frame(minHeight: 240)
      TextField("Комментарий к назначениям", text: $comment, axis: .vertical).lineLimit(3...4)
      Text("Каждый выбранный анализ будет отдельным назначением. Сохранение выполняется сразу.").font(.caption).foregroundStyle(.secondary)
      HStack {
        Button("Отмена") { dismiss() }.keyboardShortcut(.cancelAction)
        Spacer()
        Button("Назначить (\(selected.count))") { Task { await save() } }
          .buttonStyle(.borderedProminent).disabled(selected.isEmpty || selected.count > 20)
      }
    }.padding(24).frame(width: 600).disabled(isBusy)
    .task {
      do { tests = try await gateway.laboratory.listLabTests() }
      catch { errorText = errorMessage(error) }
    }
  }

  private func save() async {
    isBusy = true
    defer { isBusy = false }
    do {
      _ = try await gateway.laboratory.createLabOrders(labOrderInput: LabOrderInput(comment: comment, patientId: patient.id, testIds: Array(selected)))
      dismiss()
    } catch { errorText = errorMessage(error) }
  }
}
