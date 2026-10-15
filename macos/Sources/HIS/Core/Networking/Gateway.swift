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
}
