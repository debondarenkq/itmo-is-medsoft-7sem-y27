import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:his_desktop/core/api_client.dart';

void main() {
  test('all resources and import go through one configured gateway', () async {
    final urls = <Uri>[];
    final api = ApiClient(
      baseUrl: 'http://gateway:8080',
      client: MockClient((request) async {
        urls.add(request.url);
        if (request.method == 'POST') {
          expect(request.body, '{"diagnoses":[]}');
          return http.Response('{"items":[]}', 200);
        }
        return http.Response('{"items":[],"total":0}', 200);
      }),
    );
    await api.staff();
    await api.diagnoses();
    await api.request('GET', '/api/v1/patients');
    await api.request(
      'POST',
      '/api/v1/diagnoses/import',
      fileContents: '{"diagnoses":[]}',
    );
    expect(
      urls.every((uri) => uri.host == 'gateway' && uri.port == 8080),
      isTrue,
    );
    api.close();
  });
  test(
    'API errors preserve version-conflict and field validation data',
    () async {
      final api = ApiClient(
        client: MockClient(
          (_) async => http.Response(
            jsonEncode({
              'error': {
                'code': 'VERSION_CONFLICT',
                'message': 'Карта изменена',
                'fields': {'name': 'Неверное значение'},
              },
            }),
            409,
            headers: {'content-type': 'application/json; charset=utf-8'},
          ),
        ),
      );
      await expectLater(
        api.request('POST', '/api/v1/records/r/changes', body: {}),
        throwsA(
          isA<ApiException>()
              .having((e) => e.code, 'code', 'VERSION_CONFLICT')
              .having((e) => e.fields['name'], 'field', 'Неверное значение'),
        ),
      );
      api.close();
    },
  );
}
