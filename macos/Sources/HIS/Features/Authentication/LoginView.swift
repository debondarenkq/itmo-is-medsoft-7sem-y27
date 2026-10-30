import HISAPI
import SwiftUI

struct LoginView: View {
  let session: UserSession
  @State private var administrator = false
  @State private var username = ""
  @State private var password = ""
  @State private var code = ""
  @State private var setup: SetupStatus?
  @State private var isBusy = false
  @State private var errorText: String?

  var body: some View {
    VStack(spacing: 20) {
      Label("Госпитальная информационная система", systemImage: "cross.case")
        .font(.title2)
      Picker("Вход", selection: $administrator) {
        Text("Медработник").tag(false)
        Text("Администратор").tag(true)
      }.pickerStyle(.segmented)
      if administrator, let setup, !setup.configured, let secret = setup.secret {
        AdminSetupView(session: session, secret: secret)
      } else {
        Form {
          if administrator { LabeledContent("Логин", value: "admin") }
          else { TextField("Логин", text: $username).textContentType(.username) }
          PasswordField(title: "Пароль", text: $password)
          if administrator {
            TextField("Код аутентификатора", text: $code).textContentType(.oneTimeCode)
          }
        }
        ErrorNotice(message: errorText)
        Button("Войти") { Task { await login() } }.keyboardShortcut(.defaultAction)
          .disabled(isBusy || setup == nil)
        if isBusy { ProgressView() }
        if setup == nil { Button("Повторить подключение") { Task { await loadSetup() } } }
      }
    }
    .padding(28).frame(width: administrator && setup?.configured == false ? 740 : 520)
    .task { await loadSetup() }
    .onChange(of: administrator) { password = ""; code = ""; errorText = nil }
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
        loginInput: LoginInput(password: password, totp: administrator ? code : nil,
          username: administrator ? "admin" : username))
      try await session.enter(result)
    } catch { errorText = errorMessage(error) }
  }
}
