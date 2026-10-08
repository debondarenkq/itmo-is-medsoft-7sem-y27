//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Staff {
  /// Returns a new [Staff] instance.
  Staff({
    required this.createdAt,
    required this.deletedAt,
    required this.firstName,
    required this.id,
    required this.lastName,
    required this.position,
    required this.updatedAt,
  });

  DateTime createdAt;

  DateTime? deletedAt;

  String firstName;

  String id;

  String lastName;

  String position;

  DateTime updatedAt;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Staff &&
          other.createdAt == createdAt &&
          other.deletedAt == deletedAt &&
          other.firstName == firstName &&
          other.id == id &&
          other.lastName == lastName &&
          other.position == position &&
          other.updatedAt == updatedAt;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (createdAt.hashCode) +
      (deletedAt == null ? 0 : deletedAt!.hashCode) +
      (firstName.hashCode) +
      (id.hashCode) +
      (lastName.hashCode) +
      (position.hashCode) +
      (updatedAt.hashCode);

  @override
  String toString() =>
      'Staff[createdAt=$createdAt, deletedAt=$deletedAt, firstName=$firstName, id=$id, lastName=$lastName, position=$position, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
    if (this.deletedAt != null) {
      json[r'deleted_at'] = this.deletedAt!.toUtc().toIso8601String();
    } else {
      json[r'deleted_at'] = null;
    }
    json[r'first_name'] = this.firstName;
    json[r'id'] = this.id;
    json[r'last_name'] = this.lastName;
    json[r'position'] = this.position;
    json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [Staff] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Staff? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'created_at'),
            'Required key "Staff[created_at]" is missing from JSON.');
        assert(json[r'created_at'] != null,
            'Required key "Staff[created_at]" has a null value in JSON.');
        assert(json.containsKey(r'deleted_at'),
            'Required key "Staff[deleted_at]" is missing from JSON.');
        assert(json.containsKey(r'first_name'),
            'Required key "Staff[first_name]" is missing from JSON.');
        assert(json[r'first_name'] != null,
            'Required key "Staff[first_name]" has a null value in JSON.');
        assert(json.containsKey(r'id'),
            'Required key "Staff[id]" is missing from JSON.');
        assert(json[r'id'] != null,
            'Required key "Staff[id]" has a null value in JSON.');
        assert(json.containsKey(r'last_name'),
            'Required key "Staff[last_name]" is missing from JSON.');
        assert(json[r'last_name'] != null,
            'Required key "Staff[last_name]" has a null value in JSON.');
        assert(json.containsKey(r'position'),
            'Required key "Staff[position]" is missing from JSON.');
        assert(json[r'position'] != null,
            'Required key "Staff[position]" has a null value in JSON.');
        assert(json.containsKey(r'updated_at'),
            'Required key "Staff[updated_at]" is missing from JSON.');
        assert(json[r'updated_at'] != null,
            'Required key "Staff[updated_at]" has a null value in JSON.');
        return true;
      }());

      return Staff(
        createdAt: mapDateTime(json, r'created_at', r'')!,
        deletedAt: mapDateTime(json, r'deleted_at', r''),
        firstName: mapValueOfType<String>(json, r'first_name')!,
        id: mapValueOfType<String>(json, r'id')!,
        lastName: mapValueOfType<String>(json, r'last_name')!,
        position: mapValueOfType<String>(json, r'position')!,
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
      );
    }
    return null;
  }

  static List<Staff> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <Staff>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Staff.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Staff> mapFromJson(dynamic json) {
    final map = <String, Staff>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Staff.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Staff-objects as value to a dart map
  static Map<String, List<Staff>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<Staff>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Staff.listFromJson(
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
    'deleted_at',
    'first_name',
    'id',
    'last_name',
    'position',
    'updated_at',
  };
}
