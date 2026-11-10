import SwiftUI

struct ConnectionSettings: View {
  @Environment(\.openWindow) private var openWindow
  @AppStorage(AppConfiguration.gatewayAddressKey)
  private var savedAddress = AppConfiguration.defaultGatewayAddress
  @State private var address = ""
  @State private var errorText: String?

  var body: some View {
    Form {
      TextField("Адрес сервера", text: $address)
      Text("Например, \(AppConfiguration.defaultGatewayAddress)").foregroundStyle(.secondary)
      ErrorNotice(message: errorText)
      Button("Применить", action: applySettings)
      Button("Сертификаты сервера…") { openWindow(id: "server-certificates") }
    }
    .formStyle(.grouped)
    .frame(width: 480)
    .onAppear { address = savedAddress }
  }

  private func applySettings() {
    guard !NSApp.windows.contains(where: { $0.attachedSheet != nil }) else {
      errorText = "Закройте открытый редактор перед сменой сервера."
      return
    }
    let value = address.trimmingCharacters(in: .whitespacesAndNewlines)
    guard let url = URLComponents(string: value),
      url.scheme == "https",
      let host = url.host, !host.isEmpty,
      url.user == nil, url.password == nil, url.query == nil, url.fragment == nil
    else {
      errorText = "Укажите адрес вида \(AppConfiguration.defaultGatewayAddress)."
      return
    }
    savedAddress = value.hasSuffix("/") ? String(value.dropLast()) : value
    errorText = nil
  }
}
