import CryptoKit
import Foundation
import Security

struct ServerCertificate: Identifiable, Sendable {
  let subject: String
  let issuer: String
  let validFrom: String
  let validUntil: String
  let fingerprint: String
  let pem: String
  var id: String { fingerprint }

  init(_ certificate: SecCertificate) {
    subject = Self.name(certificate, oid: kSecOIDX509V1SubjectName)
    issuer = Self.name(certificate, oid: kSecOIDX509V1IssuerName)
    validFrom = Self.date(certificate, oid: kSecOIDX509V1ValidityNotBefore)
    validUntil = Self.date(certificate, oid: kSecOIDX509V1ValidityNotAfter)
    let data = SecCertificateCopyData(certificate) as Data
    fingerprint = SHA256.hash(data: data).map { String(format: "%02X", $0) }.joined(separator: ":")
    let encoded = data.base64EncodedString(options: [.lineLength64Characters, .endLineWithLineFeed])
    pem = "-----BEGIN CERTIFICATE-----\n\(encoded)\n-----END CERTIFICATE-----"
  }

  private static func value(_ certificate: SecCertificate, oid: CFString) -> Any? {
    guard let values = SecCertificateCopyValues(certificate, [oid] as CFArray, nil) as? [String: Any],
      let property = values[oid as String] as? [String: Any]
    else { return nil }
    return property[kSecPropertyKeyValue as String]
  }

  private static func name(_ certificate: SecCertificate, oid: CFString) -> String {
    guard let fields = value(certificate, oid: oid) as? [[String: Any]] else { return "Недоступно" }
    let names = ["2.5.4.3": "CN", "2.5.4.6": "C", "2.5.4.7": "L",
      "2.5.4.8": "ST", "2.5.4.10": "O", "2.5.4.11": "OU"]
    return fields.map { field in
      let label = field[kSecPropertyKeyLabel as String] as? String ?? ""
      let value = field[kSecPropertyKeyValue as String] as? String ?? ""
      return "\(names[label] ?? label)=\(value)"
    }.joined(separator: ", ")
  }

  private static func date(_ certificate: SecCertificate, oid: CFString) -> String {
    let raw = value(certificate, oid: oid)
    let date: Date?
    if let value = raw as? Date { date = value }
    else if let seconds = raw as? TimeInterval { date = Date(timeIntervalSinceReferenceDate: seconds) }
    else { date = nil }
    return date?.formatted(date: .numeric, time: .standard) ?? "Недоступно"
  }
}
