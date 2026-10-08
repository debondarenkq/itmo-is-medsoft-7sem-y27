//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

enum EventType {
  recordCreated._(r'record_created'),
  diagnosisAdded._(r'diagnosis_added'),
  diagnosisRemoved._(r'diagnosis_removed'),
  prescriptionAdded._(r'prescription_added'),
  prescriptionCancelled._(r'prescription_cancelled'),
  prescriptionRemoved._(r'prescription_removed'),
  ;

  /// Instantiate a new enum with the provided value.
  const EventType._(this._value);

  /// The underlying value of this enum member.
  final String _value;

  @override
  String toString() => _value;

  /// Encodes this enum as a value suitable for JSON.
  String toJson() => _value;

  /// Returns the instance of [EventType] that was successfully decoded
  /// from the passed [value] on success, null otherwise.
  static EventType? fromJson(dynamic value) =>
      EventTypeTypeTransformer().decode(value);

  /// Returns a [List] containing instances of [EventType]
  /// that were successfully decoded from the passed [JSON][json].
  static List<EventType> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <EventType>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = EventType.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [EventType] to String,
/// and [decode] dynamic data back to [EventType].
class EventTypeTypeTransformer {
  factory EventTypeTypeTransformer() =>
      _instance ??= const EventTypeTypeTransformer._();

  const EventTypeTypeTransformer._();

  /// Encodes this enum as a value suitable for JSON.
  String encode(EventType data) => data._value;

  /// Returns the instance of [EventType] that was successfully decoded
  /// from the passed [data] value on success, null otherwise.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  EventType? decode(dynamic data, {bool allowNull = true}) {
    if (data is EventType) {
      return data;
    }
    if (data != null) {
      switch (data) {
        case r'record_created':
          return EventType.recordCreated;
        case r'diagnosis_added':
          return EventType.diagnosisAdded;
        case r'diagnosis_removed':
          return EventType.diagnosisRemoved;
        case r'prescription_added':
          return EventType.prescriptionAdded;
        case r'prescription_cancelled':
          return EventType.prescriptionCancelled;
        case r'prescription_removed':
          return EventType.prescriptionRemoved;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// The singleton instance of this transformer.
  static EventTypeTypeTransformer? _instance;
}
