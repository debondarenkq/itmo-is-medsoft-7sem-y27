import HISAPI
import SwiftUI

struct DiagnosisEditorView: View {
  let gateway: Gateway
  let diagnosis: Diagnosis?
  let saved: () -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var code: String
  @State private var name: String
  @State private var isBusy = false
  @State private var errorText: String?

  init(gateway: Gateway, diagnosis: Diagnosis?, saved: @escaping () -> Void) {
    self.gateway = gateway
    self.diagnosis = diagnosis
    self.saved = saved
    _code = State(initialValue: diagnosis?.code ?? "")
    _name = State(initialValue: diagnosis?.name ?? "")
  }

  var body: some View {
    VStack(spacing: 0) {
      Form {
        Section(diagnosis == nil ? "Новый диагноз" : "Диагноз") {
          TextField("Код", text: $code)
          TextField("Наименование", text: $name)
        }
      }.formStyle(.grouped).disabled(isBusy)
      ErrorNotice(message: errorText)
      FormActions(isBusy: isBusy, cancel: { dismiss() }, save: { Task { await save() } })
    }
    .frame(width: 480)
    .interactiveDismissDisabled()
  }

  private func save() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let input = DiagnosisInput(code: code, name: name)
      if let diagnosis {
        _ = try await gateway.diagnoses.updateDiagnosis(id: diagnosis.id, diagnosisInput: input)
      } else {
        _ = try await gateway.diagnoses.createDiagnosis(diagnosisInput: input)
      }
      saved()
      dismiss()
    } catch {
      errorText = errorMessage(error)
    }
  }
}
