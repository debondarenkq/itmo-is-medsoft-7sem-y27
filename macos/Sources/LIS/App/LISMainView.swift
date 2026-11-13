import HISAPI
import SwiftUI

private enum LISSection: Hashable { case requests, history, connection }

struct LISMainView: View {
  let gateway: Gateway
  @State private var section: LISSection? = .requests
  @State private var status: LISConnectionStatus?
  @State private var errorText: String?

  var body: some View {
    NavigationSplitView {
      List(selection: $section) {
        Label("Назначения", systemImage: "list.clipboard").tag(LISSection.requests)
        Label("История анализов", systemImage: "clock.arrow.circlepath").tag(LISSection.history)
        Label("Доступ к HIS", systemImage: "network").tag(LISSection.connection)
      }
      .navigationTitle("Лаборатория · LIS")
      .safeAreaInset(edge: .bottom) {
        VStack(alignment: .leading, spacing: 8) {
          Divider()
          Label("Доступ без отдельного входа", systemImage: "cross.case")
          Text("Единицы измерения не используются.").font(.caption).foregroundStyle(.secondary)
        }.padding()
      }
      .navigationSplitViewColumnWidth(min: 240, ideal: 260, max: 320)
    } detail: {
      if let status {
        switch section {
        case .history: LISHistoryView(gateway: gateway)
        case .connection: LISConnectionView(gateway: gateway, initial: status) { self.status = $0 }
        default: LISRequestsView(gateway: gateway, configured: status.configured)
        }
      } else {
        VStack(spacing: 16) {
          ErrorNotice(message: errorText)
          if errorText == nil { ProgressView("Подключение к серверу LIS…") }
          else { Button("Повторить") { Task { await load() } } }
        }.padding()
      }
    }
    .task { await load() }
  }

  private func load() async {
    do {
      status = try await gateway.lis.getLISConnection()
      if status?.configured != true { section = .connection }
      errorText = nil
    } catch { errorText = errorMessage(error) }
  }
}
