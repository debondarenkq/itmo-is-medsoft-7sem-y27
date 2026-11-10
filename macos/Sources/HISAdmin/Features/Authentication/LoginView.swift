import HISAPI
import SwiftUI

struct LoginView: View {
  let session: UserSession
  @State private var password = ""
  @State private var code = ""
  @State private var setup: SetupStatus?
  @State private var isBusy = false
  @State private var errorText: String?

  var body: some View {
    VStack(spacing: 20) {
      Label("HIS — Администрирование", systemImage: "person.badge.key")
        .font(.title2)
      if let setup, !setup.configured, let secret = setup.secret {
        AdminSetupView(session: session, secret: secret, provisioningURI: setup.provisioningUri)
      } else {
        Form {
          LabeledContent("Логин", value: "admin")
          PasswordField(title: "Пароль", text: $password)
          TextField("Код аутентификатора", text: $code).textContentType(.oneTimeCode)
        }
        ErrorNotice(message: errorText)
        Button("Войти") { Task { await login() } }.keyboardShortcut(.defaultAction)
          .disabled(isBusy || setup == nil)
        if isBusy { ProgressView() }
        if setup == nil { Button("Повторить подключение") { Task { await loadSetup() } } }
      }
    }
    .padding(28).frame(width: setup?.configured == false ? 740 : 520)
    .task { await loadSetup() }
  }

  private func loadSetup() async {
    do { setup = try await Gateway(address: session.gateway.address).security.getSetup(); errorText = nil }
    catch { errorText = errorMessage(error) }
  }

  private func login() async {
    isBusy = true
    defer { isBusy = false; password = ""; code = "" }
    do {
      let result = try await Gateway(address: session.gateway.address).security.login(
        loginInput: LoginInput(password: password, totp: code,
          username: "admin"))
      try await session.enter(result)
    } catch { errorText = errorMessage(error) }
  }
}
