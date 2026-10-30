import Foundation
import HISAPI

extension Gateway {
  func importDiagnoses(contents: Data) async throws -> ImportResult {
    guard let url = URL(string: address + AppConfiguration.diagnosesImportPath)
    else {
      throw URLError(.badURL)
    }
    var request = URLRequest(url: url, timeoutInterval: AppConfiguration.requestTimeout)
    request.httpMethod = "POST"
    request.setValue("application/json", forHTTPHeaderField: "Content-Type")
    request.httpBody = contents

    if let tokens {
      let session = try await tokens.current()
      request.setValue("Bearer \(session.accessToken)", forHTTPHeaderField: "Authorization")
    }
    let (data, response) = try await CertificateSession.shared.data(for: request)
    guard let response = response as? HTTPURLResponse else {
      throw URLError(.badServerResponse)
    }
    guard (200..<300).contains(response.statusCode) else {
      throw ErrorResponse.error(response.statusCode, data, response, URLError(.badServerResponse))
    }
    return try diagnoses.apiConfiguration.codableHelper.decode(ImportResult.self, from: data).get()
  }
}
