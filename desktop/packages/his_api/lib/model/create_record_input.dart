//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class CreateRecordInput {
  /// Returns a new [CreateRecordInput] instance.
  CreateRecordInput({
    required this.staffId,
  });

  String staffId;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is CreateRecordInput && other.staffId == staffId;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (staffId.hashCode);

  @override
  String toString() => 'CreateRecordInput[staffId=$staffId]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'staff_id'] = this.staffId;
    return json;
  }

  /// Returns a new [CreateRecordInput] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static CreateRecordInput? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'staff_id'),
            'Required key "CreateRecordInput[staff_id]" is missing from JSON.');
        assert(json[r'staff_id'] != null,
            'Required key "CreateRecordInput[staff_id]" has a null value in JSON.');
        return true;
      }());

      return CreateRecordInput(
        staffId: mapValueOfType<String>(json, r'staff_id')!,
      );
    }
    return null;
  }

  static List<CreateRecordInput> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <CreateRecordInput>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = CreateRecordInput.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, CreateRecordInput> mapFromJson(dynamic json) {
    final map = <String, CreateRecordInput>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = CreateRecordInput.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of CreateRecordInput-objects as value to a dart map
  static Map<String, List<CreateRecordInput>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<CreateRecordInput>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = CreateRecordInput.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'staff_id',
  };
}
