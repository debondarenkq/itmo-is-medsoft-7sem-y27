enum AppConfiguration {
  #if HIS_ADMIN
  static let isAdministrator = true
  #else
  static let isAdministrator = false
  #endif

  static func accepts(role: String) -> Bool {
    isAdministrator ? role == "admin" : role != "admin"
  }

  #if LIS_APP
  static let defaultGatewayAddress = "https://localhost:8444"
  static let gatewayAddressKey = "lis.backend.address"
  #else
  static let defaultGatewayAddress = "https://localhost:8443"
  static let gatewayAddressKey = "gatewayAddress"
  #endif
  static let diagnosesImportPath = "/api/v1/diagnoses/import"
  static let requestTimeout: Double = 20
}
