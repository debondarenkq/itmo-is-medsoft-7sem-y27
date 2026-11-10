import SwiftUI

struct FormActions: View {
  let isBusy: Bool
  var saveTitle = "Сохранить"
  let cancel: () -> Void
  let save: () -> Void

  var body: some View {
    HStack {
      if isBusy { ProgressView().controlSize(.small) }
      Spacer()
      Button("Отмена", role: .cancel, action: cancel).keyboardShortcut(.cancelAction)
      Button(saveTitle, action: save).keyboardShortcut(.defaultAction)
    }
    .disabled(isBusy)
    .padding()
  }
}
