import HISAPI
import SwiftUI

struct AccountEditorView: View {
  let gateway: Gateway
  let account: Account
  @Environment(\.dismiss) private var dismiss
  @State private var username: String
  @State private var role: StaffRole
  @State private var isActive: Bool
  @State private var isBusy = false
  @State private var errorText: String?

  init(gateway: Gateway, account: Account) {
    self.gateway = gateway
    self.account = account
    _username = State(initialValue: account.username)
    _role = State(initialValue: StaffRole(rawValue: account.role.rawValue) ?? .other)
    _isActive = State(initialValue: account.isActive)
  }

  var body: some View {
    VStack(spacing: 0) {
      Form {
        Section("Аккаунт медработника") {
          LabeledContent("Медработник", value: account.staffName)
          TextField("Логин", text: $username)
          Picker("Роль", selection: $role) {
            ForEach(StaffRole.allCases, id: \.self) { Text($0.label).tag($0) }
          }
          Toggle("Аккаунт активен", isOn: $isActive)
          Text("Изменение роли или отключение аккаунта завершит возможность обновления его сессий.")
            .font(.caption).foregroundStyle(.secondary)
        }
      }.formStyle(.grouped)
      ErrorNotice(message: errorText)
      FormActions(isBusy: isBusy, cancel: { dismiss() }, save: { Task { await save() } })
    }.frame(width: 500).disabled(isBusy).interactiveDismissDisabled()
  }

  private func save() async {
    isBusy = true
    defer { isBusy = false }
    do {
      _ = try await gateway.security.updateAccount(id: account.id,
        accountUpdate: AccountUpdate(isActive: isActive,
          role: AccountUpdate.Role(rawValue: role.rawValue), username: username))
      dismiss()
    } catch { errorText = errorMessage(error) }
  }
}
