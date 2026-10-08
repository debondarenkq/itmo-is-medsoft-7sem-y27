//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class State {
  /// Returns a new [State] instance.
  State({
    this.diagnoses = const [],
    this.prescriptions = const [],
  });

  List<RecordDiagnosis> diagnoses;

  List<Prescription> prescriptions;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is State &&
          _deepEquality.equals(other.diagnoses, diagnoses) &&
          _deepEquality.equals(other.prescriptions, prescriptions);

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (diagnoses.hashCode) + (prescriptions.hashCode);

  @override
  String toString() =>
      'State[diagnoses=$diagnoses, prescriptions=$prescriptions]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'diagnoses'] = this.diagnoses;
    json[r'prescriptions'] = this.prescriptions;
    return json;
  }

  /// Returns a new [State] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static State? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'diagnoses'),
            'Required key "State[diagnoses]" is missing from JSON.');
        assert(json[r'diagnoses'] != null,
            'Required key "State[diagnoses]" has a null value in JSON.');
        assert(json.containsKey(r'prescriptions'),
            'Required key "State[prescriptions]" is missing from JSON.');
        assert(json[r'prescriptions'] != null,
            'Required key "State[prescriptions]" has a null value in JSON.');
        return true;
      }());

      return State(
        diagnoses: RecordDiagnosis.listFromJson(json[r'diagnoses']),
        prescriptions: Prescription.listFromJson(json[r'prescriptions']),
      );
    }
    return null;
  }

  static List<State> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <State>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = State.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, State> mapFromJson(dynamic json) {
    final map = <String, State>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = State.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of State-objects as value to a dart map
  static Map<String, List<State>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<State>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = State.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'diagnoses',
    'prescriptions',
  };
}
