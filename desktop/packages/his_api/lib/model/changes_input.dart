//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ChangesInput {
  /// Returns a new [ChangesInput] instance.
  ChangesInput({
    this.commands = const [],
    required this.expectedVersion,
    required this.staffId,
  });

  List<Command> commands;

  /// Minimum value: 1
  int expectedVersion;

  String staffId;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is ChangesInput &&
          _deepEquality.equals(other.commands, commands) &&
          other.expectedVersion == expectedVersion &&
          other.staffId == staffId;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (commands.hashCode) + (expectedVersion.hashCode) + (staffId.hashCode);

  @override
  String toString() =>
      'ChangesInput[commands=$commands, expectedVersion=$expectedVersion, staffId=$staffId]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'commands'] = this.commands;
    json[r'expected_version'] = this.expectedVersion;
    json[r'staff_id'] = this.staffId;
    return json;
  }

  /// Returns a new [ChangesInput] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ChangesInput? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'commands'),
            'Required key "ChangesInput[commands]" is missing from JSON.');
        assert(json[r'commands'] != null,
            'Required key "ChangesInput[commands]" has a null value in JSON.');
        assert(json.containsKey(r'expected_version'),
            'Required key "ChangesInput[expected_version]" is missing from JSON.');
        assert(json[r'expected_version'] != null,
            'Required key "ChangesInput[expected_version]" has a null value in JSON.');
        assert(json.containsKey(r'staff_id'),
            'Required key "ChangesInput[staff_id]" is missing from JSON.');
        assert(json[r'staff_id'] != null,
            'Required key "ChangesInput[staff_id]" has a null value in JSON.');
        return true;
      }());

      return ChangesInput(
        commands: Command.listFromJson(json[r'commands']),
        expectedVersion: mapValueOfType<int>(json, r'expected_version')!,
        staffId: mapValueOfType<String>(json, r'staff_id')!,
      );
    }
    return null;
  }

  static List<ChangesInput> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <ChangesInput>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ChangesInput.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ChangesInput> mapFromJson(dynamic json) {
    final map = <String, ChangesInput>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ChangesInput.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ChangesInput-objects as value to a dart map
  static Map<String, List<ChangesInput>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<ChangesInput>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ChangesInput.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'commands',
    'expected_version',
    'staff_id',
  };
}
