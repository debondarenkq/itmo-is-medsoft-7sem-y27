//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

enum CommandType {
  addDiagnosis._(r'add_diagnosis'),
  removeDiagnosis._(r'remove_diagnosis'),
  addPrescription._(r'add_prescription'),
  editPrescription._(r'edit_prescription'),
  cancelPrescription._(r'cancel_prescription'),
  removePrescription._(r'remove_prescription'),
  ;

  /// Instantiate a new enum with the provided value.
  const CommandType._(this._value);

  /// The underlying value of this enum member.
  final String _value;

  @override
  String toString() => _value;

  /// Encodes this enum as a value suitable for JSON.
  String toJson() => _value;

  /// Returns the instance of [CommandType] that was successfully decoded
  /// from the passed [value] on success, null otherwise.
  static CommandType? fromJson(dynamic value) =>
      CommandTypeTypeTransformer().decode(value);

  /// Returns a [List] containing instances of [CommandType]
  /// that were successfully decoded from the passed [JSON][json].
  static List<CommandType> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <CommandType>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = CommandType.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [CommandType] to String,
/// and [decode] dynamic data back to [CommandType].
class CommandTypeTypeTransformer {
  factory CommandTypeTypeTransformer() =>
      _instance ??= const CommandTypeTypeTransformer._();

  const CommandTypeTypeTransformer._();

  /// Encodes this enum as a value suitable for JSON.
  String encode(CommandType data) => data._value;

  /// Returns the instance of [CommandType] that was successfully decoded
  /// from the passed [data] value on success, null otherwise.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  CommandType? decode(dynamic data, {bool allowNull = true}) {
    if (data is CommandType) {
      return data;
    }
    if (data != null) {
      switch (data) {
        case r'add_diagnosis':
          return CommandType.addDiagnosis;
        case r'remove_diagnosis':
          return CommandType.removeDiagnosis;
        case r'add_prescription':
          return CommandType.addPrescription;
        case r'edit_prescription':
          return CommandType.editPrescription;
        case r'cancel_prescription':
          return CommandType.cancelPrescription;
        case r'remove_prescription':
          return CommandType.removePrescription;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// The singleton instance of this transformer.
  static CommandTypeTypeTransformer? _instance;
}
