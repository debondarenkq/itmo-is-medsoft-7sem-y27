import Foundation
import HISAPI

extension FhirFHIRBundle {
  func resources<T: Decodable>(_ type: T.Type) throws -> [T] {
    let coder = CodableHelper()
    return try entry.map { try coder.jsonDecoder.decode(type, from: coder.jsonEncoder.encode($0.resource)) }
  }
}

extension FhirFHIRResource {
  var fullName: String {
    guard let name = name?.first else { return id.uuidString }
    return ([name.family] + name.given).joined(separator: " ")
  }
}

extension FhirFHIRServiceRequest {
  var testName: String { code.concept.coding.first?.display ?? "Анализ" }
  var loincCode: String { code.concept.coding.first?.code ?? "" }
}
