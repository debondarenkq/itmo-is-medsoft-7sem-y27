import HISAPI
import SwiftUI

extension AuditEventKind {
  var label: String {
    switch self {
    case .adminSetup: return "Настройка администратора"
    case .login: return "Вход в систему"
    case .loginFailed: return "Неуспешный вход"
    case .logout: return "Выход из системы"
    case .accessDenied: return "Отказ в доступе"
    case .staffCreated: return "Создание медработника"
    case .staffUpdated: return "Изменение медработника"
    case .staffDeleted: return "Удаление медработника"
    case .patientCreated: return "Создание пациента"
    case .patientUpdated: return "Изменение пациента"
    case .patientDeleted: return "Удаление пациента"
    case .recordCreated: return "Создание ЭМК"
    case .recordUpdated: return "Изменение ЭМК"
    case .accountCreated: return "Создание аккаунта"
    case .accountUpdated: return "Изменение аккаунта"
    case .tokenRefreshed: return "Обновление сессии"
    case .other: return "Другое событие"
    }
  }

  var icon: String {
    switch self {
    case .loginFailed, .accessDenied: return "exclamationmark.shield"
    case .login, .logout, .adminSetup, .tokenRefreshed: return "person.badge.key"
    case .patientCreated, .patientUpdated, .patientDeleted: return "person.crop.rectangle"
    case .recordCreated, .recordUpdated: return "heart.text.clipboard"
    case .staffDeleted: return "person.fill.xmark"
    case .staffCreated, .staffUpdated: return "stethoscope"
    case .accountCreated, .accountUpdated: return "person.crop.circle"
    case .other: return "info.circle"
    }
  }

  var color: Color {
    switch self {
    case .loginFailed, .accessDenied, .staffDeleted, .patientDeleted: return .red
    case .staffCreated, .patientCreated, .recordCreated, .accountCreated, .adminSetup: return .green
    default: return .accentColor
    }
  }
}
