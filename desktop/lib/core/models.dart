import 'package:his_api/api.dart' as contract;

typedef Json = Map<String, dynamic>;

class Staff {
  const Staff({
    required this.id,
    required this.firstName,
    required this.lastName,
    required this.position,
    this.deletedAt,
  });
  factory Staff.fromContract(contract.Staff value) => Staff(
    id: value.id,
    firstName: value.firstName,
    lastName: value.lastName,
    position: value.position,
    deletedAt: value.deletedAt,
  );
  factory Staff.fromActor(contract.Actor value) => Staff(
    id: value.id,
    firstName: value.firstName,
    lastName: value.lastName,
    position: value.position,
  );
  final String id, firstName, lastName, position;
  final DateTime? deletedAt;
  String get name => '$lastName $firstName';
  bool get active => deletedAt == null;
}

class Diagnosis {
  const Diagnosis({
    required this.id,
    required this.code,
    required this.name,
    this.deletedAt,
  });
  factory Diagnosis.fromContract(contract.Diagnosis value) => Diagnosis(
    id: value.id,
    code: value.code,
    name: value.name,
    deletedAt: value.deletedAt,
  );
  final String id, code, name;
  final DateTime? deletedAt;
  bool get active => deletedAt == null;
}

class Patient {
  const Patient({
    required this.id,
    required this.firstName,
    required this.lastName,
    this.middleName,
    required this.birthDate,
    required this.sex,
    required this.comment,
    required this.hasRecord,
  });
  factory Patient.fromContract(contract.Patient value) => Patient(
    id: value.id,
    firstName: value.firstName,
    lastName: value.lastName,
    middleName: value.middleName,
    birthDate: value.birthDate,
    sex: value.administrativeSex.toJson(),
    comment: value.comment,
    hasRecord: value.hasRecord,
  );
  final String id, firstName, lastName, sex, comment;
  final String? middleName;
  final DateTime birthDate;
  final bool hasRecord;
  String get name =>
      [lastName, firstName, if (middleName != null) middleName!].join(' ');
  String get sexLabel => switch (sex) {
    'M' => 'М',
    'F' => 'Ж',
    _ => 'НУ',
  };
}

class RecordDiagnosis {
  const RecordDiagnosis({
    required this.id,
    required this.diagnosisId,
    required this.code,
    required this.name,
  });
  factory RecordDiagnosis.fromContract(contract.RecordDiagnosis value) =>
      RecordDiagnosis(
        id: value.id,
        diagnosisId: value.diagnosisId,
        code: value.code,
        name: value.name,
      );
  factory RecordDiagnosis.fromSnapshot(Json value) =>
      RecordDiagnosis.fromContract(contract.RecordDiagnosis.fromJson(value)!);
  final String id, diagnosisId, code, name;
}

class Prescription {
  const Prescription({
    required this.id,
    required this.text,
    required this.status,
  });
  factory Prescription.fromContract(contract.Prescription value) =>
      Prescription(
        id: value.id,
        text: value.text,
        status: value.status.toJson(),
      );
  final String id, text, status;
  bool get active => status == 'active';
  Prescription copyWith({String? text, String? status}) => Prescription(
    id: id,
    text: text ?? this.text,
    status: status ?? this.status,
  );
}

class MedicalRecord {
  const MedicalRecord({
    required this.id,
    required this.patientId,
    required this.version,
    required this.createdAt,
    required this.updatedAt,
    required this.diagnoses,
    required this.prescriptions,
  });
  factory MedicalRecord.fromContract(contract.Record value) => MedicalRecord(
    id: value.id,
    patientId: value.patientId,
    version: value.version,
    createdAt: value.createdAt,
    updatedAt: value.updatedAt,
    diagnoses: value.state.diagnoses.map(RecordDiagnosis.fromContract).toList(),
    prescriptions: value.state.prescriptions
        .map(Prescription.fromContract)
        .toList(),
  );
  final String id, patientId;
  final int version;
  final DateTime createdAt, updatedAt;
  final List<RecordDiagnosis> diagnoses;
  final List<Prescription> prescriptions;
}

class RecordEvent {
  const RecordEvent({
    required this.sequence,
    required this.version,
    required this.at,
    required this.actor,
    required this.type,
    required this.entityId,
    this.before,
    this.after,
  });
  factory RecordEvent.fromContract(contract.Event value) => RecordEvent(
    sequence: value.sequence,
    version: value.version,
    at: value.occurredAt,
    actor: Staff.fromActor(value.actor),
    type: value.type.toJson(),
    entityId: value.entityId,
    before: value.before?.map((key, value) => MapEntry(key, value)),
    after: value.after?.map((key, value) => MapEntry(key, value)),
  );
  final int sequence, version;
  final DateTime at;
  final Staff actor;
  final String type, entityId;
  final Json? before, after;
  String get label => switch (type) {
    'record_created' => 'Создана медицинская карта',
    'diagnosis_added' => 'Добавлен диагноз',
    'diagnosis_removed' => 'Удалён диагноз',
    'prescription_added' => 'Добавлено назначение',
    'prescription_removed' => 'Удалено назначение',
    'prescription_cancelled' => 'Отменено назначение',
    _ => 'Изменение карты',
  };
  String get description =>
      (after?['text'] ??
              after?['name'] ??
              before?['text'] ??
              before?['name'] ??
              '')
          as String;
}

class PageResult<T> {
  const PageResult(this.items, this.total);
  final List<T> items;
  final int total;
}

String calendarDate(DateTime value) =>
    '${value.year.toString().padLeft(4, '0')}-${value.month.toString().padLeft(2, '0')}-${value.day.toString().padLeft(2, '0')}';
String displayDate(DateTime value) =>
    '${value.day.toString().padLeft(2, '0')}.${value.month.toString().padLeft(2, '0')}.${value.year}';
String displayTime(DateTime value) {
  final local = value.toLocal();
  return '${displayDate(local)} · ${local.hour.toString().padLeft(2, '0')}:${local.minute.toString().padLeft(2, '0')}:${local.second.toString().padLeft(2, '0')}';
}
