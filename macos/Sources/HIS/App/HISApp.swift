import SwiftUI

@main
struct HISApp: App {
  @AppStorage(AppConfiguration.gatewayAddressKey)
  private var address = AppConfiguration.defaultGatewayAddress

  var body: some Scene {
    Window("Госпитальная информационная система", id: "main") {
      MainView(gateway: Gateway(address: address))
        .frame(minWidth: 850, minHeight: 560)
    }
    .defaultSize(width: 1100, height: 740)
    Settings { ConnectionSettings() }
  }
}
