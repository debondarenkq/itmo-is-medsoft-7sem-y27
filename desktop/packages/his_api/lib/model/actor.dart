//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Actor {
  /// Returns a new [Actor] instance.
  Actor({
    required this.firstName,
    required this.id,
    required this.lastName,
    required this.position,
  });

  String firstName;

  String id;

  String lastName;

  String position;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Actor &&
          other.firstName == firstName &&
          other.id == id &&
          other.lastName == lastName &&
          other.position == position;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (firstName.hashCode) +
      (id.hashCode) +
      (lastName.hashCode) +
      (position.hashCode);

  @override
  String toString() =>
      'Actor[firstName=$firstName, id=$id, lastName=$lastName, position=$position]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'first_name'] = this.firstName;
    json[r'id'] = this.id;
    json[r'last_name'] = this.lastName;
    json[r'position'] = this.position;
    return json;
  }

  /// Returns a new [Actor] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Actor? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'first_name'),
            'Required key "Actor[first_name]" is missing from JSON.');
        assert(json[r'first_name'] != null,
            'Required key "Actor[first_name]" has a null value in JSON.');
        assert(json.containsKey(r'id'),
            'Required key "Actor[id]" is missing from JSON.');
        assert(json[r'id'] != null,
            'Required key "Actor[id]" has a null value in JSON.');
        assert(json.containsKey(r'last_name'),
            'Required key "Actor[last_name]" is missing from JSON.');
        assert(json[r'last_name'] != null,
            'Required key "Actor[last_name]" has a null value in JSON.');
        assert(json.containsKey(r'position'),
            'Required key "Actor[position]" is missing from JSON.');
        assert(json[r'position'] != null,
            'Required key "Actor[position]" has a null value in JSON.');
        return true;
      }());

      return Actor(
        firstName: mapValueOfType<String>(json, r'first_name')!,
        id: mapValueOfType<String>(json, r'id')!,
        lastName: mapValueOfType<String>(json, r'last_name')!,
        position: mapValueOfType<String>(json, r'position')!,
      );
    }
    return null;
  }

  static List<Actor> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <Actor>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Actor.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Actor> mapFromJson(dynamic json) {
    final map = <String, Actor>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Actor.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Actor-objects as value to a dart map
  static Map<String, List<Actor>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<Actor>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Actor.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'first_name',
    'id',
    'last_name',
    'position',
  };
}
