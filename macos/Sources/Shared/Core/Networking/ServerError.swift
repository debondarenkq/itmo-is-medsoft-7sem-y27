import Foundation
import HISAPI

func serverDetail(_ error: any Error) -> ErrorDetail? {
  guard let responseError = error as? ErrorResponse else { return nil }
  switch responseError {
  case .error(_, let body, _, _):
    guard let body else { return nil }
    let response = try? JSONDecoder().decode(ModelError.self, from: body)
    return response?.error
  }
}
