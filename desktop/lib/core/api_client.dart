import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'models.dart';

class ApiException implements Exception {
  const ApiException(
    this.message, {
    this.status = 0,
    this.code = '',
    this.fields = const {},
  });
  final String message, code;
  final int status;
  final Map<String, String> fields;
  @override
  String toString() => [
    message,
    ...fields.entries.map((e) => '${e.key}: ${e.value}'),
  ].join('\n');
}

// The desktop client has exactly one backend address: the public gateway.
class ApiClient {
  ApiClient({
    String baseUrl = const String.fromEnvironment(
      'API_BASE_URL',
      defaultValue: 'http://localhost:8080',
    ),
    http.Client? client,
  }) : _client = client ?? http.Client() {
    configure(baseUrl);
  }
  final http.Client _client;
  late Uri _base;
  String get baseUrl => _base.toString();
  void configure(String value) {
    final uri = Uri.tryParse(value.trim());
    if (uri == null ||
        !['http', 'https'].contains(uri.scheme) ||
        uri.host.isEmpty ||
        uri.userInfo.isNotEmpty ||
        uri.hasQuery ||
        uri.hasFragment) {
      throw const FormatException('Укажите адрес вида http://localhost:8080');
    }
    _base = uri.replace(path: uri.path.replaceFirst(RegExp(r'/+$'), ''));
  }

  void close() => _client.close();

  Future<Json> request(
    String method,
    String path, {
    Json? body,
    String? fileContents,
    Map<String, String>? query,
  }) async {
    final uri = _base.replace(
      path: '${_base.path}$path',
      queryParameters: query,
    );
    final request = http.Request(method, uri);
    if (body != null || fileContents != null) {
      request.headers['Content-Type'] = 'application/json';
      request.body = fileContents ?? jsonEncode(body);
    }
    try {
      final streamed = await _client
          .send(request)
          .timeout(const Duration(seconds: 20));
      final response = await http.Response.fromStream(
        streamed,
      ).timeout(const Duration(seconds: 20));
      Json json = {};
      if (response.bodyBytes.isNotEmpty) {
        try {
          json = jsonDecode(utf8.decode(response.bodyBytes)) as Json;
        } on FormatException {
          throw ApiException(
            'Сервер временно недоступен. Повторите запрос.',
            status: response.statusCode,
          );
        }
      }
      if (response.statusCode < 200 || response.statusCode >= 300) {
        final error = json['error'] as Json? ?? {};
        throw ApiException(
          error['message'] as String? ?? 'Не удалось выполнить запрос.',
          status: response.statusCode,
          code: error['code'] as String? ?? '',
          fields: (error['fields'] as Map? ?? {}).map(
            (key, value) => MapEntry(key.toString(), value.toString()),
          ),
        );
      }
      return json;
    } on TimeoutException {
      throw const ApiException(
        'Сервер не ответил вовремя. Попробуйте ещё раз.',
      );
    } on http.ClientException {
      throw const ApiException(
        'Нет соединения с сервером. Проверьте подключение и адрес сервера.',
      );
    } on SocketException {
      throw const ApiException('Не удалось подключиться к серверу.');
    }
  }

  Future<PageResult<T>> page<T>(
    String path,
    T Function(Json) parse, {
    int offset = 0,
    String query = '',
    bool includeDeleted = false,
  }) async {
    final data = await request(
      'GET',
      path,
      query: {
        'limit': '50',
        'offset': '$offset',
        'q': query,
        if (includeDeleted) 'include_deleted': 'true',
      },
    );
    return PageResult(
      (data['items'] as List).map((v) => parse(v as Json)).toList(),
      data['total'] as int,
    );
  }

  Future<List<T>> all<T>(String path, T Function(Json) parse) async {
    final result = <T>[];
    while (true) {
      final data = await request(
        'GET',
        path,
        query: {'limit': '200', 'offset': '${result.length}'},
      );
      final items = (data['items'] as List)
          .map((v) => parse(v as Json))
          .toList();
      result.addAll(items);
      if (items.isEmpty || result.length >= (data['total'] as int)) {
        return result;
      }
    }
  }

  Future<List<Staff>> staff() => all('/api/v1/staff', Staff.fromJson);
  Future<List<Diagnosis>> diagnoses() =>
      all('/api/v1/diagnoses', Diagnosis.fromJson);
  Future<Patient> patient(String id) async =>
      Patient.fromJson(await request('GET', '/api/v1/patients/$id'));
  Future<MedicalRecord> patientRecord(String id) async =>
      MedicalRecord.fromJson(
        await request('GET', '/api/v1/patients/$id/record'),
      );
  Future<MedicalRecord> createRecord(String id, String staffId) async =>
      MedicalRecord.fromJson(
        await request(
          'POST',
          '/api/v1/patients/$id/record',
          body: {'staff_id': staffId},
        ),
      );
  Future<List<RecordEvent>> history(String id) =>
      all('/api/v1/records/$id/history', RecordEvent.fromJson);
  Future<MedicalRecord> saveRecord(
    MedicalRecord base,
    String staffId,
    List<Json> commands,
  ) async => MedicalRecord.fromJson(
    await request(
      'POST',
      '/api/v1/records/${base.id}/changes',
      body: {
        'staff_id': staffId,
        'expected_version': base.version,
        'commands': commands,
      },
    ),
  );
  Future<MedicalRecord> historicalRecord(
    String id, {
    int? version,
    DateTime? at,
  }) async => MedicalRecord.fromJson(
    await request(
      'GET',
      '/api/v1/records/$id/state',
      query: {
        if (version != null) 'version': '$version',
        if (at != null) 'at': at.toUtc().toIso8601String(),
      },
    ),
  );
}
