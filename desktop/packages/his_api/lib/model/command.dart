//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Command {
  /// Returns a new [Command] instance.
  Command({
    this.diagnosisId = const Optional.absent(),
    this.id = const Optional.absent(),
    this.text = const Optional.absent(),
    required this.type,
  });

  /// Обязателен для add_diagnosis; ID из справочника.
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  Optional<String?> diagnosisId;

  /// Обязателен для удаления диагноза и изменения/отмены/удаления назначения; ID записи внутри ЭМК.
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  Optional<String?> id;

  /// Обязателен для add_prescription и edit_prescription.
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  Optional<String?> text;

  CommandType type;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Command &&
          other.diagnosisId == diagnosisId &&
          other.id == id &&
          other.text == text &&
          other.type == type;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (diagnosisId == null ? 0 : diagnosisId!.hashCode) +
      (id == null ? 0 : id!.hashCode) +
      (text == null ? 0 : text!.hashCode) +
      (type.hashCode);

  @override
  String toString() =>
      'Command[diagnosisId=$diagnosisId, id=$id, text=$text, type=$type]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.diagnosisId.isPresent) {
      final value = this.diagnosisId.value;
      json[r'diagnosis_id'] = value;
    }
    if (this.id.isPresent) {
      final value = this.id.value;
      json[r'id'] = value;
    }
    if (this.text.isPresent) {
      final value = this.text.value;
      json[r'text'] = value;
    }
    json[r'type'] = this.type;
    return json;
  }

  /// Returns a new [Command] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Command? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'type'),
            'Required key "Command[type]" is missing from JSON.');
        assert(json[r'type'] != null,
            'Required key "Command[type]" has a null value in JSON.');
        return true;
      }());

      return Command(
        diagnosisId: json.containsKey(r'diagnosis_id')
            ? Optional.present(mapValueOfType<String>(json, r'diagnosis_id'))
            : const Optional.absent(),
        id: json.containsKey(r'id')
            ? Optional.present(mapValueOfType<String>(json, r'id'))
            : const Optional.absent(),
        text: json.containsKey(r'text')
            ? Optional.present(mapValueOfType<String>(json, r'text'))
            : const Optional.absent(),
        type: CommandType.fromJson(json[r'type'])!,
      );
    }
    return null;
  }

  static List<Command> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <Command>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Command.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Command> mapFromJson(dynamic json) {
    final map = <String, Command>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Command.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Command-objects as value to a dart map
  static Map<String, List<Command>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<Command>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Command.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'type',
  };
}
