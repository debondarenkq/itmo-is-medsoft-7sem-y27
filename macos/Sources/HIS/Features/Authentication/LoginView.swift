import HISAPI
import SwiftUI

struct LoginView: View {
  let session: UserSession
  @State private var username = ""
  @State private var password = ""
  @State private var isBusy = false
  @State private var errorText: String?

  var body: some View {
    VStack(spacing: 20) {
      Label("HIS — Рабочее место медработника", systemImage: "stethoscope")
        .font(.title2)
      Form {
        TextField("Логин", text: $username).textContentType(.username)
        PasswordField(title: "Пароль", text: $password)
      }
      ErrorNotice(message: errorText)
      Button("Войти") { Task { await login() } }.keyboardShortcut(.defaultAction)
        .disabled(isBusy || username.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty || password.isEmpty)
      if isBusy { ProgressView() }
    }
    .padding(28).frame(width: 520)
  }

  private func login() async {
    isBusy = true
    defer { isBusy = false; password = "" }
    do {
      let result = try await Gateway(address: session.gateway.address).security.login(
        loginInput: LoginInput(password: password, username: username))
      try await session.enter(result)
    } catch { errorText = errorMessage(error) }
  }
}
