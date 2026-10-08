//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Patient {
  /// Returns a new [Patient] instance.
  Patient({
    required this.administrativeSex,
    required this.birthDate,
    required this.comment,
    required this.createdAt,
    required this.firstName,
    required this.hasRecord,
    required this.id,
    required this.lastName,
    required this.middleName,
    required this.updatedAt,
  });

  AdministrativeSex administrativeSex;

  DateTime birthDate;

  String comment;

  DateTime createdAt;

  String firstName;

  bool hasRecord;

  String id;

  String lastName;

  String? middleName;

  DateTime updatedAt;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Patient &&
          other.administrativeSex == administrativeSex &&
          other.birthDate == birthDate &&
          other.comment == comment &&
          other.createdAt == createdAt &&
          other.firstName == firstName &&
          other.hasRecord == hasRecord &&
          other.id == id &&
          other.lastName == lastName &&
          other.middleName == middleName &&
          other.updatedAt == updatedAt;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (administrativeSex.hashCode) +
      (birthDate.hashCode) +
      (comment.hashCode) +
      (createdAt.hashCode) +
      (firstName.hashCode) +
      (hasRecord.hashCode) +
      (id.hashCode) +
      (lastName.hashCode) +
      (middleName == null ? 0 : middleName!.hashCode) +
      (updatedAt.hashCode);

  @override
  String toString() =>
      'Patient[administrativeSex=$administrativeSex, birthDate=$birthDate, comment=$comment, createdAt=$createdAt, firstName=$firstName, hasRecord=$hasRecord, id=$id, lastName=$lastName, middleName=$middleName, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'administrative_sex'] = this.administrativeSex;
    json[r'birth_date'] = _dateFormatter.format(this.birthDate);
    json[r'comment'] = this.comment;
    json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
    json[r'first_name'] = this.firstName;
    json[r'has_record'] = this.hasRecord;
    json[r'id'] = this.id;
    json[r'last_name'] = this.lastName;
    if (this.middleName != null) {
      json[r'middle_name'] = this.middleName;
    } else {
      json[r'middle_name'] = null;
    }
    json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [Patient] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Patient? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'administrative_sex'),
            'Required key "Patient[administrative_sex]" is missing from JSON.');
        assert(json[r'administrative_sex'] != null,
            'Required key "Patient[administrative_sex]" has a null value in JSON.');
        assert(json.containsKey(r'birth_date'),
            'Required key "Patient[birth_date]" is missing from JSON.');
        assert(json[r'birth_date'] != null,
            'Required key "Patient[birth_date]" has a null value in JSON.');
        assert(json.containsKey(r'comment'),
            'Required key "Patient[comment]" is missing from JSON.');
        assert(json[r'comment'] != null,
            'Required key "Patient[comment]" has a null value in JSON.');
        assert(json.containsKey(r'created_at'),
            'Required key "Patient[created_at]" is missing from JSON.');
        assert(json[r'created_at'] != null,
            'Required key "Patient[created_at]" has a null value in JSON.');
        assert(json.containsKey(r'first_name'),
            'Required key "Patient[first_name]" is missing from JSON.');
        assert(json[r'first_name'] != null,
            'Required key "Patient[first_name]" has a null value in JSON.');
        assert(json.containsKey(r'has_record'),
            'Required key "Patient[has_record]" is missing from JSON.');
        assert(json[r'has_record'] != null,
            'Required key "Patient[has_record]" has a null value in JSON.');
        assert(json.containsKey(r'id'),
            'Required key "Patient[id]" is missing from JSON.');
        assert(json[r'id'] != null,
            'Required key "Patient[id]" has a null value in JSON.');
        assert(json.containsKey(r'last_name'),
            'Required key "Patient[last_name]" is missing from JSON.');
        assert(json[r'last_name'] != null,
            'Required key "Patient[last_name]" has a null value in JSON.');
        assert(json.containsKey(r'middle_name'),
            'Required key "Patient[middle_name]" is missing from JSON.');
        assert(json.containsKey(r'updated_at'),
            'Required key "Patient[updated_at]" is missing from JSON.');
        assert(json[r'updated_at'] != null,
            'Required key "Patient[updated_at]" has a null value in JSON.');
        return true;
      }());

      return Patient(
        administrativeSex:
            AdministrativeSex.fromJson(json[r'administrative_sex'])!,
        birthDate: mapDateTime(json, r'birth_date', r'')!,
        comment: mapValueOfType<String>(json, r'comment')!,
        createdAt: mapDateTime(json, r'created_at', r'')!,
        firstName: mapValueOfType<String>(json, r'first_name')!,
        hasRecord: mapValueOfType<bool>(json, r'has_record')!,
        id: mapValueOfType<String>(json, r'id')!,
        lastName: mapValueOfType<String>(json, r'last_name')!,
        middleName: mapValueOfType<String>(json, r'middle_name'),
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
      );
    }
    return null;
  }

  static List<Patient> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <Patient>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Patient.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Patient> mapFromJson(dynamic json) {
    final map = <String, Patient>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Patient.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Patient-objects as value to a dart map
  static Map<String, List<Patient>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<Patient>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Patient.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'administrative_sex',
    'birth_date',
    'comment',
    'created_at',
    'first_name',
    'has_record',
    'id',
    'last_name',
    'middle_name',
    'updated_at',
  };
}
