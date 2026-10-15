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
    defer {
      if accessed { file.stopAccessingSecurityScopedResource() }
    }
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
    let result = try configuration.codableHelper.decode(ImportResult.self, from: data).get()
    return result.items.count
  }
}

func serverDetail(_ error: any Error) -> ErrorDetail? {
  guard let responseError = error as? ErrorResponse else { return nil }
  switch responseError {
  case .error(_, let body, _, _):
    guard let body else { return nil }
    let response = try? JSONDecoder().decode(ModelError.self, from: body)
    return response?.error
  }
}

func errorMessage(_ error: any Error) -> String {
  if let detail = serverDetail(error) {
    var lines = [detail.message]
    let fields = detail.fields ?? [:]
    for field in fields.keys.sorted() {
      if let description = fields[field] {
        lines.append("\(fieldLabel(field)): \(description)")
      }
    }
    return lines.joined(separator: "\n")
  }
  if let responseError = error as? ErrorResponse {
    switch responseError {
    case .error(let statusCode, _, _, let underlyingError):
      if statusCode > 0 { return "Сервер вернул ошибку HTTP \(statusCode)." }
      return underlyingError.localizedDescription
    }
  }
  return error.localizedDescription
}

private func fieldLabel(_ field: String) -> String {
  let path = field.split(separator: ".")
  guard let property = path.last else { return field }
  let label = fieldName(String(property))

  if let segment = path.dropLast().last, let index = Int(segment) {
    return "\(label) №\(index + 1)"
  }
  return label
}

private func fieldName(_ property: String) -> String {
  switch property {
  case "first_name": return "Имя"
  case "last_name": return "Фамилия"
  case "middle_name": return "Отчество"
  case "position": return "Должность"
  case "birth_date": return "Дата рождения"
  case "administrative_sex": return "Пол"
  case "comment": return "Комментарий"
  case "code": return "Код"
  case "name": return "Наименование"
  case "text": return "Назначение"
  default: return property
  }
}
