//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Diagnosis {
  /// Returns a new [Diagnosis] instance.
  Diagnosis({
    required this.code,
    required this.createdAt,
    required this.deletedAt,
    required this.id,
    required this.name,
    required this.updatedAt,
  });

  String code;

  DateTime createdAt;

  DateTime? deletedAt;

  String id;

  String name;

  DateTime updatedAt;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Diagnosis &&
          other.code == code &&
          other.createdAt == createdAt &&
          other.deletedAt == deletedAt &&
          other.id == id &&
          other.name == name &&
          other.updatedAt == updatedAt;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (code.hashCode) +
      (createdAt.hashCode) +
      (deletedAt == null ? 0 : deletedAt!.hashCode) +
      (id.hashCode) +
      (name.hashCode) +
      (updatedAt.hashCode);

  @override
  String toString() =>
      'Diagnosis[code=$code, createdAt=$createdAt, deletedAt=$deletedAt, id=$id, name=$name, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'code'] = this.code;
    json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
    if (this.deletedAt != null) {
      json[r'deleted_at'] = this.deletedAt!.toUtc().toIso8601String();
    } else {
      json[r'deleted_at'] = null;
    }
    json[r'id'] = this.id;
    json[r'name'] = this.name;
    json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [Diagnosis] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Diagnosis? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'code'),
            'Required key "Diagnosis[code]" is missing from JSON.');
        assert(json[r'code'] != null,
            'Required key "Diagnosis[code]" has a null value in JSON.');
        assert(json.containsKey(r'created_at'),
            'Required key "Diagnosis[created_at]" is missing from JSON.');
        assert(json[r'created_at'] != null,
            'Required key "Diagnosis[created_at]" has a null value in JSON.');
        assert(json.containsKey(r'deleted_at'),
            'Required key "Diagnosis[deleted_at]" is missing from JSON.');
        assert(json.containsKey(r'id'),
            'Required key "Diagnosis[id]" is missing from JSON.');
        assert(json[r'id'] != null,
            'Required key "Diagnosis[id]" has a null value in JSON.');
        assert(json.containsKey(r'name'),
            'Required key "Diagnosis[name]" is missing from JSON.');
        assert(json[r'name'] != null,
            'Required key "Diagnosis[name]" has a null value in JSON.');
        assert(json.containsKey(r'updated_at'),
            'Required key "Diagnosis[updated_at]" is missing from JSON.');
        assert(json[r'updated_at'] != null,
            'Required key "Diagnosis[updated_at]" has a null value in JSON.');
        return true;
      }());

      return Diagnosis(
        code: mapValueOfType<String>(json, r'code')!,
        createdAt: mapDateTime(json, r'created_at', r'')!,
        deletedAt: mapDateTime(json, r'deleted_at', r''),
        id: mapValueOfType<String>(json, r'id')!,
        name: mapValueOfType<String>(json, r'name')!,
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
      );
    }
    return null;
  }

  static List<Diagnosis> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <Diagnosis>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Diagnosis.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Diagnosis> mapFromJson(dynamic json) {
    final map = <String, Diagnosis>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Diagnosis.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Diagnosis-objects as value to a dart map
  static Map<String, List<Diagnosis>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<Diagnosis>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Diagnosis.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'code',
    'created_at',
    'deleted_at',
    'id',
    'name',
    'updated_at',
  };
}
