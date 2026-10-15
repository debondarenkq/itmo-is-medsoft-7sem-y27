import HISAPI
import SwiftUI

extension Patient {
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
