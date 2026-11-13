import HISAPI
import SwiftUI

struct LISHistoryView: View {
  let gateway: Gateway
  @State private var history: [LISHistory] = []
  @State private var search = ""
  @State private var selectedID: UUID?
  @State private var opened: LISHistory?
  @State private var errorText: String?

  var body: some View {
    VStack(spacing: 0) {
      ErrorNotice(message: errorText)
      Table(history, selection: $selectedID) {
        TableColumn("Дата проведения") { Text($0.result.effectiveDateTime, format: .dateTime.day().month().year().hour().minute()) }
        TableColumn("Пациент") { Text($0.patient.fullName) }
        TableColumn("Анализ") { Text($0.request.testName) }
        TableColumn("Результат") { Text($0.result.valueString) }
        TableColumn("Доставка") { Text($0.delivered ? "В HIS" : "Не подтверждена").foregroundStyle($0.delivered ? Color.green : Color.orange) }
      }
      .contextMenu(forSelectionType: UUID.self) { ids in
        Button("Подробности") { open(ids.first) }
      } primaryAction: { ids in open(ids.first) }
      HStack {
        Text("Результатов: \(history.count)").foregroundStyle(.secondary)
        Spacer()
        Button("Подробности") { open(selectedID) }.disabled(selectedID == nil)
      }.padding()
    }
    .navigationTitle("История анализов")
    .searchable(text: $search, prompt: "Пациент, анализ или результат")
    .toolbar { Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } } }
    .task(id: search) { await load() }
    .sheet(item: $opened) { item in LISResultView(gateway: gateway, requestID: item.request.id, snapshot: item) }
  }

  private func open(_ id: UUID?) { opened = history.first { $0.id == id } }

  private func load() async {
    do { history = try await gateway.lis.listLISHistory(search: search); errorText = nil }
    catch { errorText = errorMessage(error) }
  }
}
