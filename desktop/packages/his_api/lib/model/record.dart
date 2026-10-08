//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Record {
  /// Returns a new [Record] instance.
  Record({
    required this.createdAt,
    required this.id,
    required this.patientId,
    required this.state,
    required this.updatedAt,
    required this.version,
  });

  DateTime createdAt;

  String id;

  String patientId;

  State state;

  DateTime updatedAt;

  /// Minimum value: 1
  int version;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Record &&
          other.createdAt == createdAt &&
          other.id == id &&
          other.patientId == patientId &&
          other.state == state &&
          other.updatedAt == updatedAt &&
          other.version == version;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (createdAt.hashCode) +
      (id.hashCode) +
      (patientId.hashCode) +
      (state.hashCode) +
      (updatedAt.hashCode) +
      (version.hashCode);

  @override
  String toString() =>
      'Record[createdAt=$createdAt, id=$id, patientId=$patientId, state=$state, updatedAt=$updatedAt, version=$version]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
    json[r'id'] = this.id;
    json[r'patient_id'] = this.patientId;
    json[r'state'] = this.state;
    json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    json[r'version'] = this.version;
    return json;
  }

  /// Returns a new [Record] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Record? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'created_at'),
            'Required key "Record[created_at]" is missing from JSON.');
        assert(json[r'created_at'] != null,
            'Required key "Record[created_at]" has a null value in JSON.');
        assert(json.containsKey(r'id'),
            'Required key "Record[id]" is missing from JSON.');
        assert(json[r'id'] != null,
            'Required key "Record[id]" has a null value in JSON.');
        assert(json.containsKey(r'patient_id'),
            'Required key "Record[patient_id]" is missing from JSON.');
        assert(json[r'patient_id'] != null,
            'Required key "Record[patient_id]" has a null value in JSON.');
        assert(json.containsKey(r'state'),
            'Required key "Record[state]" is missing from JSON.');
        assert(json[r'state'] != null,
            'Required key "Record[state]" has a null value in JSON.');
        assert(json.containsKey(r'updated_at'),
            'Required key "Record[updated_at]" is missing from JSON.');
        assert(json[r'updated_at'] != null,
            'Required key "Record[updated_at]" has a null value in JSON.');
        assert(json.containsKey(r'version'),
            'Required key "Record[version]" is missing from JSON.');
        assert(json[r'version'] != null,
            'Required key "Record[version]" has a null value in JSON.');
        return true;
      }());

      return Record(
        createdAt: mapDateTime(json, r'created_at', r'')!,
        id: mapValueOfType<String>(json, r'id')!,
        patientId: mapValueOfType<String>(json, r'patient_id')!,
        state: State.fromJson(json[r'state'])!,
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
        version: mapValueOfType<int>(json, r'version')!,
      );
    }
    return null;
  }

  static List<Record> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <Record>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Record.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Record> mapFromJson(dynamic json) {
    final map = <String, Record>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Record.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Record-objects as value to a dart map
  static Map<String, List<Record>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<Record>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Record.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'created_at',
    'id',
    'patient_id',
    'state',
    'updated_at',
    'version',
  };
}
