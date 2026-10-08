import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:his_desktop/core/api_client.dart';
import 'package:his_desktop/main.dart';

void main() {
  testWidgets(
    'empty hospital starts with empty lists and real creation controls',
    (tester) async {
      tester.view.physicalSize = const Size(1280, 840);
      tester.view.devicePixelRatio = 1;
      addTearDown(tester.view.resetPhysicalSize);
      addTearDown(tester.view.resetDevicePixelRatio);
      final api = ApiClient(
        client: MockClient(
          (_) async => http.Response(
            '{"items":[],"total":0,"limit":50,"offset":0}',
            200,
          ),
        ),
      );
      await tester.pumpWidget(HisApp(api: api));
      await tester.pumpAndSettle();
      expect(find.text('Добавьте первого медработника'), findsOneWidget);
      await tester.tap(find.text('Добавить медработника'));
      await tester.pumpAndSettle();
      expect(find.text('Новый медработник'), findsOneWidget);
      await tester.tap(find.text('Сохранить'));
      await tester.pumpAndSettle();
      expect(find.text('От 3 до 24 символов'), findsNWidgets(3));
      await tester.tap(find.text('Отмена'));
      await tester.pumpAndSettle();
      await tester.tap(find.text('Диагнозы'));
      await tester.pumpAndSettle();
      expect(find.text('Справочник пока пуст'), findsOneWidget);
      expect(find.text('Загрузить из файла'), findsOneWidget);
      expect(tester.takeException(), isNull);
    },
  );
  testWidgets('connection failure is actionable and does not add mock data', (
    tester,
  ) async {
    tester.view.physicalSize = const Size(1080, 760);
    tester.view.devicePixelRatio = 1;
    addTearDown(tester.view.resetPhysicalSize);
    addTearDown(tester.view.resetDevicePixelRatio);
    final api = ApiClient(
      client: MockClient((_) async => throw http.ClientException('offline')),
    );
    await tester.pumpWidget(HisApp(api: api));
    await tester.pumpAndSettle();
    expect(find.text('Не удалось загрузить данные'), findsOneWidget);
    expect(find.text('Попробовать снова'), findsOneWidget);
    expect(tester.takeException(), isNull);
  });
}
