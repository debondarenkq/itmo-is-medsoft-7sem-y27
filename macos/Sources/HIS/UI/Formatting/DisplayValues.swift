import Foundation
import HISAPI

extension Staff {
  var fullName: String { "\(lastName) \(firstName)" }
}

extension Patient {
  var age: Int { Calendar.current.dateComponents([.year], from: birthDate.calendarDate, to: Date()).year ?? 0 }
  var recordLabel: String { recordNumber.map { String(format: "ЭМК-%06lld", $0) } ?? "Не создана" }
  var fullName: String {
    var parts = [lastName, firstName]
    if let middleName { parts.append(middleName) }
    return parts.joined(separator: " ")
  }
}

extension AdministrativeSex {
  var label: String {
    switch self {
    case .m: return "Мужской"
    case .f: return "Женский"
    case .unknown: return "Не указан"
    }
  }
}

extension OpenAPIDateWithoutTime {
  var calendarDate: Date {
    let utc = TimeZone(secondsFromGMT: 0)!
    var calendar = Calendar(identifier: .gregorian)
    calendar.timeZone = utc
    let components = calendar.dateComponents([.year, .month, .day], from: wrappedDate)
    calendar.timeZone = .current
    return calendar.date(from: components) ?? wrappedDate
  }
}

extension EventType {
  var label: String {
    switch self {
    case .recordCreated: return "Создана медицинская карта"
    case .diagnosisAdded: return "Добавлен диагноз"
    case .diagnosisUpdated: return "Изменён диагноз"
    case .diagnosisRemoved: return "Удалён диагноз"
    case .prescriptionAdded: return "Добавлено назначение"
    case .prescriptionUpdated: return "Изменено назначение"
    case .prescriptionCancelled: return "Отменено назначение"
    case .prescriptionRemoved: return "Удалено назначение"
    }
  }
}
