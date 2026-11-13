import HISAPI
import SwiftUI

struct LISConnectionView: View {
  let gateway: Gateway
  let initial: LISConnectionStatus
  let changed: (LISConnectionStatus) -> Void
  @State private var apiKey = ""
  @State private var status: LISConnectionStatus?
  @State private var errorText: String?
  @State private var busy = false

  var body: some View {
    VStack(alignment: .leading, spacing: 22) {
      Label("Авторизация LIS в HIS", systemImage: "key.horizontal").font(.largeTitle)
      Text("Создайте сервисный аккаунт с ролью LIS в HISAdmin и введите его API-ключ.")
        .foregroundStyle(.secondary)
      Form {
        SecureField("API-ключ", text: $apiKey)
      }
      ErrorNotice(message: errorText)
      if let status, errorText == nil {
        Label(status.message, systemImage: status.configured ? "checkmark.circle.fill" : "info.circle")
          .foregroundStyle(status.configured ? Color.green : Color.secondary)
      }
      HStack {
        Spacer()
        Button("Авторизовать") { Task { await configure() } }
          .buttonStyle(.borderedProminent).disabled(apiKey.isEmpty)
      }
      Text("API-ключ не сохраняется. Сервер LIS хранит выданные токены. Ошибки связи показываются без автоматических повторов.")
        .font(.caption).foregroundStyle(.secondary)
      Spacer()
    }.padding(32).frame(maxWidth: 760).disabled(busy)
    .onAppear { status = initial }
  }

  private func configure() async {
    busy = true
    defer { busy = false }
    do {
      let updated = try await gateway.lis.configureLISConnection(lISConnectionInput: LISConnectionInput(apiKey: apiKey))
      status = updated; changed(updated); apiKey = ""; errorText = nil
    } catch { errorText = errorMessage(error) }
  }

}
