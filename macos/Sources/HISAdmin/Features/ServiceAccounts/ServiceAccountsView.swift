import HISAPI
import SwiftUI

struct AccountDirectoryView: View {
  let gateway: Gateway
  var body: some View {
    TabView {
      AccountsView(gateway: gateway).tabItem { Label("Пользовательские", systemImage: "person.crop.circle") }
      ServiceAccountsView(gateway: gateway).tabItem { Label("Сервисные", systemImage: "key.horizontal") }
    }
  }
}

struct ServiceAccountsView: View {
  let gateway: Gateway
  @State private var accounts: [ServiceAccount] = []
  @State private var selectedID: UUID?
  @State private var creating = false
  @State private var editing: ServiceAccount?
  @State private var generated: ServiceAccountKey?
  @State private var errorText: String?
  @State private var busy = false
  @State private var rotating = false
  private var selected: ServiceAccount? { accounts.first { $0.id == selectedID } }

  var body: some View {
    VStack {
      ErrorNotice(message: errorText)
      Table(accounts, selection: $selectedID) {
        TableColumn("Название", value: \.name)
        TableColumn("Роль") { Text($0.role.rawValue.uppercased()) }
        TableColumn("Статус") { Text($0.isActive ? "Активен" : "Отключён") }
        TableColumn("Создан") { Text($0.createdAt, format: .dateTime.day().month().year()) }
      }
      Text("LIS: чтение Patient, Practitioner, ServiceRequest и Observation; создание Observation.")
        .font(.caption).foregroundStyle(.secondary).padding()
    }
    .navigationTitle("Сервисные аккаунты")
    .toolbar {
      Button("Создать", systemImage: "plus") { creating = true }
      Button("Изменить", systemImage: "square.and.pencil") { editing = selected }.disabled(selected == nil)
      Button("Заменить API-ключ", systemImage: "key") { rotating = true }.disabled(selected == nil)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .disabled(busy)
    .task { await load() }
    .sheet(isPresented: $creating, onDismiss: { Task { await load() } }) {
      ServiceAccountEditor(gateway: gateway, account: nil) { generated = $0 }
    }
    .sheet(item: $editing, onDismiss: { Task { await load() } }) {
      ServiceAccountEditor(gateway: gateway, account: $0) { generated = $0 }
    }
    .confirmationDialog("Заменить API-ключ?", isPresented: $rotating) {
      Button("Заменить ключ") { Task { await rotate() } }
    } message: { Text("Старый ключ и выданные сессии перестанут действовать. Новый ключ нужно ввести в LIS.") }
    .sheet(isPresented: Binding(get: { generated != nil && !creating && editing == nil }, set: { if !$0 { generated = nil } })) {
      VStack(alignment: .leading, spacing: 18) {
        Text("API-ключ: \(generated?.account.name ?? "")").font(.title2)
        Text(generated?.apiKey ?? "").monospaced().textSelection(.enabled)
        Text("Ключ показывается только сейчас. Введите его в настройках LIS.").foregroundStyle(.secondary)
        Button("Готово") { generated = nil }
      }.padding(24).frame(width: 620)
    }
  }

  private func load() async {
    do { accounts = try await gateway.serviceAccounts.listServiceAccounts(); errorText = nil }
    catch { errorText = errorMessage(error) }
  }

  private func rotate() async {
    guard let selected else { return }
    busy = true
    defer { busy = false }
    do { generated = try await gateway.serviceAccounts.rotateServiceAccountKey(id: selected.id); await load() }
    catch { errorText = errorMessage(error) }
  }
}

private struct ServiceAccountEditor: View {
  let gateway: Gateway
  let account: ServiceAccount?
  let created: (ServiceAccountKey) -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var name = ""
  @State private var active = true
  @State private var errorText: String?
  @State private var busy = false

  var body: some View {
    VStack(alignment: .leading, spacing: 16) {
      Text(account == nil ? "Новый сервисный аккаунт" : "Сервисный аккаунт").font(.title2)
      Form {
        TextField("Название", text: $name)
        LabeledContent("Роль", value: "LIS")
        if account != nil { Toggle("Активен", isOn: $active) }
      }
      ErrorNotice(message: errorText)
      HStack {
        Button("Отмена") { dismiss() }.keyboardShortcut(.cancelAction)
        Spacer()
        Button("Сохранить") { Task { await save() } }.buttonStyle(.borderedProminent)
      }
    }.padding(24).frame(width: 460).disabled(busy)
    .onAppear { name = account?.name ?? ""; active = account?.isActive ?? true }
  }

  private func save() async {
    busy = true
    defer { busy = false }
    do {
      if let account {
        _ = try await gateway.serviceAccounts.updateServiceAccount(id: account.id, serviceAccountUpdate: ServiceAccountUpdate(isActive: active, name: name))
      } else {
        created(try await gateway.serviceAccounts.createServiceAccount(serviceAccountInput: ServiceAccountInput(name: name)))
      }
      dismiss()
    } catch { errorText = errorMessage(error) }
  }
}
