typedef Json = Map<String, dynamic>;

class Staff {
  const Staff({
    required this.id,
    required this.firstName,
    required this.lastName,
    required this.position,
    this.deletedAt,
  });
  factory Staff.fromJson(Json json) => Staff(
    id: json['id'] as String,
    firstName: json['first_name'] as String,
    lastName: json['last_name'] as String,
    position: json['position'] as String,
    deletedAt: json['deleted_at'] == null
        ? null
        : DateTime.parse(json['deleted_at'] as String),
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
  factory Diagnosis.fromJson(Json json) => Diagnosis(
    id: json['id'] as String,
    code: json['code'] as String,
    name: json['name'] as String,
    deletedAt: json['deleted_at'] == null
        ? null
        : DateTime.parse(json['deleted_at'] as String),
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
  factory Patient.fromJson(Json json) => Patient(
    id: json['id'] as String,
    firstName: json['first_name'] as String,
    lastName: json['last_name'] as String,
    middleName: json['middle_name'] as String?,
    birthDate: DateTime.parse(json['birth_date'] as String),
    sex: json['administrative_sex'] as String,
    comment: json['comment'] as String,
    hasRecord: json['has_record'] as bool,
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
  factory RecordDiagnosis.fromJson(Json json) => RecordDiagnosis(
    id: json['id'] as String,
    diagnosisId: json['diagnosis_id'] as String,
    code: json['code'] as String,
    name: json['name'] as String,
  );
  final String id, diagnosisId, code, name;
}

class Prescription {
  const Prescription({
    required this.id,
    required this.text,
    required this.status,
  });
  factory Prescription.fromJson(Json json) => Prescription(
    id: json['id'] as String,
    text: json['text'] as String,
    status: json['status'] as String,
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
  factory MedicalRecord.fromJson(Json json) {
    final state = json['state'] as Json;
    return MedicalRecord(
      id: json['id'] as String,
      patientId: json['patient_id'] as String,
      version: json['version'] as int,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
      diagnoses: (state['diagnoses'] as List)
          .map((item) => RecordDiagnosis.fromJson(item as Json))
          .toList(),
      prescriptions: (state['prescriptions'] as List)
          .map((item) => Prescription.fromJson(item as Json))
          .toList(),
    );
  }
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
  factory RecordEvent.fromJson(Json json) => RecordEvent(
    sequence: json['sequence'] as int,
    version: json['version'] as int,
    at: DateTime.parse(json['occurred_at'] as String),
    actor: Staff.fromJson(json['actor'] as Json),
    type: json['type'] as String,
    entityId: json['entity_id'] as String,
    before: json['before'] as Json?,
    after: json['after'] as Json?,
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
