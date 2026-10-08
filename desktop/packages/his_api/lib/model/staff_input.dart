//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class StaffInput {
  /// Returns a new [StaffInput] instance.
  StaffInput({
    required this.firstName,
    required this.lastName,
    required this.position,
  });

  String firstName;

  String lastName;

  String position;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is StaffInput &&
          other.firstName == firstName &&
          other.lastName == lastName &&
          other.position == position;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (firstName.hashCode) + (lastName.hashCode) + (position.hashCode);

  @override
  String toString() =>
      'StaffInput[firstName=$firstName, lastName=$lastName, position=$position]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'first_name'] = this.firstName;
    json[r'last_name'] = this.lastName;
    json[r'position'] = this.position;
    return json;
  }

  /// Returns a new [StaffInput] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static StaffInput? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'first_name'),
            'Required key "StaffInput[first_name]" is missing from JSON.');
        assert(json[r'first_name'] != null,
            'Required key "StaffInput[first_name]" has a null value in JSON.');
        assert(json.containsKey(r'last_name'),
            'Required key "StaffInput[last_name]" is missing from JSON.');
        assert(json[r'last_name'] != null,
            'Required key "StaffInput[last_name]" has a null value in JSON.');
        assert(json.containsKey(r'position'),
            'Required key "StaffInput[position]" is missing from JSON.');
        assert(json[r'position'] != null,
            'Required key "StaffInput[position]" has a null value in JSON.');
        return true;
      }());

      return StaffInput(
        firstName: mapValueOfType<String>(json, r'first_name')!,
        lastName: mapValueOfType<String>(json, r'last_name')!,
        position: mapValueOfType<String>(json, r'position')!,
      );
    }
    return null;
  }

  static List<StaffInput> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <StaffInput>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = StaffInput.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, StaffInput> mapFromJson(dynamic json) {
    final map = <String, StaffInput>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = StaffInput.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of StaffInput-objects as value to a dart map
  static Map<String, List<StaffInput>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<StaffInput>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = StaffInput.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'first_name',
    'last_name',
    'position',
  };
}
