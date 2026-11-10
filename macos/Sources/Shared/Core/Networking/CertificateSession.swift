import Foundation
import HISAPI
import Security

// Observe the actual TLS challenge, then leave certificate validation to macOS.
final class ServerCertificateDelegate: NSObject, URLSessionDelegate {
  func urlSession(_ session: URLSession, didReceive challenge: URLAuthenticationChallenge,
    completionHandler: @Sendable @escaping (URLSession.AuthChallengeDisposition, URLCredential?) -> Void) {
    if challenge.protectionSpace.authenticationMethod == NSURLAuthenticationMethodServerTrust,
      let trust = challenge.protectionSpace.serverTrust,
      let chain = SecTrustCopyCertificateChain(trust) as? [SecCertificate] {
      let connection = CertificateConnection(
        id: "\(challenge.protectionSpace.host):\(challenge.protectionSpace.port)",
        receivedAt: Date(), certificates: chain.map { ServerCertificate($0) })
      Task { @MainActor in ServerCertificates.shared.record(connection) }
    }
    completionHandler(.performDefaultHandling, nil)
  }
}

enum CertificateSession {
  static let shared = URLSession(configuration: .default, delegate: ServerCertificateDelegate(), delegateQueue: nil)
}

// The generated client provides these extension points; generated code stays untouched.
struct CertificateRequestBuilderFactory: RequestBuilderFactory {
  func getNonDecodableBuilder<T>() -> RequestBuilder<T>.Type { CertificateRequestBuilder<T>.self }
  func getBuilder<T: Decodable>() -> RequestBuilder<T>.Type { CertificateDecodableRequestBuilder<T>.self }
}

private final class CertificateRequestBuilder<T: Sendable>: URLSessionRequestBuilder<T>, @unchecked Sendable {
  override func createURLSession() -> URLSessionProtocol { CertificateSession.shared }
}

private final class CertificateDecodableRequestBuilder<T: Decodable & Sendable>:
  URLSessionDecodableRequestBuilder<T>, @unchecked Sendable {
  override func createURLSession() -> URLSessionProtocol { CertificateSession.shared }
}
