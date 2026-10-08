import 'package:flutter_test/flutter_test.dart';
import 'package:his_desktop/core/models.dart';
import 'package:his_desktop/core/record_draft.dart';

MedicalRecord record() => MedicalRecord(
  id: 'record',
  patientId: 'patient',
  version: 1,
  createdAt: DateTime.utc(2026),
  updatedAt: DateTime.utc(2026),
  diagnoses: const [
    RecordDiagnosis(
      id: 'entry',
      diagnosisId: 'diagnosis',
      code: 'TEST01',
      name: 'Диагноз',
    ),
  ],
  prescriptions: const [
    Prescription(id: 'rx', text: 'Пить воду', status: 'active'),
  ],
);
void main() {
  test('undo diagnosis removal returns to unchanged draft', () {
    final draft = RecordDraft(record());
    draft.removeDiagnosis('entry');
    expect(draft.commands.single, {'type': 'remove_diagnosis', 'id': 'entry'});
    draft.addDiagnosis(
      const Diagnosis(id: 'diagnosis', code: 'TEST01', name: 'Новое название'),
    );
    expect(draft.dirty, isFalse);
    expect(draft.diagnoses.single.name, 'Диагноз');
  });
  test('new prescriptions can be edited and discarded before saving', () {
    final draft = RecordDraft(record());
    draft.addPrescription('Первый вариант');
    final id = draft.prescriptions.last.id;
    draft.editPrescription(id, 'Второй вариант');
    expect(draft.commands.single, {
      'type': 'add_prescription',
      'text': 'Второй вариант',
    });
    draft.cancelPrescription(id);
    expect(draft.dirty, isFalse);
  });
  test(
    'edit and cancellation of persisted prescriptions produce correct commands',
    () {
      final edited = RecordDraft(record());
      edited.editPrescription('rx', 'Пить больше воды');
      expect(edited.commands.single, {
        'type': 'edit_prescription',
        'id': 'rx',
        'text': 'Пить больше воды',
      });
      edited.cancelPrescription('rx');
      expect(edited.commands.single['type'], 'edit_prescription');
      final cancelled = RecordDraft(record());
      cancelled.cancelPrescription('rx');
      expect(cancelled.commands.single, {
        'type': 'cancel_prescription',
        'id': 'rx',
      });
      expect(cancelled.base.prescriptions.single.status, 'active');
    },
  );
  test('adding then removing a diagnosis produces no server command', () {
    final draft = RecordDraft(record());
    draft.addDiagnosis(
      const Diagnosis(id: 'other', code: 'TEST02', name: 'Второй диагноз'),
    );
    draft.removeDiagnosis(draft.diagnoses.last.id);
    expect(draft.dirty, isFalse);
  });
}
