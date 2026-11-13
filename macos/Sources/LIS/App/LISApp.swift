import SwiftUI

@main
struct LISApp: App {
  @AppStorage(AppConfiguration.gatewayAddressKey) private var address = AppConfiguration.defaultGatewayAddress

  var body: some Scene {
    Window("LIS — Лабораторная информационная система", id: "main") {
      LISMainView(gateway: Gateway(address: address)).id(address)
        .frame(minWidth: 980, minHeight: 640)
        .onAppear { configureApplicationIcon() }
    }
    .defaultSize(width: 1180, height: 760)
    .commands { ServerCertificateCommands() }
    Window("Сертификаты сервера", id: "server-certificates") { ServerCertificatesView() }
    Settings {
      Form { TextField("Адрес сервера LIS", text: $address) }.padding(24).frame(width: 460)
    }
  }
}
