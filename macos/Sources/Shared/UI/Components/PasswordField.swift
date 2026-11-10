import SwiftUI

struct PasswordField: View {
  let title: String
  @Binding var text: String
  @State private var revealed = false

  var body: some View {
    HStack {
      if revealed { TextField(title, text: $text) }
      else { SecureField(title, text: $text) }
      Button(revealed ? "Скрыть пароль" : "Показать пароль",
        systemImage: revealed ? "eye.slash" : "eye") { revealed.toggle() }
        .labelStyle(.iconOnly).buttonStyle(.borderless)
    }
  }
}
