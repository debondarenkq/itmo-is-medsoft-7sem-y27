import HISAPI
import SwiftUI

private enum SidebarSection { case patients, records, staff, diagnoses, accounts, audit }

struct MainView: View {
  let gateway: Gateway
  let profile: AppProfile
  let logout: () -> Void
  @State private var editingRecord = false
  @State private var confirmingLogout = false
  @State private var selectedSection: SidebarSection? = .staff

  var body: some View {
    NavigationSplitView {
      List(selection: $selectedSection) {
        switch profile {
        case .administrator:
          Section("Администрирование") {
            Label("Медработники", systemImage: "person.badge.key").tag(SidebarSection.staff)
            Label("Диагнозы", systemImage: "list.clipboard").tag(SidebarSection.diagnoses)
            Label("Аккаунты", systemImage: "person.crop.circle").tag(SidebarSection.accounts)
            Label("Аудит", systemImage: "list.bullet.rectangle").tag(SidebarSection.audit)
          }
        case .worker:
          Section("Рабочее место") {
            Label("Пациенты", systemImage: "person.2").tag(SidebarSection.patients)
            Label("ЭМК", systemImage: "heart.text.clipboard").tag(SidebarSection.records)
          }
        }
      }
      .disabled(editingRecord)
      .navigationTitle("HIS")
      .safeAreaInset(edge: .bottom) {
        VStack(alignment: .leading) {
          Divider()
          Label(profile.name, systemImage: profile.icon)
          Text(profile.subtitle).font(.caption).foregroundStyle(.secondary)
          Button("Выйти") { if editingRecord { confirmingLogout = true } else { logout() } }
        }.padding()
      }
    } detail: {
      switch profile {
      case .administrator:
        switch selectedSection {
        case .diagnoses: DiagnosesView(gateway: gateway)
        case .accounts: AccountsView(gateway: gateway)
        case .audit: AuditView(gateway: gateway)
        default: StaffView(gateway: gateway)
        }
      case .worker(let worker):
        PatientsView(gateway: gateway, worker: worker, recordsOnly: selectedSection == .records, editingChanged: { editingRecord = $0 }).id(selectedSection)
      }
    }
    .confirmationDialog("Закрыть редактирование и выйти?", isPresented: $confirmingLogout) {
      Button("Выйти", role: .destructive, action: logout)
      Button("Продолжить редактирование", role: .cancel) {}
    } message: { Text("Несохранённые изменения не попадут в ЭМК.") }
  }
}
