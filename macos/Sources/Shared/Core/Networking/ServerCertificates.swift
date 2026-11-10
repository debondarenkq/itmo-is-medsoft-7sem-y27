import Foundation
import Observation

struct CertificateConnection: Identifiable, Sendable {
  let id: String
  let receivedAt: Date
  let certificates: [ServerCertificate]
}

@MainActor @Observable
final class ServerCertificates {
  static let shared = ServerCertificates()
  private(set) var connections: [CertificateConnection] = []

  func record(_ connection: CertificateConnection) {
    connections.removeAll { $0.id == connection.id }
    connections.insert(connection, at: 0)
  }
}
