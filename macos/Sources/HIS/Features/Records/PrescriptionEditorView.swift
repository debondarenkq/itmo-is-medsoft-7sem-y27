import HISAPI
import SwiftUI

struct PrescriptionEditorView: View {
  let prescription: Prescription?
  let worker: Staff
  let save: (String) -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var text: String

  init(prescription: Prescription?, worker: Staff, save: @escaping (String) -> Void) {
    self.prescription = prescription
    self.worker = worker
    self.save = save
    _text = State(initialValue: prescription?.text ?? "")
  }

  var body: some View {
    HStack(alignment: .top, spacing: 20) {
      VStack(alignment: .leading, spacing: 16) {
        Button("К ЭМК пациента", systemImage: "chevron.left") { dismiss() }
        Text(prescription == nil ? "Новое назначение" : "Редактирование назначения").font(.title2)
        GroupBox("Содержание *") {
          TextEditor(text: $text).frame(minHeight: 160).padding(4)
        }
        HStack {
          Text("От 3 до 128 символов").foregroundStyle(.secondary)
          Spacer()
          Text("\(text.count) / 128").foregroundStyle(text.count > 128 ? Color.red : Color.secondary)
        }.font(.caption)
        Text("После сохранения черновика примените изменения в ЭМК.").font(.caption).foregroundStyle(.secondary)
        FormActions(isBusy: false, saveTitle: "Сохранить в черновик", cancel: { dismiss() }, save: {
          save(text)
          dismiss()
        })
      }.frame(width: 520)
      WorkerPermissionsView(worker: worker)
    }.padding(24).interactiveDismissDisabled()
  }
}
