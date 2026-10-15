import HISAPI
import SwiftUI

private enum SidebarSection {
  case patients
  case staff
  case diagnoses
}

struct MainView: View {
  let gateway: Gateway
  @State private var selectedProfile: AppProfile?
  @State private var selectedSection: SidebarSection? = .staff
  @State private var showingWorkerPicker = false

  var body: some View {
    NavigationSplitView {
      List(selection: $selectedSection) {
        switch selectedProfile {
        case .administrator:
          Section("Администрирование") {
            Label("Медработники", systemImage: "person.badge.key").tag(SidebarSection.staff)
            Label("Диагнозы", systemImage: "list.clipboard").tag(SidebarSection.diagnoses)
          }
        case .worker:
          Label("Пациенты", systemImage: "person.2").tag(SidebarSection.patients)
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
              Button("Администратор", systemImage: "person.badge.key") {
                selectProfile(.administrator)
              }
              Button("Медработник…", systemImage: "stethoscope") { showingWorkerPicker = true }
            } label: {
              Label(
                selectedProfile?.name ?? "Выбрать профиль",
                systemImage: selectedProfile?.icon ?? "person.crop.circle")
            }
            .menuStyle(.borderlessButton)
            .help(profileHelpText)
            if let selectedProfile {
              Text(selectedProfile.subtitle).font(.caption).foregroundStyle(.secondary).lineLimit(1)
            }
          }
          .padding()
        }
      }
    } detail: {
      Group {
        switch selectedProfile {
        case .administrator:
          if selectedSection == .diagnoses {
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
            Button("Администратор") { selectProfile(.administrator) }
            Button("Медработник…") { showingWorkerPicker = true }
          }
        }
      }
      .id(selectedProfile)
      .id(gateway.address)
    }
    .sheet(isPresented: $showingWorkerPicker) {
      DirectoryPicker<Staff>(
        title: "Выберите профиль медработника",
        label: { worker in "\(worker.fullName) · \(worker.position)" },
        fetch: { offset in
          let page = try await gateway.staff.listStaff(limit: 50, offset: offset)
          return (page.items, page.total)
        }, select: { worker in selectProfile(.worker(worker)) })
    }
    .onChange(of: gateway.address) { selectedProfile = nil }
  }

  private var profileHelpText: String {
    guard let selectedProfile else { return "Выбрать профиль" }
    return "\(selectedProfile.name) · \(selectedProfile.subtitle). Сменить профиль"
  }

  private func selectProfile(_ profile: AppProfile) {
    selectedProfile = profile
    switch profile {
    case .administrator: selectedSection = .staff
    case .worker: selectedSection = .patients
    }
  }
}
