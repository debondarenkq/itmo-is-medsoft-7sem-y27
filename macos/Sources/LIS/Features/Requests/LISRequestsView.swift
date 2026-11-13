import HISAPI
import SwiftUI

private enum RequestFilter: String, CaseIterable { case new = "Новые", processed = "Обработанные", all = "Все" }

struct LISRequestsView: View {
  let gateway: Gateway
  let configured: Bool
  @State private var requests: [FhirFHIRServiceRequest] = []
  @State private var processed: Set<UUID> = []
  @State private var filter = RequestFilter.new
  @State private var search = ""
  @State private var selectedID: UUID?
  @State private var opened: FhirFHIRServiceRequest?
  @State private var errorText: String?
  @State private var busy = false

  private var visible: [FhirFHIRServiceRequest] {
    requests.filter { request in
      let done = processed.contains(request.id) || request.status == .completed
      let matchesStatus = filter == .all || (filter == .processed ? done : !done)
      let text = "\(request.subject.display ?? "") \(request.testName) \(request.loincCode) \(request.id)"
      return matchesStatus && (search.isEmpty || text.localizedCaseInsensitiveContains(search))
    }
  }

  var body: some View {
    VStack(alignment: .leading, spacing: 0) {
      Picker("Назначения", selection: $filter) {
        ForEach(RequestFilter.allCases, id: \.self) { Text($0.rawValue).tag($0) }
      }.pickerStyle(.segmented).padding()
      ErrorNotice(message: errorText)
      if !configured {
        ContentUnavailableView("Настройте подключение к HIS", systemImage: "network", description: Text("Откройте раздел «Доступ к HIS»."))
      } else {
        Table(visible, selection: $selectedID) {
          TableColumn("Назначен") { Text($0.authoredOn, format: .dateTime.day().month().year().hour().minute()) }
          TableColumn("Пациент") { Text($0.subject.display ?? $0.subject.reference) }
          TableColumn("Анализ") { Text($0.testName) }
          TableColumn("LOINC", value: \.loincCode)
          TableColumn("Статус") { request in
            if processed.contains(request.id) {
              Label("Уже обработан в LIS", systemImage: "checkmark.lock").foregroundStyle(.green)
            } else { Text(request.status == .completed ? "Выполнен в HIS" : "Новый") }
          }
        }
        .contextMenu(forSelectionType: UUID.self) { ids in
          Button("Открыть") { open(ids.first) }
        } primaryAction: { ids in open(ids.first) }
        HStack {
          Text("Назначений: \(visible.count)").foregroundStyle(.secondary)
          Spacer()
          Button("Открыть анализ", systemImage: "testtube.2") { open(selectedID) }
            .buttonStyle(.borderedProminent).disabled(selectedID == nil)
        }.padding()
      }
    }
    .navigationTitle("Назначения из HIS")
    .searchable(text: $search, prompt: "Пациент, анализ или LOINC")
    .toolbar { Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }.disabled(!configured || busy) }
    .task { if configured { await load() } }
    .sheet(item: $opened, onDismiss: { Task { await load() } }) { request in
      LISResultView(gateway: gateway, requestID: request.id)
    }
  }

  private func open(_ id: UUID?) { opened = requests.first { $0.id == id } }

  private func load() async {
    busy = true
    defer { busy = false }
    do {
      processed = Set(try await gateway.lis.listLISHistory().map { $0.request.id })
      requests = try await gateway.lis.listLISRequests().resources(FhirFHIRServiceRequest.self)
      errorText = nil
    } catch { errorText = errorMessage(error) }
  }
}
