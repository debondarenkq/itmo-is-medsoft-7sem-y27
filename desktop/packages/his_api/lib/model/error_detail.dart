//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ErrorDetail {
  /// Returns a new [ErrorDetail] instance.
  ErrorDetail({
    required this.code,
    this.fields = const Optional.present(const {}),
    required this.message,
  });

  String code;

  Optional<Map<String, String>?> fields;

  String message;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is ErrorDetail &&
          other.code == code &&
          _deepEquality.equals(other.fields, fields) &&
          other.message == message;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (code.hashCode) + (fields.hashCode) + (message.hashCode);

  @override
  String toString() =>
      'ErrorDetail[code=$code, fields=$fields, message=$message]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'code'] = this.code;
    if (this.fields.isPresent) {
      final value = this.fields.value;
      json[r'fields'] = value;
    }
    json[r'message'] = this.message;
    return json;
  }

  /// Returns a new [ErrorDetail] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ErrorDetail? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'code'),
            'Required key "ErrorDetail[code]" is missing from JSON.');
        assert(json[r'code'] != null,
            'Required key "ErrorDetail[code]" has a null value in JSON.');
        assert(json.containsKey(r'message'),
            'Required key "ErrorDetail[message]" is missing from JSON.');
        assert(json[r'message'] != null,
            'Required key "ErrorDetail[message]" has a null value in JSON.');
        return true;
      }());

      return ErrorDetail(
        code: mapValueOfType<String>(json, r'code')!,
        fields: json.containsKey(r'fields')
            ? Optional.present(mapCastOfType<String, String>(json, r'fields'))
            : const Optional.absent(),
        message: mapValueOfType<String>(json, r'message')!,
      );
    }
    return null;
  }

  static List<ErrorDetail> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <ErrorDetail>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ErrorDetail.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ErrorDetail> mapFromJson(dynamic json) {
    final map = <String, ErrorDetail>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ErrorDetail.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ErrorDetail-objects as value to a dart map
  static Map<String, List<ErrorDetail>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<ErrorDetail>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ErrorDetail.listFromJson(
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
    'message',
  };
}
