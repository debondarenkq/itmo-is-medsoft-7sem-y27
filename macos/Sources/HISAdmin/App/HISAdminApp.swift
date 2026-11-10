import SwiftUI

@main
struct HISAdminApp: App {
  @AppStorage(AppConfiguration.gatewayAddressKey)
  private var address = AppConfiguration.defaultGatewayAddress

  var body: some Scene {
    Window("HIS — Администрирование", id: "main") {
      SessionView(address: address).id(address)
        .frame(minWidth: 850, minHeight: 560)
    }
    .defaultSize(width: 1100, height: 740)
    .commands { ServerCertificateCommands() }
    Window("Сертификаты сервера", id: "server-certificates") {
      ServerCertificatesView()
    }
    Settings { ConnectionSettings() }
  }
}
