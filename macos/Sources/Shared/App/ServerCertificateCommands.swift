import SwiftUI

struct ServerCertificateCommands: Commands {
  @Environment(\.openWindow) private var openWindow

  var body: some Commands {
    CommandGroup(after: .appSettings) {
      Button("Сертификаты сервера…") { openWindow(id: "server-certificates") }
    }
  }
}
