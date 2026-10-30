import HISAPI
import SwiftUI

struct AdminSetupView: View {
  let session: UserSession
  let secret: String
  @State private var step = 0
  @State private var password = ""
  @State private var confirmation = ""
  @State private var code = ""
  @State private var issuedSession: Session?
  @State private var isBusy = false
  @State private var errorText: String?
  @State private var copied = false

  private var passwordReady: Bool {
    password.count >= 8 && password.utf8.count <= 72 && password == confirmation
      && password.contains { $0.isUppercase }
      && password.contains { $0.isLowercase }
      && password.contains { $0.isNumber }
  }

  var body: some View {
    VStack(alignment: .leading, spacing: 20) {
      Text("Первичная настройка администратора").font(.title2)
      HStack {
        stage("1. Пароль", index: 0)
        Image(systemName: "chevron.right").foregroundStyle(.secondary)
        stage("2. Двухфакторная защита", index: 1)
        Image(systemName: "chevron.right").foregroundStyle(.secondary)
        stage("3. Готово", index: 2)
      }.font(.callout)
      Divider()
      switch step {
      case 0:
        Form {
          PasswordField(title: "Новый пароль", text: $password)
          PasswordField(title: "Повторите пароль", text: $confirmation)
          AdminPasswordRequirements(password: password)
          Text("Не более 72 байт. Окончательную проверку выполнит сервер.")
            .font(.caption).foregroundStyle(.secondary)
        }
      case 1:
        Text("Отсканируйте QR-код приложением-аутентификатором или добавьте секрет вручную.")
        HStack(alignment: .top, spacing: 20) {
          if let image = totpQRCode(secret: secret) {
            Image(nsImage: image).interpolation(.none).resizable().scaledToFit()
              .padding(12).background(.white).frame(width: 200, height: 200)
              .accessibilityLabel("QR-код настройки TOTP для администратора HIS")
          }
          VStack(alignment: .leading, spacing: 12) {
            Text("Секретный ключ").font(.headline)
            Text(secret).font(.system(.body, design: .monospaced)).textSelection(.enabled)
            Button(copied ? "Ключ скопирован" : "Копировать ключ", systemImage: "doc.on.doc") {
              NSPasteboard.general.clearContents()
              NSPasteboard.general.setString(secret, forType: .string)
              copied = true
            }
            Text("TOTP · 6 цифр · 30 секунд · SHA-1").font(.caption).foregroundStyle(.secondary)
            TextField("Одноразовый код", text: $code).textContentType(.oneTimeCode)
          }
        }
      default:
        Label("Администратор настроен", systemImage: "checkmark.circle.fill")
          .font(.title3).foregroundStyle(.green)
        Text("Пароль сохранён, одноразовый код проверен. Для последующих входов нужны пароль и текущий код аутентификатора.")
      }
      ErrorNotice(message: errorText)
      HStack {
        if step == 1 { Button("Назад") { step = 0; errorText = nil } }
        Spacer()
        if isBusy { ProgressView().controlSize(.small) }
        switch step {
        case 0: Button("Продолжить") { step = 1; errorText = nil }.keyboardShortcut(.defaultAction).disabled(!passwordReady)
        case 1: Button("Проверить код") { Task { await configure() } }.keyboardShortcut(.defaultAction)
        default: Button("Открыть администрирование") { Task { await enter() } }.keyboardShortcut(.defaultAction)
        }
      }
    }.padding(28).frame(width: 660).disabled(isBusy)
  }

  private func stage(_ title: String, index: Int) -> some View {
    Text(title).fontWeight(step == index ? .semibold : .regular)
      .foregroundStyle(step == index ? Color.accentColor : Color.secondary)
  }

  private func configure() async {
    isBusy = true
    defer { isBusy = false; code = "" }
    do {
      issuedSession = try await Gateway(address: session.gateway.address).security.setupAdmin(
        setupInput: SetupInput(confirmation: confirmation, password: password, secret: secret, totp: code))
      password = ""; confirmation = ""; errorText = nil; step = 2
    } catch { errorText = errorMessage(error) }
  }

  private func enter() async {
    guard let issuedSession else { return }
    isBusy = true
    defer { isBusy = false }
    do { try await session.enter(issuedSession) }
    catch { errorText = errorMessage(error) }
  }
}
