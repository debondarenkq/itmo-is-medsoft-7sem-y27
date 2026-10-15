import HISAPI
import SwiftUI

@main
struct HISApp: App {
  @NSApplicationDelegateAdaptor(ApplicationDelegate.self) private var delegate
  @AppStorage("gatewayAddress") private var address = "http://localhost:8080"

  var body: some Scene {
    Window("Госпитальная информационная система", id: "main") {
      MainView(gateway: Gateway(address: address))
        .frame(minWidth: 850, minHeight: 560)
    }
    .defaultSize(width: 1100, height: 740)
    Settings { ConnectionSettings() }
  }
}

private enum Workspace: String, CaseIterable, Identifiable {
  case patients = "Пациенты"
  case staff = "Медработники"
  case diagnoses = "Диагнозы"
  var id: Self { self }
  var icon: String {
    switch self {
    case .patients: "person.2"
    case .staff: "person.badge.key"
    case .diagnoses: "list.clipboard"
    }
  }
}

private struct MainView: View {
  let gateway: Gateway
  @State private var profile: AppProfile?
  @State private var selection: Workspace? = .staff
  @State private var choosingWorker = false

  var body: some View {
    NavigationSplitView {
      List(selection: $selection) {
        switch profile {
        case .administrator:
          Section("Администрирование") {
            ForEach([Workspace.staff, .diagnoses]) { section in
              Label(section.rawValue, systemImage: section.icon).tag(section)
            }
          }
        case .worker:
          Label(Workspace.patients.rawValue, systemImage: Workspace.patients.icon).tag(
            Workspace.patients)
        case nil:
          EmptyView()
        }
      }
      .navigationTitle("HIS")
      .safeAreaInset(edge: .bottom, spacing: 0) {
        VStack(spacing: 0) {
          Divider()
          VStack(alignment: .leading) {
            Menu {
              Button("Администратор", systemImage: "person.badge.key") { select(.administrator) }
              Button("Медработник…", systemImage: "stethoscope") { choosingWorker = true }
            } label: {
              Label(
                profile?.name ?? "Выбрать профиль",
                systemImage: profile?.icon ?? "person.crop.circle")
            }
            .menuStyle(.borderlessButton)
            .help(
              profile.map { "\($0.name) · \($0.subtitle). Сменить профиль" } ?? "Выбрать профиль")
            if let profile {
              Text(profile.subtitle).font(.caption).foregroundStyle(.secondary).lineLimit(1)
            }
          }
          .padding()
        }
      }
    } detail: {
      Group {
        switch profile {
        case .administrator:
          if selection == .diagnoses {
            DiagnosesView(gateway: gateway)
          } else {
            StaffView(gateway: gateway)
          }
        case .worker(let worker):
          PatientsView(gateway: gateway, worker: worker)
        case nil:
          ContentUnavailableView {
            Label("Выберите профиль", systemImage: "person.crop.circle")
          } description: {
            Text("Войдите как администратор или выберите медработника.")
          } actions: {
            Button("Администратор") { select(.administrator) }
            Button("Медработник…") { choosingWorker = true }
          }
        }
      }
      .id(profile)
      .id(gateway.address)
    }
    .sheet(isPresented: $choosingWorker) {
      DirectoryPicker<Staff>(
        title: "Выберите профиль медработника",
        label: { "\($0.fullName) · \($0.position)" },
        fetch: { search, offset in
          let page = try await gateway.staff.listStaff(limit: 50, offset: offset, search: search)
          return (page.items, page.total)
        }, select: { select(.worker($0)) })
    }
    .onChange(of: gateway.address) { profile = nil }
  }

  private func select(_ profile: AppProfile) {
    self.profile = profile
    switch profile {
    case .administrator: selection = .staff
    case .worker: selection = .patients
    }
  }
}

private struct ConnectionSettings: View {
  @AppStorage("gatewayAddress") private var savedAddress = "http://localhost:8080"
  @State private var address = ""
  @State private var message: String?

  var body: some View {
    Form {
      TextField("Адрес сервера", text: $address)
      Text("Например, http://localhost:8080").foregroundStyle(.secondary)
      ErrorNotice(message: message)
      Button("Применить") {
        guard !NSApp.windows.contains(where: { $0.attachedSheet != nil }) else {
          message = "Закройте открытый редактор перед сменой сервера."
          return
        }
        let value = address.trimmingCharacters(in: .whitespacesAndNewlines)
        guard let url = URLComponents(string: value),
          ["http", "https"].contains(url.scheme),
          let host = url.host, !host.isEmpty,
          url.user == nil, url.password == nil, url.query == nil, url.fragment == nil
        else {
          message = "Укажите адрес вида http://localhost:8080."
          return
        }
        savedAddress = value.hasSuffix("/") ? String(value.dropLast()) : value
        message = nil
      }
    }
    .formStyle(.grouped)
    .frame(width: 480)
    .onAppear { address = savedAddress }
  }
}
