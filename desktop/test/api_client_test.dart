import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:his_api/api.dart' as contract;
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:his_desktop/core/api_client.dart';

void main() {
  test('all generated APIs and import use one configured gateway', () async {
    final urls = <Uri>[];
    final api = ApiClient(
      baseUrl: 'http://gateway:8080',
      client: MockClient((request) async {
        urls.add(request.url);
        if (request.method == 'POST') {
          expect(jsonDecode(request.body), {'diagnoses': []});
          return http.Response('{"items":[]}', 200);
        }
        return http.Response(
          '{"items":[],"total":0,"limit":50,"offset":0}',
          200,
        );
      }),
    );
    await api.staff();
    await api.diagnoses();
    await api.patientPage();
    await api.importDiagnoses('{"diagnoses":[]}');
    expect(
      urls.every((uri) => uri.host == 'gateway' && uri.port == 8080),
      isTrue,
    );
    api.close();
  });
  test('generated command serialization omits unrelated optional fields', () {
    final command = contract.Command(
      type: contract.CommandType.addPrescription,
      text: const contract.Optional.present('Пить воду'),
    );
    expect(jsonDecode(jsonEncode(command)), {
      'type': 'add_prescription',
      'text': 'Пить воду',
    });
  });
  test(
    'generated client keeps business error code and validation fields',
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
        api.patient('12e575d0-0047-4be4-9e09-07f7e5aa8ffb'),
        throwsA(
          isA<ApiException>()
              .having((e) => e.code, 'code', 'VERSION_CONFLICT')
              .having((e) => e.fields['name'], 'field', 'Неверное значение'),
        ),
      );
      api.close();
    },
  );
  test(
    'business errors without optional fields preserve the error code',
    () async {
      final api = ApiClient(
        client: MockClient(
          (_) async => http.Response(
            '{"error":{"code":"VERSION_CONFLICT","message":"Changed"}}',
            409,
            headers: {'content-type': 'application/json; charset=utf-8'},
          ),
        ),
      );
      await expectLater(
        api.patient('12e575d0-0047-4be4-9e09-07f7e5aa8ffb'),
        throwsA(
          isA<ApiException>()
              .having((e) => e.code, 'code', 'VERSION_CONFLICT')
              .having((e) => e.fields, 'fields', isEmpty),
        ),
      );
      api.close();
    },
  );
}
