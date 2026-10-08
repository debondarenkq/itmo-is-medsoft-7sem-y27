//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Prescription {
  /// Returns a new [Prescription] instance.
  Prescription({
    required this.id,
    required this.status,
    required this.text,
  });

  String id;

  PrescriptionStatus status;

  String text;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Prescription &&
          other.id == id &&
          other.status == status &&
          other.text == text;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (id.hashCode) + (status.hashCode) + (text.hashCode);

  @override
  String toString() => 'Prescription[id=$id, status=$status, text=$text]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'id'] = this.id;
    json[r'status'] = this.status;
    json[r'text'] = this.text;
    return json;
  }

  /// Returns a new [Prescription] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Prescription? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'id'),
            'Required key "Prescription[id]" is missing from JSON.');
        assert(json[r'id'] != null,
            'Required key "Prescription[id]" has a null value in JSON.');
        assert(json.containsKey(r'status'),
            'Required key "Prescription[status]" is missing from JSON.');
        assert(json[r'status'] != null,
            'Required key "Prescription[status]" has a null value in JSON.');
        assert(json.containsKey(r'text'),
            'Required key "Prescription[text]" is missing from JSON.');
        assert(json[r'text'] != null,
            'Required key "Prescription[text]" has a null value in JSON.');
        return true;
      }());

      return Prescription(
        id: mapValueOfType<String>(json, r'id')!,
        status: PrescriptionStatus.fromJson(json[r'status'])!,
        text: mapValueOfType<String>(json, r'text')!,
      );
    }
    return null;
  }

  static List<Prescription> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <Prescription>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Prescription.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Prescription> mapFromJson(dynamic json) {
    final map = <String, Prescription>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Prescription.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Prescription-objects as value to a dart map
  static Map<String, List<Prescription>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<Prescription>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Prescription.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'id',
    'status',
    'text',
  };
}
