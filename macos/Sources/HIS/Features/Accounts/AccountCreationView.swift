import HISAPI
import SwiftUI

struct AccountCreationView: View {
  let gateway: Gateway
  @Environment(\.dismiss) private var dismiss
  @State private var worker: Staff?
  @State private var choosingWorker = false
  @State private var username = ""
  @State private var created: AccountPassword?
  @State private var errorText: String?
  @State private var isBusy = false

  init(gateway: Gateway, staff: Staff?) {
    self.gateway = gateway
    _worker = State(initialValue: staff)
  }

  var body: some View {
    VStack(spacing: 16) {
      Text("Создание аккаунта").font(.title2)
      if let created {
        Text("Передайте логин и пароль медработнику. Пароль показывается только сейчас.")
        LabeledContent("Логин", value: created.account.username)
        LabeledContent("Пароль", value: created.password).monospaced().textSelection(.enabled)
      } else {
        Form {
          if let worker { LabeledContent("Медработник", value: worker.fullName) }
          Button(worker == nil ? "Выбрать медработника…" : "Выбрать другого…") { choosingWorker = true }
          TextField("Логин", text: $username)
          Text("От 3 до 16 символов. Пароль из 12 букв и цифр будет создан сервером.")
            .font(.caption).foregroundStyle(.secondary)
        }
        Button("Создать аккаунт") { Task { await create() } }.disabled(isBusy || worker == nil)
      }
      ErrorNotice(message: errorText)
      Button("Закрыть") { dismiss() }
    }.padding(24).frame(width: 500).interactiveDismissDisabled(isBusy)
    .sheet(isPresented: $choosingWorker) {
      DirectoryPicker<Staff>(title: "Выберите медработника",
        label: { "\($0.fullName) · \($0.role.label)" }, fetch: { offset in
          let page = try await gateway.staff.listStaff(limit: 50, offset: offset)
          return (page.items, page.total)
        }, select: { worker = $0 })
    }
  }

  private func create() async {
    guard let worker else { return }
    isBusy = true
    defer { isBusy = false }
    do {
      created = try await gateway.security.createAccount(accountInput: AccountInput(staffId: worker.id, username: username))
      errorText = nil
    } catch { errorText = errorMessage(error) }
  }
}
