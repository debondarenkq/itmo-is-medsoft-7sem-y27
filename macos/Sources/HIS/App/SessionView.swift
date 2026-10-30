import SwiftUI

struct SessionView: View {
  @State private var session: UserSession

  init(address: String) { _session = State(initialValue: UserSession(address: address)) }

  var body: some View {
    Group {
      if let profile = session.profile {
        MainView(gateway: session.gateway, profile: profile) { Task { await session.logout() } }
          .task {
            while !Task.isCancelled {
              try? await Task.sleep(for: .seconds(30))
              if !Task.isCancelled { await session.check() }
            }
          }
      } else {
        LoginView(session: session)
      }
    }
  }
}
