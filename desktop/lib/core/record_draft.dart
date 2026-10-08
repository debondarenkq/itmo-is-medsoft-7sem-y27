import 'package:flutter/foundation.dart';
import 'models.dart';

class RecordDraft extends ChangeNotifier {
  RecordDraft(this.base)
    : diagnoses = List.of(base.diagnoses),
      prescriptions = List.of(base.prescriptions);
  final MedicalRecord base;
  final List<RecordDiagnosis> diagnoses;
  final List<Prescription> prescriptions;
  int _nextId = 0;
  String _localId() => 'draft-${_nextId++}';
  bool get dirty => commands.isNotEmpty;
  bool containsDiagnosis(String id) =>
      diagnoses.any((d) => d.diagnosisId == id);

  void addDiagnosis(Diagnosis value) {
    if (containsDiagnosis(value.id)) return;
    // Undo a pending removal without introducing a duplicate historical entry.
    final original = base.diagnoses
        .where((d) => d.diagnosisId == value.id)
        .firstOrNull;
    diagnoses.add(
      original ??
          RecordDiagnosis(
            id: _localId(),
            diagnosisId: value.id,
            code: value.code,
            name: value.name,
          ),
    );
    notifyListeners();
  }

  void removeDiagnosis(String id) {
    diagnoses.removeWhere((d) => d.id == id);
    notifyListeners();
  }

  void undoDiagnosisRemoval(RecordDiagnosis value) {
    if (!containsDiagnosis(value.diagnosisId)) {
      diagnoses.add(value);
      notifyListeners();
    }
  }

  void addPrescription(String text) {
    prescriptions.add(
      Prescription(id: _localId(), text: text.trim(), status: 'active'),
    );
    notifyListeners();
  }

  void editPrescription(String id, String text) {
    final index = prescriptions.indexWhere((p) => p.id == id);
    if (index < 0 || !prescriptions[index].active) return;
    prescriptions[index] = prescriptions[index].copyWith(text: text.trim());
    notifyListeners();
  }

  bool isNewPrescription(String id) =>
      !base.prescriptions.any((p) => p.id == id);
  bool isEditedPrescription(String id) {
    final current = prescriptions.where((p) => p.id == id).firstOrNull;
    final original = base.prescriptions.where((p) => p.id == id).firstOrNull;
    return current != null && original != null && current.text != original.text;
  }

  void cancelPrescription(String id) {
    final index = prescriptions.indexWhere((p) => p.id == id);
    if (index < 0 || isEditedPrescription(id)) return;
    if (isNewPrescription(id)) {
      prescriptions.removeAt(index);
    } else {
      prescriptions[index] = prescriptions[index].copyWith(status: 'cancelled');
    }
    notifyListeners();
  }

  List<Json> get commands {
    final result = <Json>[];
    for (final old in base.diagnoses) {
      if (!diagnoses.any((d) => d.id == old.id)) {
        result.add({'type': 'remove_diagnosis', 'id': old.id});
      }
    }
    for (final current in diagnoses) {
      if (!base.diagnoses.any((d) => d.id == current.id)) {
        result.add({
          'type': 'add_diagnosis',
          'diagnosis_id': current.diagnosisId,
        });
      }
    }
    for (final old in base.prescriptions) {
      final current = prescriptions.where((p) => p.id == old.id).firstOrNull;
      if (current == null) {
        result.add({'type': 'remove_prescription', 'id': old.id});
      } else if (current.text != old.text) {
        result.add({
          'type': 'edit_prescription',
          'id': old.id,
          'text': current.text,
        });
      } else if (current.status != old.status) {
        result.add({'type': 'cancel_prescription', 'id': old.id});
      }
    }
    for (final current in prescriptions) {
      if (isNewPrescription(current.id)) {
        result.add({'type': 'add_prescription', 'text': current.text});
      }
    }
    return result;
  }
}
