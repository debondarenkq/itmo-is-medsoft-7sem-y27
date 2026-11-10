import SwiftUI

struct AdminPasswordRequirements: View {
  let password: String

  var body: some View {
    VStack(alignment: .leading, spacing: 6) {
      requirement("Минимум 8 символов", met: password.count >= 8)
      requirement("Большие и маленькие буквы",
        met: password.contains { $0.isUppercase }
          && password.contains { $0.isLowercase })
      requirement("Хотя бы одна цифра", met: password.contains { $0.isNumber })
    }.font(.caption)
  }

  private func requirement(_ title: String, met: Bool) -> some View {
    Label(title, systemImage: met ? "checkmark.circle.fill" : "circle")
      .foregroundStyle(met ? Color.green : Color.secondary)
  }
}
