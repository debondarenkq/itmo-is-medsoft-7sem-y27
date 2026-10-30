import HISAPI

extension StaffRole {
  var label: String {
    switch self {
    case .doctor: return "Врач"
    case .laboratoryTechnician: return "Лаборант"
    case .nurse: return "Медсестра / медбрат"
    case .other: return "Другая"
    }
  }
  var canCreate: Bool { self == .doctor || self == .nurse }
}
