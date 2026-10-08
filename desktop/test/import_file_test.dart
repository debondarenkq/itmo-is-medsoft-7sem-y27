import 'dart:convert';
import 'package:file_selector_platform_interface/file_selector_platform_interface.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:his_desktop/core/api_client.dart';
import 'package:his_desktop/main.dart';

class _FilePicker extends FileSelectorPlatform {
  _FilePicker(this.file);
  final XFile? file;
  @override
  Future<XFile?> openFile({
    List<XTypeGroup>? acceptedTypeGroups,
    String? initialDirectory,
    String? confirmButtonText,
  }) async => file;
}

void main() {
  testWidgets(
    'selected JSON file is sent to the gateway without automatic seeding',
    (tester) async {
      tester.view.physicalSize = const Size(1280, 840);
      tester.view.devicePixelRatio = 1;
      addTearDown(tester.view.resetPhysicalSize);
      addTearDown(tester.view.resetDevicePixelRatio);
      final original = FileSelectorPlatform.instance;
      addTearDown(() => FileSelectorPlatform.instance = original);
      const contents =
          '{"diagnoses":[{"code":"TEST01","name":"Первый диагноз"}]}';
      FileSelectorPlatform.instance = _FilePicker(
        XFile.fromData(utf8.encode(contents), name: 'diagnoses.json'),
      );
      final posted = <http.Request>[];
      final api = ApiClient(
        client: MockClient((request) async {
          if (request.method == 'POST') {
            posted.add(request);
            return http.Response('{"items":[{}]}', 200);
          }
          return http.Response('{"items":[],"total":0}', 200);
        }),
      );
      await tester.pumpWidget(HisApp(api: api));
      await tester.pumpAndSettle();
      await tester.tap(find.text('Диагнозы'));
      await tester.pumpAndSettle();
      expect(posted, isEmpty);
      await tester.tap(find.text('Загрузить из файла'));
      await tester.pumpAndSettle();
      expect(posted.single.url.path, '/api/v1/diagnoses/import');
      expect(posted.single.body, contents);
      expect(posted.single.headers['Content-Type'], 'application/json');
      expect(
        find.text('Импорт завершён. Загружено записей: 1'),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    },
  );
}
