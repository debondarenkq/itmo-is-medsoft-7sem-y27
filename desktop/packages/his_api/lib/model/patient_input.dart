//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class PatientInput {
  /// Returns a new [PatientInput] instance.
  PatientInput({
    required this.administrativeSex,
    required this.birthDate,
    this.comment = const Optional.absent(),
    required this.firstName,
    required this.lastName,
    this.middleName = const Optional.absent(),
  });

  AdministrativeSex administrativeSex;

  DateTime birthDate;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  Optional<String?> comment;

  String firstName;

  String lastName;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  Optional<String?> middleName;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is PatientInput &&
          other.administrativeSex == administrativeSex &&
          other.birthDate == birthDate &&
          other.comment == comment &&
          other.firstName == firstName &&
          other.lastName == lastName &&
          other.middleName == middleName;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (administrativeSex.hashCode) +
      (birthDate.hashCode) +
      (comment == null ? 0 : comment!.hashCode) +
      (firstName.hashCode) +
      (lastName.hashCode) +
      (middleName == null ? 0 : middleName!.hashCode);

  @override
  String toString() =>
      'PatientInput[administrativeSex=$administrativeSex, birthDate=$birthDate, comment=$comment, firstName=$firstName, lastName=$lastName, middleName=$middleName]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'administrative_sex'] = this.administrativeSex;
    json[r'birth_date'] = _dateFormatter.format(this.birthDate);
    if (this.comment.isPresent) {
      final value = this.comment.value;
      json[r'comment'] = value;
    }
    json[r'first_name'] = this.firstName;
    json[r'last_name'] = this.lastName;
    if (this.middleName.isPresent) {
      final value = this.middleName.value;
      json[r'middle_name'] = value;
    }
    return json;
  }

  /// Returns a new [PatientInput] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static PatientInput? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'administrative_sex'),
            'Required key "PatientInput[administrative_sex]" is missing from JSON.');
        assert(json[r'administrative_sex'] != null,
            'Required key "PatientInput[administrative_sex]" has a null value in JSON.');
        assert(json.containsKey(r'birth_date'),
            'Required key "PatientInput[birth_date]" is missing from JSON.');
        assert(json[r'birth_date'] != null,
            'Required key "PatientInput[birth_date]" has a null value in JSON.');
        assert(json.containsKey(r'first_name'),
            'Required key "PatientInput[first_name]" is missing from JSON.');
        assert(json[r'first_name'] != null,
            'Required key "PatientInput[first_name]" has a null value in JSON.');
        assert(json.containsKey(r'last_name'),
            'Required key "PatientInput[last_name]" is missing from JSON.');
        assert(json[r'last_name'] != null,
            'Required key "PatientInput[last_name]" has a null value in JSON.');
        return true;
      }());

      return PatientInput(
        administrativeSex:
            AdministrativeSex.fromJson(json[r'administrative_sex'])!,
        birthDate: mapDateTime(json, r'birth_date', r'')!,
        comment: json.containsKey(r'comment')
            ? Optional.present(mapValueOfType<String>(json, r'comment'))
            : const Optional.absent(),
        firstName: mapValueOfType<String>(json, r'first_name')!,
        lastName: mapValueOfType<String>(json, r'last_name')!,
        middleName: json.containsKey(r'middle_name')
            ? Optional.present(mapValueOfType<String>(json, r'middle_name'))
            : const Optional.absent(),
      );
    }
    return null;
  }

  static List<PatientInput> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <PatientInput>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = PatientInput.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, PatientInput> mapFromJson(dynamic json) {
    final map = <String, PatientInput>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = PatientInput.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of PatientInput-objects as value to a dart map
  static Map<String, List<PatientInput>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<PatientInput>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = PatientInput.listFromJson(
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
    'first_name',
    'last_name',
  };
}
