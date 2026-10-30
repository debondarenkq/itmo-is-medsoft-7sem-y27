import HISAPI
import SwiftUI

private enum AccountStatusFilter: String, CaseIterable {
  case all = "Все статусы", active = "Активные", inactive = "Неактивные"
  var value: Bool? {
    switch self { case .all: return nil; case .active: return true; case .inactive: return false }
  }
}

struct AccountsView: View {
  let gateway: Gateway
  @State private var accounts: [Account] = []
  @State private var selectedID: UUID?
  @State private var search = ""
  @State private var status = AccountStatusFilter.all
  @State private var offset = 0
  @State private var total = 0
  @State private var errorText: String?
  @State private var isBusy = false
  @State private var creating = false
  @State private var editing: Account?
  @State private var resettingPassword = false
  @State private var changingStatus = false
  @State private var statusAccount: Account?
  @State private var desiredActivity = false
  @State private var generated: AccountPassword?
  private var selected: Account? { accounts.first { $0.id == selectedID } }

  var body: some View {
    VStack(spacing: 0) {
      HStack {
        Picker("Статус", selection: $status) {
          ForEach(AccountStatusFilter.allCases, id: \.self) { Text($0.rawValue).tag($0) }
        }.frame(width: 240)
        Spacer()
        Text("Учётные записи пользователей").foregroundStyle(.secondary)
      }.padding()
      Divider()
      ErrorNotice(message: errorText)
      Table(accounts, selection: $selectedID) {
        TableColumn("Логин", value: \.username)
        TableColumn("Медработник") { account in
          Text(account.staffName.isEmpty ? (account.role == .admin ? "Администратор" : "Не связан") : account.staffName)
        }
        TableColumn("Роль") { Text($0.role == .admin ? "Администратор" : (StaffRole(rawValue: $0.role.rawValue)?.label ?? "—")) }
        TableColumn("Статус") { account in
          HStack {
            Toggle("Активен", isOn: Binding(get: { account.isActive }, set: { value in
              statusAccount = account; desiredActivity = value; changingStatus = true
            })).toggleStyle(.switch).labelsHidden().disabled(account.staffId == nil || isBusy)
            Text(account.isActive ? "Активен" : "Неактивен").foregroundStyle(account.isActive ? Color.green : Color.secondary)
          }
        }
        TableColumn("Создан") { Text($0.createdAt, format: .dateTime.day().month().year()) }
        TableColumn("Действия") { account in
          Button("Изменить", systemImage: "square.and.pencil") { editing = account }
            .labelStyle(.iconOnly).buttonStyle(.borderless).disabled(account.staffId == nil || isBusy)
        }.width(70)
      }
      PageControls(offset: $offset, total: total, count: accounts.count)
    }
    .navigationTitle("Аккаунты")
    .searchable(text: $search, prompt: "Логин или медработник")
    .toolbar {
      Button("Создать аккаунт", systemImage: "person.badge.plus") { creating = true }
      Button("Изменить", systemImage: "square.and.pencil") { editing = selected }.disabled(selected?.staffId == nil)
      Button("Новый пароль", systemImage: "key") { resettingPassword = true }.disabled(selected?.staffId == nil)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .task(id: "\(offset)|\(search)|\(status)") {
      do { try await Task.sleep(for: .milliseconds(250)); await load() } catch {}
    }
    .onChange(of: search) { offset = 0; selectedID = nil }
    .onChange(of: status) { offset = 0; selectedID = nil }
    .sheet(isPresented: $creating, onDismiss: { Task { await load() } }) {
      AccountCreationView(gateway: gateway, staff: nil)
    }
    .sheet(item: $editing, onDismiss: { Task { await load() } }) { account in
      AccountEditorView(gateway: gateway, account: account)
    }
    .confirmationDialog("Сгенерировать новый пароль?", isPresented: $resettingPassword) {
      Button("Сгенерировать", role: .destructive) { Task { await resetPassword() } }
    } message: { Text("Передайте новый пароль медработнику. Обновление прежних сессий будет прекращено.") }
    .confirmationDialog(desiredActivity ? "Активировать аккаунт?" : "Деактивировать аккаунт?", isPresented: $changingStatus) {
      Button("Подтвердить") { Task { await changeStatus() } }
    } message: { Text(statusAccount?.username ?? "") }
    .sheet(isPresented: Binding(get: { generated != nil }, set: { if !$0 { generated = nil } })) {
      VStack(spacing: 16) {
        Text("Новый пароль для \(generated?.account.username ?? "")").font(.headline)
        Text(generated?.password ?? "").monospaced().textSelection(.enabled)
        Text("Пароль показывается только сейчас.").foregroundStyle(.secondary)
        Button("Закрыть") { generated = nil }
      }.padding(24).frame(width: 420)
    }
  }

  private func load() async {
    do {
      let page = try await gateway.security.listAccounts(limit: 50, offset: offset, search: search, isActive: status.value)
      try Task.checkCancellation()
      accounts = page.items; total = page.total; errorText = nil
    } catch { if !Task.isCancelled { errorText = errorMessage(error) } }
  }

  private func resetPassword() async {
    guard let selected else { return }
    isBusy = true
    defer { isBusy = false }
    do { generated = try await gateway.security.resetPassword(id: selected.id); await load() }
    catch { errorText = errorMessage(error) }
  }

  private func changeStatus() async {
    guard let statusAccount else { return }
    isBusy = true
    defer { isBusy = false }
    do {
      _ = try await gateway.security.updateAccount(id: statusAccount.id, accountUpdate: AccountUpdate(isActive: desiredActivity))
      await load()
    } catch { errorText = errorMessage(error) }
  }
}

extension Account: Identifiable {}
