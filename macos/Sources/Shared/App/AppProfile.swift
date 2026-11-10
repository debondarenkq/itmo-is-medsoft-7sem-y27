import HISAPI

enum AppProfile: Hashable {
  case administrator
  case worker(Staff)

  var name: String {
    switch self {
    case .administrator: return "Администратор"
    case .worker(let worker): return worker.fullName
    }
  }

  var subtitle: String {
    switch self {
    case .administrator: return "Администрирование"
    case .worker(let worker): return worker.position
    }
  }

  var icon: String {
    switch self {
    case .administrator: return "person.badge.key"
    case .worker: return "stethoscope"
    }
  }
}
