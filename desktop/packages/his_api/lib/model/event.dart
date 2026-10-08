//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Event {
  /// Returns a new [Event] instance.
  Event({
    required this.actor,
    this.after = const {},
    this.before = const {},
    required this.entityId,
    required this.id,
    required this.occurredAt,
    required this.recordId,
    required this.sequence,
    required this.type,
    required this.version,
  });

  Actor actor;

  /// Полный снимок данных записи; null означает отсутствие записи.
  Map<String, Object>? after;

  /// Полный снимок данных записи; null означает отсутствие записи.
  Map<String, Object>? before;

  String entityId;

  String id;

  DateTime occurredAt;

  String recordId;

  int sequence;

  EventType type;

  int version;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Event &&
          other.actor == actor &&
          _deepEquality.equals(other.after, after) &&
          _deepEquality.equals(other.before, before) &&
          other.entityId == entityId &&
          other.id == id &&
          other.occurredAt == occurredAt &&
          other.recordId == recordId &&
          other.sequence == sequence &&
          other.type == type &&
          other.version == version;

  @override
  int get hashCode =>
      // ignore: unnecessary_parenthesis
      (actor.hashCode) +
      (after == null ? 0 : after!.hashCode) +
      (before == null ? 0 : before!.hashCode) +
      (entityId.hashCode) +
      (id.hashCode) +
      (occurredAt.hashCode) +
      (recordId.hashCode) +
      (sequence.hashCode) +
      (type.hashCode) +
      (version.hashCode);

  @override
  String toString() =>
      'Event[actor=$actor, after=$after, before=$before, entityId=$entityId, id=$id, occurredAt=$occurredAt, recordId=$recordId, sequence=$sequence, type=$type, version=$version]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    json[r'actor'] = this.actor;
    if (this.after != null) {
      json[r'after'] = this.after;
    } else {
      json[r'after'] = null;
    }
    if (this.before != null) {
      json[r'before'] = this.before;
    } else {
      json[r'before'] = null;
    }
    json[r'entity_id'] = this.entityId;
    json[r'id'] = this.id;
    json[r'occurred_at'] = this.occurredAt.toUtc().toIso8601String();
    json[r'record_id'] = this.recordId;
    json[r'sequence'] = this.sequence;
    json[r'type'] = this.type;
    json[r'version'] = this.version;
    return json;
  }

  /// Returns a new [Event] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Event? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'actor'),
            'Required key "Event[actor]" is missing from JSON.');
        assert(json[r'actor'] != null,
            'Required key "Event[actor]" has a null value in JSON.');
        assert(json.containsKey(r'after'),
            'Required key "Event[after]" is missing from JSON.');
        assert(json.containsKey(r'before'),
            'Required key "Event[before]" is missing from JSON.');
        assert(json.containsKey(r'entity_id'),
            'Required key "Event[entity_id]" is missing from JSON.');
        assert(json[r'entity_id'] != null,
            'Required key "Event[entity_id]" has a null value in JSON.');
        assert(json.containsKey(r'id'),
            'Required key "Event[id]" is missing from JSON.');
        assert(json[r'id'] != null,
            'Required key "Event[id]" has a null value in JSON.');
        assert(json.containsKey(r'occurred_at'),
            'Required key "Event[occurred_at]" is missing from JSON.');
        assert(json[r'occurred_at'] != null,
            'Required key "Event[occurred_at]" has a null value in JSON.');
        assert(json.containsKey(r'record_id'),
            'Required key "Event[record_id]" is missing from JSON.');
        assert(json[r'record_id'] != null,
            'Required key "Event[record_id]" has a null value in JSON.');
        assert(json.containsKey(r'sequence'),
            'Required key "Event[sequence]" is missing from JSON.');
        assert(json[r'sequence'] != null,
            'Required key "Event[sequence]" has a null value in JSON.');
        assert(json.containsKey(r'type'),
            'Required key "Event[type]" is missing from JSON.');
        assert(json[r'type'] != null,
            'Required key "Event[type]" has a null value in JSON.');
        assert(json.containsKey(r'version'),
            'Required key "Event[version]" is missing from JSON.');
        assert(json[r'version'] != null,
            'Required key "Event[version]" has a null value in JSON.');
        return true;
      }());

      return Event(
        actor: Actor.fromJson(json[r'actor'])!,
        after: mapCastOfType<String, Object>(json, r'after'),
        before: mapCastOfType<String, Object>(json, r'before'),
        entityId: mapValueOfType<String>(json, r'entity_id')!,
        id: mapValueOfType<String>(json, r'id')!,
        occurredAt: mapDateTime(json, r'occurred_at', r'')!,
        recordId: mapValueOfType<String>(json, r'record_id')!,
        sequence: mapValueOfType<int>(json, r'sequence')!,
        type: EventType.fromJson(json[r'type'])!,
        version: mapValueOfType<int>(json, r'version')!,
      );
    }
    return null;
  }

  static List<Event> listFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final result = <Event>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Event.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Event> mapFromJson(dynamic json) {
    final map = <String, Event>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Event.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Event-objects as value to a dart map
  static Map<String, List<Event>> mapListFromJson(
    dynamic json, {
    bool growable = false,
  }) {
    final map = <String, List<Event>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Event.listFromJson(
          entry.value,
          growable: growable,
        );
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'actor',
    'after',
    'before',
    'entity_id',
    'id',
    'occurred_at',
    'record_id',
    'sequence',
    'type',
    'version',
  };
}
