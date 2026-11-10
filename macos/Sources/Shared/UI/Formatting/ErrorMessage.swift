import Foundation
import HISAPI

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
    case .error(let statusCode, let data, _, let underlyingError):
      if let data, let problem = try? JSONDecoder().decode(SecurityError.self, from: data) { return problem.error.message }
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
