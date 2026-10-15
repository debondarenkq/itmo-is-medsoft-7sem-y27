import HISAPI
import SwiftUI

struct PrescriptionEditorView: View {
  let prescription: Prescription?
  let save: (String) -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var text: String

  init(prescription: Prescription?, save: @escaping (String) -> Void) {
    self.prescription = prescription
    self.save = save
    _text = State(initialValue: prescription?.text ?? "")
  }

  var body: some View {
    VStack(spacing: 0) {
      Form {
        TextField("Назначение", text: $text, axis: .vertical).lineLimit(3...8)
      }.formStyle(.grouped)
      FormActions(
        isBusy: false, saveTitle: "Добавить в черновик", cancel: { dismiss() },
        save: {
          save(text)
          dismiss()
        })
    }
    .frame(width: 520)
    .interactiveDismissDisabled()
  }
}
