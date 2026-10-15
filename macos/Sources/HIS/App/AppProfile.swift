import HISAPI

enum AppProfile: Hashable {
  case administrator
  case worker(Staff)

  var name: String {
    switch self {
    case .administrator: "Администратор"
    case .worker(let worker): worker.fullName
    }
  }

  var subtitle: String {
    switch self {
    case .administrator: "Управление справочниками"
    case .worker(let worker): worker.position
    }
  }

  var icon: String {
    switch self {
    case .administrator: "person.badge.key"
    case .worker: "stethoscope"
    }
  }
}
