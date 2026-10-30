import HISAPI
import Foundation

struct Gateway: Sendable {
  let address: String
  let tokens: SessionTokens?

  init(address: String, tokens: SessionTokens? = nil) {
    self.address = URLComponents(string: address)?.scheme == "https"
      ? address : AppConfiguration.defaultGatewayAddress
    self.tokens = tokens
  }

  private var configuration: HISAPIAPIConfiguration {
    let configuration = HISAPIAPIConfiguration()
    configuration.basePath = address
    configuration.requestBuilderFactory = CertificateRequestBuilderFactory()
    if let tokens { configuration.interceptor = SessionInterceptor(tokens: tokens) }
    return configuration
  }

  var security: SecurityAPI { SecurityAPI(apiConfiguration: configuration) }

  var staff: StaffAPI { StaffAPI(apiConfiguration: configuration) }
  var diagnoses: DiagnosesAPI { DiagnosesAPI(apiConfiguration: configuration) }
  var patients: PatientsAPI { PatientsAPI(apiConfiguration: configuration) }
  var records: RecordsAPI { RecordsAPI(apiConfiguration: configuration) }
}
