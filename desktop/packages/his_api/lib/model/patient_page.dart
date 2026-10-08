//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class PatientPage {
  /// Returns a new [PatientPage] instance.
  PatientPage({
    this.items = const [],
    required this.limit,
    required this.offset,
    required this.total,
  });

  List<Patient> items;

  int limit;

  int offset;

  int total;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is PatientPage &&
          _deepEquality.equals(other.items, items) &&
          other.limit == limit &&
          other.offset == offset &&
          other.total == total;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (items.hashCode) +
      (limit.hashCode) +
      (offset.hashCode) +
      (total.hashCode);

  @override
  String toString() =>
      'PatientPage[items=$items, limit=$limit, offset=$offset, total=$total]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'items'] = this.items;
    json[r'limit'] = this.limit;
    json[r'offset'] = this.offset;
    json[r'total'] = this.total;
    return json;
  }

  /// Returns a new [PatientPage] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static PatientPage? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'items'),
            'Required key "PatientPage[items]" is missing from JSON.');
        assert(json[r'items'] != null,
            'Required key "PatientPage[items]" has a null value in JSON.');
        assert(json.containsKey(r'limit'),
            'Required key "PatientPage[limit]" is missing from JSON.');
        assert(json[r'limit'] != null,
            'Required key "PatientPage[limit]" has a null value in JSON.');
        assert(json.containsKey(r'offset'),
            'Required key "PatientPage[offset]" is missing from JSON.');
        assert(json[r'offset'] != null,
            'Required key "PatientPage[offset]" has a null value in JSON.');
        assert(json.containsKey(r'total'),
            'Required key "PatientPage[total]" is missing from JSON.');
        assert(json[r'total'] != null,
            'Required key "PatientPage[total]" has a null value in JSON.');
        return true;
      }());

      return PatientPage(
        items: Patient.listFromJson(json[r'items']),
        limit: mapValueOfType<int>(json, r'limit')!,
        offset: mapValueOfType<int>(json, r'offset')!,
        total: mapValueOfType<int>(json, r'total')!,
      );
    }
    return null;
  }

  static List<PatientPage> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <PatientPage>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = PatientPage.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, PatientPage> mapFromJson(dynamic json) {
    final map = <String, PatientPage>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = PatientPage.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of PatientPage-objects as value to a dart map
  static Map<String, List<PatientPage>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<PatientPage>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = PatientPage.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'items',
    'limit',
    'offset',
    'total',
  };
}
