import Foundation
import HISAPI

extension DiagnosesAPI {
  func importDiagnoses(contents: Data) async throws -> ImportResult {
    guard let url = URL(string: apiConfiguration.basePath + AppConfiguration.diagnosesImportPath)
    else {
      throw URLError(.badURL)
    }
    var request = URLRequest(url: url, timeoutInterval: AppConfiguration.requestTimeout)
    request.httpMethod = "POST"
    request.setValue("application/json", forHTTPHeaderField: "Content-Type")
    request.httpBody = contents

    let (data, response) = try await URLSession.shared.data(for: request)
    guard let response = response as? HTTPURLResponse else {
      throw URLError(.badServerResponse)
    }
    guard (200..<300).contains(response.statusCode) else {
      throw ErrorResponse.error(response.statusCode, data, response, URLError(.badServerResponse))
    }
    return try apiConfiguration.codableHelper.decode(ImportResult.self, from: data).get()
  }
}
