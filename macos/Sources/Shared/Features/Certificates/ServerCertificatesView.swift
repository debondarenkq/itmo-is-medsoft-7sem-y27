import SwiftUI

struct ServerCertificatesView: View {
  @AppStorage(AppConfiguration.gatewayAddressKey)
  private var address = AppConfiguration.defaultGatewayAddress
  private let certificates = ServerCertificates.shared
  @State private var isBusy = false
  @State private var errorText: String?

  var body: some View {
    ScrollView {
      VStack(alignment: .leading, spacing: 20) {
        Text("Сертификаты HTTPS-соединений").font(.title2)
        Text("Данные получены клиентом при TLS-подключении. Цепочка serverTrust может включать корневой сертификат из хранилища доверия macOS. Показ сертификата сам по себе не означает, что он доверенный.")
          .foregroundStyle(.secondary)
        if certificates.connections.isEmpty {
          Text("Сертификаты ещё не получены. Подключитесь к серверу или нажмите «Проверить подключение».")
        }
        ForEach(certificates.connections) { connection in
          GroupBox {
            VStack(alignment: .leading, spacing: 16) {
              Text(connection.id).font(.headline)
              Text("Получены: \(connection.receivedAt.formatted(date: .numeric, time: .standard))")
              ForEach(Array(connection.certificates.enumerated()), id: \.offset) { index, certificate in
                Divider()
                Text(certificateTitle(certificate, index: index)).font(.headline)
                LabeledContent("Subject", value: certificate.subject)
                LabeledContent("Issuer", value: certificate.issuer)
                LabeledContent("Действует с", value: certificate.validFrom)
                LabeledContent("Действует до", value: certificate.validUntil)
                Text("SHA-256: \(certificate.fingerprint)").font(.system(.caption, design: .monospaced))
                DisclosureGroup("Полный сертификат (PEM)") {
                  Text(certificate.pem).font(.system(.caption, design: .monospaced))
                }
              }
            }.padding(8).frame(maxWidth: .infinity, alignment: .leading)
          }
        }
        ErrorNotice(message: errorText)
        Button("Проверить подключение") { Task { await connect() } }.disabled(isBusy)
        if isBusy { ProgressView() }
      }.padding().textSelection(.enabled)
    }
    .frame(minWidth: 650, idealWidth: 800, minHeight: 450, idealHeight: 650)
  }

  private func certificateTitle(_ certificate: ServerCertificate, index: Int) -> String {
    if index == 0 { return "Сертификат сервера" }
    if certificate.subject == certificate.issuer && certificate.subject != "Недоступно" {
      return "Корневой сертификат CA"
    }
    return "Промежуточный сертификат CA"
  }

  private func connect() async {
    guard let url = URL(string: Gateway(address: address).address + "/live") else { return }
    isBusy = true
    defer { isBusy = false }
    // A fresh session performs a new handshake, even when the API reuses a connection.
    let session = URLSession(configuration: .ephemeral, delegate: ServerCertificateDelegate(), delegateQueue: nil)
    defer { session.invalidateAndCancel() }
    do {
      _ = try await session.data(from: url)
      errorText = nil
    } catch { errorText = errorMessage(error) }
  }
}
