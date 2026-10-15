import Foundation
import HISAPI

struct RecordDraft {
  let base: Record
  var diagnoses: [RecordDiagnosis]
  var prescriptions: [Prescription]

  init(_ record: Record) {
    base = record
    diagnoses = record.state.diagnoses
    prescriptions = record.state.prescriptions
  }

  var commands: [Command] {
    var result: [Command] = []
    for original in base.state.diagnoses where !diagnoses.contains(where: { $0.id == original.id })
    {
      result.append(Command(id: original.id, type: .removeDiagnosis))
    }
    for current in diagnoses where !base.state.diagnoses.contains(where: { $0.id == current.id }) {
      result.append(Command(diagnosisId: current.diagnosisId, type: .addDiagnosis))
    }
    for original in base.state.prescriptions {
      guard let current = prescriptions.first(where: { $0.id == original.id }) else {
        result.append(Command(id: original.id, type: .removePrescription))
        continue
      }
      if current.text != original.text {
        result.append(Command(id: original.id, text: current.text, type: .editPrescription))
      } else if current.status != original.status {
        result.append(Command(id: original.id, type: .cancelPrescription))
      }
    }
    for current in prescriptions
    where !base.state.prescriptions.contains(where: { $0.id == current.id }) {
      result.append(Command(text: current.text, type: .addPrescription))
    }
    return result
  }

  var isDirty: Bool { !commands.isEmpty }

  mutating func addDiagnosis(_ diagnosis: Diagnosis) {
    let removed = base.state.diagnoses.first {
      $0.diagnosisId == diagnosis.id && !diagnoses.contains($0)
    }
    diagnoses.append(
      removed
        ?? RecordDiagnosis(
          code: diagnosis.code, diagnosisId: diagnosis.id,
          id: UUID(), name: diagnosis.name))
  }

  func isEdited(_ prescription: Prescription) -> Bool {
    base.state.prescriptions.contains { $0.id == prescription.id && $0.text != prescription.text }
  }

  mutating func cancelPrescription(_ prescription: Prescription) {
    if base.state.prescriptions.contains(where: { $0.id == prescription.id }) {
      if let index = prescriptions.firstIndex(where: { $0.id == prescription.id }) {
        prescriptions[index].status = .cancelled
      }
    } else {
      prescriptions.removeAll { $0.id == prescription.id }
    }
  }

  mutating func setPrescription(_ prescription: Prescription?, text: String) {
    if let prescription, let index = prescriptions.firstIndex(where: { $0.id == prescription.id }) {
      prescriptions[index].text = text
    } else {
      prescriptions.append(Prescription(id: UUID(), status: .active, text: text))
    }
  }
}
