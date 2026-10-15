import Foundation
import HISAPI

struct Gateway: Sendable {
  let address: String

  private var configuration: HISAPIAPIConfiguration {
    let configuration = HISAPIAPIConfiguration()
    configuration.basePath = address
    return configuration
  }

  var staff: StaffAPI { StaffAPI(apiConfiguration: configuration) }
  var diagnoses: DiagnosesAPI { DiagnosesAPI(apiConfiguration: configuration) }
  var patients: PatientsAPI { PatientsAPI(apiConfiguration: configuration) }
  var records: RecordsAPI { RecordsAPI(apiConfiguration: configuration) }

  func importDiagnoses(from file: URL) async throws -> Int {
    let accessed = file.startAccessingSecurityScopedResource()
    defer { if accessed { file.stopAccessingSecurityScopedResource() } }
    let contents = try Data(contentsOf: file)
    guard let url = URL(string: address + "/api/v1/diagnoses/import") else {
      throw URLError(.badURL)
    }
    var request = URLRequest(url: url, timeoutInterval: 20)
    request.httpMethod = "POST"
    request.setValue("application/json", forHTTPHeaderField: "Content-Type")
    request.httpBody = contents
    let (data, response) = try await URLSession.shared.data(for: request)
    guard let response = response as? HTTPURLResponse else { throw URLError(.badServerResponse) }
    guard (200..<300).contains(response.statusCode) else {
      throw ErrorResponse.error(response.statusCode, data, response, URLError(.badServerResponse))
    }
    return try configuration.codableHelper.decode(ImportResult.self, from: data).get().items.count
  }
}

func serverDetail(_ error: any Error) -> ErrorDetail? {
  guard case ErrorResponse.error(_, let data?, _, _) = error else { return nil }
  return try? JSONDecoder().decode(ModelError.self, from: data).error
}

func errorMessage(_ error: any Error) -> String {
  if let detail = serverDetail(error) {
    return
      ([detail.message]
      + (detail.fields ?? [:]).sorted { $0.key < $1.key }
      .map { "\($0.key): \($0.value)" }).joined(separator: "\n")
  }
  if case ErrorResponse.error(let status, _, _, let underlying) = error {
    return status > 0 ? "Сервер вернул ошибку HTTP \(status)." : underlying.localizedDescription
  }
  return error.localizedDescription
}
