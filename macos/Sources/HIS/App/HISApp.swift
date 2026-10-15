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
  @State private var selection: Workspace? = .patients

  var body: some View {
    NavigationSplitView {
      List(selection: $selection) {
        Section("Медработник") {
          Label(Workspace.patients.rawValue, systemImage: Workspace.patients.icon).tag(
            Workspace.patients)
        }
        Section("Администрирование") {
          ForEach([Workspace.staff, .diagnoses]) { section in
            Label(section.rawValue, systemImage: section.icon).tag(section)
          }
        }
      }
      .navigationTitle("HIS")
    } detail: {
      Group {
        switch selection {
        case .staff: StaffView(gateway: gateway)
        case .diagnoses: DiagnosesView(gateway: gateway)
        default: PatientsView(gateway: gateway)
        }
      }
      .id(gateway.address)
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
