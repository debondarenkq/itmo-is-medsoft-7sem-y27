import SwiftUI

struct ErrorNotice: View {
  let message: String?

  var body: some View {
    if let message {
      Label(message, systemImage: "exclamationmark.triangle")
        .foregroundStyle(.red)
        .textSelection(.enabled)
        .padding()
        .frame(maxWidth: .infinity, alignment: .leading)
    }
  }
}
