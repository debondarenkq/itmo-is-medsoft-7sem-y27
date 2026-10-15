import SwiftUI

struct ConnectionSettings: View {
  @AppStorage("gatewayAddress") private var savedAddress = "http://localhost:8080"
  @State private var address = ""
  @State private var errorText: String?

  var body: some View {
    Form {
      TextField("Адрес сервера", text: $address)
      Text("Например, http://localhost:8080").foregroundStyle(.secondary)
      ErrorNotice(message: errorText)
      Button("Применить", action: applySettings)
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
      ["http", "https"].contains(url.scheme),
      let host = url.host, !host.isEmpty,
      url.user == nil, url.password == nil, url.query == nil, url.fragment == nil
    else {
      errorText = "Укажите адрес вида http://localhost:8080."
      return
    }
    savedAddress = value.hasSuffix("/") ? String(value.dropLast()) : value
    errorText = nil
  }
}
