import HISAPI
import SwiftUI

extension Staff {
  var fullName: String { "\(lastName) \(firstName)" }
}

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

struct ErrorNotice: View {
  let message: String?

  var body: some View {
    if let message {
      Label(message, systemImage: "exclamationmark.triangle")
        .foregroundStyle(.red)
        .textSelection(.enabled)
        .padding()
        .frame(maxWidth: .infinity, alignment: .leading)
    }
  }
}

struct PageControls: View {
  @Binding var offset: Int
  let total: Int
  let count: Int
  var limit = 50

  var body: some View {
    HStack {
      Text(total == 0 ? "Нет записей" : "\(offset + 1)–\(offset + count) из \(total)")
        .foregroundStyle(.secondary)
      Spacer()
      Button("Предыдущая", systemImage: "chevron.left") { offset = max(0, offset - limit) }
        .disabled(offset == 0)
      Button("Следующая", systemImage: "chevron.right") { offset += limit }
        .disabled(offset + count >= total)
    }
    .padding()
  }
}

struct FormActions: View {
  let isBusy: Bool
  var saveTitle = "Сохранить"
  let cancel: () -> Void
  let save: () -> Void

  var body: some View {
    HStack {
      if isBusy { ProgressView().controlSize(.small) }
      Spacer()
      Button("Отмена", role: .cancel, action: cancel).keyboardShortcut(.cancelAction)
      Button(saveTitle, action: save).keyboardShortcut(.defaultAction)
    }
    .disabled(isBusy)
    .padding()
  }
}
