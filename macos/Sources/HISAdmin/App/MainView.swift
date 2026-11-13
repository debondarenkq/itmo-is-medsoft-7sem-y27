import SwiftUI

private enum SidebarSection { case staff, diagnoses, labTests, accounts, audit }

struct MainView: View {
  let gateway: Gateway
  let profile: AppProfile
  let logout: () -> Void
  @State private var selectedSection: SidebarSection? = .staff

  var body: some View {
    NavigationSplitView {
      List(selection: $selectedSection) {
        Section("Администрирование") {
          Label("Медработники", systemImage: "person.badge.key").tag(SidebarSection.staff)
          Label("Диагнозы", systemImage: "list.clipboard").tag(SidebarSection.diagnoses)
          Label("Лабораторные анализы", systemImage: "testtube.2").tag(SidebarSection.labTests)
          Label("Аккаунты", systemImage: "person.crop.circle").tag(SidebarSection.accounts)
          Label("Аудит", systemImage: "list.bullet.rectangle").tag(SidebarSection.audit)
        }
      }
      .navigationTitle("HIS — Администрирование")
      .safeAreaInset(edge: .bottom) {
        VStack(alignment: .leading) {
          Divider()
          Label(profile.name, systemImage: profile.icon)
          Button("Выйти", action: logout)
        }.padding()
      }
    } detail: {
      switch selectedSection {
      case .diagnoses: DiagnosesView(gateway: gateway)
      case .labTests: LabTestsView(gateway: gateway)
      case .accounts: AccountDirectoryView(gateway: gateway)
      case .audit: AuditView(gateway: gateway)
      default: StaffView(gateway: gateway)
      }
    }
  }
}
