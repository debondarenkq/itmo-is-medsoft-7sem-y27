import HISAPI
import SwiftUI

private enum SidebarSection { case patients, records }

struct MainView: View {
  let gateway: Gateway
  let profile: AppProfile
  let logout: () -> Void
  @State private var editingRecord = false
  @State private var confirmingLogout = false
  @State private var selectedSection: SidebarSection? = .patients

  var body: some View {
    NavigationSplitView {
      List(selection: $selectedSection) {
        Section("Рабочее место") {
          Label("Пациенты", systemImage: "person.2").tag(SidebarSection.patients)
          Label("ЭМК", systemImage: "heart.text.clipboard").tag(SidebarSection.records)
        }
      }
      .disabled(editingRecord)
      .navigationTitle("HIS — Медработник")
      .safeAreaInset(edge: .bottom) {
        VStack(alignment: .leading) {
          Divider()
          Label(profile.name, systemImage: profile.icon)
          Text(profile.subtitle).font(.caption).foregroundStyle(.secondary)
          Button("Выйти") { if editingRecord { confirmingLogout = true } else { logout() } }
        }.padding()
      }
    } detail: {
      if case .worker(let worker) = profile {
        PatientsView(gateway: gateway, worker: worker, recordsOnly: selectedSection == .records,
          editingChanged: { editingRecord = $0 }).id(selectedSection)
      }
    }
    .confirmationDialog("Закрыть редактирование и выйти?", isPresented: $confirmingLogout) {
      Button("Выйти", role: .destructive, action: logout)
      Button("Продолжить редактирование", role: .cancel) {}
    } message: { Text("Несохранённые изменения не попадут в ЭМК.") }
  }
}
