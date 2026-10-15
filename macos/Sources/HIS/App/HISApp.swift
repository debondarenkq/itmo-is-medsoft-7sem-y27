import SwiftUI

@main
struct HISApp: App {
  @AppStorage("gatewayAddress") private var address = "http://localhost:8080"

  var body: some Scene {
    Window("Госпитальная информационная система", id: "main") {
      MainView(gateway: Gateway(address: address))
        .frame(minWidth: 850, minHeight: 560)
    }
    .defaultSize(width: 1100, height: 740)
    Settings { ConnectionSettings() }
  }
}
