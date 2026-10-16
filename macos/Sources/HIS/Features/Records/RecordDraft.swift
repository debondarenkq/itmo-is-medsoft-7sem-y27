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
    for original in base.state.diagnoses {
      if let current = diagnoses.first(where: { $0.id == original.id }) {
        if current.diagnosisId != original.diagnosisId {
          result.append(
            Command(diagnosisId: current.diagnosisId, id: original.id, type: .editDiagnosis))
        }
      } else {
        result.append(Command(id: original.id, type: .removeDiagnosis))
      }
    }
    for diagnosis in diagnoses {
      let alreadySaved = base.state.diagnoses.contains { original in original.id == diagnosis.id }
      if !alreadySaved {
        result.append(Command(diagnosisId: diagnosis.diagnosisId, type: .addDiagnosis))
      }
    }
    for original in base.state.prescriptions {
      guard
        let current = prescriptions.first(where: { prescription in prescription.id == original.id })
      else {
        result.append(Command(id: original.id, type: .removePrescription))
        continue
      }
      if current.text != original.text {
        result.append(Command(id: original.id, text: current.text, type: .editPrescription))
      } else if current.status != original.status {
        result.append(Command(id: original.id, type: .cancelPrescription))
      }
    }
    for prescription in prescriptions {
      let alreadySaved = base.state.prescriptions.contains { original in
        original.id == prescription.id
      }
      if !alreadySaved {
        result.append(Command(text: prescription.text, type: .addPrescription))
      }
    }
    return result
  }

  var isDirty: Bool { !commands.isEmpty }

  mutating func setDiagnosis(_ original: RecordDiagnosis?, diagnosis: Diagnosis) {
    if let original, let index = diagnoses.firstIndex(where: { $0.id == original.id }) {
      diagnoses[index] = RecordDiagnosis(
        code: diagnosis.code, diagnosisId: diagnosis.id, id: original.id, name: diagnosis.name)
    } else {
      addDiagnosis(diagnosis)
    }
  }

  mutating func addDiagnosis(_ diagnosis: Diagnosis) {
    let removedDiagnosis = base.state.diagnoses.first { original in
      let stillPresent = diagnoses.contains(original)
      return original.diagnosisId == diagnosis.id && !stillPresent
    }
    if let removedDiagnosis {
      diagnoses.append(removedDiagnosis)
    } else {
      diagnoses.append(
        RecordDiagnosis(
          code: diagnosis.code, diagnosisId: diagnosis.id,
          id: UUID(), name: diagnosis.name))
    }
  }

  mutating func removeDiagnosis(id: UUID) {
    diagnoses.removeAll { diagnosis in diagnosis.id == id }
  }

  mutating func restoreDiagnosis(_ diagnosis: RecordDiagnosis) {
    let belongsToCurrentVersion = base.state.diagnoses.contains { original in
      original.id == diagnosis.id
    }
    if belongsToCurrentVersion {
      diagnoses.append(diagnosis)
    } else {
      diagnoses.append(
        RecordDiagnosis(
          code: diagnosis.code, diagnosisId: diagnosis.diagnosisId,
          id: UUID(), name: diagnosis.name))
    }
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
