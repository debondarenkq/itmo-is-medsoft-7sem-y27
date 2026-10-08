import 'dart:io';
import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/io_client.dart';
import 'package:his_desktop/core/api_client.dart';
import 'package:his_desktop/core/models.dart';
import 'package:his_desktop/main.dart';

class _RealHttp extends HttpOverrides {}

Future<void> settleNetwork(WidgetTester tester) async {
  for (var index = 0; index < 20; index++) {
    await tester.runAsync(
      () => Future<void>.delayed(const Duration(milliseconds: 30)),
    );
    await tester.pump(const Duration(milliseconds: 50));
  }
  await tester.pumpAndSettle();
}

Future<void> click(WidgetTester tester, Finder finder) async {
  await tester.tap(finder);
  await settleNetwork(tester);
}

Future<void> capture(WidgetTester tester, GlobalKey key, String name) async {
  if (!const bool.fromEnvironment('SAVE_PREVIEWS')) return;
  await tester.runAsync(() async {
    final boundary =
        key.currentContext!.findRenderObject() as RenderRepaintBoundary;
    final image = await boundary.toImage(pixelRatio: 1);
    final bytes = await image.toByteData(format: ui.ImageByteFormat.png);
    final file = File('build/previews/$name.png');
    await file.parent.create(recursive: true);
    await file.writeAsBytes(bytes!.buffer.asUint8List());
    image.dispose();
  });
}

void main() {
  testWidgets(
    'desktop workflow against real Go services through the gateway',
    (tester) async {
      tester.view.physicalSize = const Size(1440, 960);
      tester.view.devicePixelRatio = 1;
      addTearDown(tester.view.resetPhysicalSize);
      addTearDown(tester.view.resetDevicePixelRatio);
      for (final font in ['Regular', 'Medium', 'Bold']) {
        final loader = FontLoader('HospitalSans')
          ..addFont(rootBundle.load('assets/fonts/Roboto-$font.ttf'));
        await loader.load();
      }
      final icons = FontLoader('MaterialIcons')
        ..addFont(rootBundle.load('fonts/MaterialIcons-Regular.otf'));
      await icons.load();
      final api = ApiClient(
        client: IOClient(_RealHttp().createHttpClient(null)),
      );
      final data = await tester.runAsync(() async {
        final suffix = DateTime.now().microsecondsSinceEpoch
            .toRadixString(16)
            .toUpperCase()
            .substring(8);
        final staff = await api.request(
          'POST',
          '/api/v1/staff',
          body: {
            'first_name': 'Анна',
            'last_name': 'Петрова$suffix',
            'position': 'Терапевт',
          },
        );
        final diagnosis = await api.request(
          'POST',
          '/api/v1/diagnoses',
          body: {
            'code': suffix.padLeft(6, '0').substring(0, 6),
            'name': 'Учебный диагноз $suffix',
          },
        );
        final patient = await api.request(
          'POST',
          '/api/v1/patients',
          body: {
            'first_name': 'Иван',
            'last_name': 'Иванов$suffix',
            'middle_name': 'Иванович',
            'birth_date': '1994-05-12',
            'administrative_sex': 'M',
            'comment': 'Учебная карточка',
          },
        );
        return {'staff': staff, 'patient': patient, 'diagnosis': diagnosis};
      });
      final staffName = '${data!['staff']!['last_name']} Анна';
      final patientName = '${data['patient']!['last_name']} Иван Иванович';
      final diagnosisName = data['diagnosis']!['name'] as String;
      final previewKey = GlobalKey();
      await tester.pumpWidget(
        RepaintBoundary(
          key: previewKey,
          child: HisApp(api: api),
        ),
      );
      await settleNetwork(tester);
      expect(find.text('Анна'), findsNothing);
      expect(find.text(staffName), findsOneWidget);
      await capture(tester, previewKey, '01-staff-list');
      await click(tester, find.text(staffName));
      expect(find.text('Карточка медработника'), findsOneWidget);
      await capture(tester, previewKey, '02-staff-card');
      await click(tester, find.text('Отмена'));
      await click(tester, find.text('Диагнозы'));
      expect(find.text('Справочник диагнозов'), findsOneWidget);
      await capture(tester, previewKey, '03-diagnoses');

      await click(
        tester,
        find.byWidgetPredicate((w) => w is DropdownButtonFormField<bool>),
      );
      await click(tester, find.text('Медработник').last);
      expect(find.text(patientName), findsOneWidget);
      await capture(tester, previewKey, '04-patients');
      await click(tester, find.text(patientName));
      expect(find.text('Карточка пациента'), findsOneWidget);
      await capture(tester, previewKey, '05-patient-card');
      await click(tester, find.text('Создать ЭМК'));
      await click(
        tester,
        find.byWidgetPredicate((w) => w is DropdownButtonFormField<Staff>),
      );
      await click(tester, find.text('$staffName · Терапевт').last);
      await click(tester, find.widgetWithText(FilledButton, 'Создать ЭМК'));
      expect(find.text('Только просмотр'), findsOneWidget);
      expect(
        tester
            .widget<FilledButton>(
              find.widgetWithText(FilledButton, 'Добавить диагноз'),
            )
            .onPressed,
        isNull,
      );
      await capture(tester, previewKey, '06-record-locked');
      await click(
        tester,
        find.widgetWithText(OutlinedButton, 'Разблокировать'),
      );
      await click(tester, find.widgetWithText(FilledButton, 'Разблокировать'));
      await click(
        tester,
        find.widgetWithText(FilledButton, 'Добавить диагноз'),
      );
      await tester.enterText(find.byType(TextField), diagnosisName);
      await tester.pump();
      await click(tester, find.widgetWithText(ListTile, diagnosisName));
      await click(tester, find.widgetWithText(Tab, 'Назначения'));
      await click(tester, find.text('Добавить назначение'));
      await tester.enterText(find.byType(TextFormField), 'Пить больше воды');
      await click(tester, find.text('Добавить в карту'));
      await capture(tester, previewKey, '07-record-editing');
      await click(tester, find.text('Сохранить изменения'));
      expect(find.text('Только просмотр'), findsOneWidget);
      final patientId = (data['patient'] as Json)['id'] as String;
      final saved = await tester.runAsync(() => api.patientRecord(patientId));
      expect(saved!.version, 2);
      expect(saved.diagnoses.length, 1);
      expect(saved.prescriptions.single.text, 'Пить больше воды');

      // Edit, save and cancel a persisted prescription through the UI.
      await click(
        tester,
        find.widgetWithText(OutlinedButton, 'Разблокировать'),
      );
      await click(tester, find.widgetWithText(FilledButton, 'Разблокировать'));
      await click(tester, find.widgetWithText(Tab, 'Назначения'));
      await click(tester, find.byTooltip('Редактировать назначение'));
      await tester.enterText(find.byType(TextFormField), 'Пить воду после еды');
      await click(tester, find.text('Добавить в карту'));
      await click(tester, find.text('Сохранить изменения'));
      await click(
        tester,
        find.widgetWithText(OutlinedButton, 'Разблокировать'),
      );
      await click(tester, find.widgetWithText(FilledButton, 'Разблокировать'));
      await click(tester, find.widgetWithText(Tab, 'Назначения'));
      await click(tester, find.byTooltip('Отменить назначение'));
      await click(
        tester,
        find.widgetWithText(FilledButton, 'Отменить назначение'),
      );
      await click(tester, find.text('Сохранить изменения'));
      expect(find.text('Отменено'), findsOneWidget);
      await click(tester, find.widgetWithText(Tab, 'История изменений'));
      expect(find.text('Отменено назначение'), findsOneWidget);
      await capture(tester, previewKey, '08-record-history');
      final history = await tester.runAsync(() => api.history(saved.id));
      expect(history!.where((e) => e.version == 3).map((e) => e.type), [
        'prescription_removed',
        'prescription_added',
      ]);

      // Historical state is read-only and does not change current EMR data.
      await tester.scrollUntilVisible(
        find.text('Добавлено назначение').last,
        300,
        scrollable: find
            .descendant(
              of: find.byType(TabBarView),
              matching: find.byType(Scrollable),
            )
            .last,
      );
      await click(tester, find.text('Добавлено назначение').last);
      await Scrollable.ensureVisible(
        tester.element(find.text('Посмотреть карту · версия 2')),
        alignment: .25,
      );
      await tester.pumpAndSettle();
      await click(tester, find.text('Посмотреть карту · версия 2'));
      expect(find.text('Историческое состояние · версия 2'), findsOneWidget);
      await click(tester, find.widgetWithText(Tab, 'Назначения'));
      expect(find.text('Пить больше воды'), findsOneWidget);
      await capture(tester, previewKey, '09-historical-state');
      await click(tester, find.text('К текущей карте'));
      expect(find.text('Пить воду после еды'), findsOneWidget);
      expect(find.text('Отменено'), findsOneWidget);
      // Removing and restoring diagnoses preserves the history and creates a new entry.
      await click(
        tester,
        find.widgetWithText(OutlinedButton, 'Разблокировать'),
      );
      await click(tester, find.widgetWithText(FilledButton, 'Разблокировать'));
      await click(tester, find.widgetWithText(Tab, 'Диагнозы'));
      await click(tester, find.byTooltip('Удалить диагноз из ЭМК'));
      await click(tester, find.text('К карточке пациента'));
      expect(find.text('Отменить несохранённые изменения?'), findsOneWidget);
      await click(tester, find.widgetWithText(TextButton, 'Отмена'));
      expect(find.text('Редактирование карты'), findsOneWidget);
      await click(tester, find.text('Сохранить изменения'));
      await click(
        tester,
        find.widgetWithText(OutlinedButton, 'Разблокировать'),
      );
      await click(tester, find.widgetWithText(FilledButton, 'Разблокировать'));
      await click(tester, find.byTooltip('Восстановить диагноз'));
      await click(tester, find.text('Сохранить изменения'));
      final restored = await tester.runAsync(
        () => api.patientRecord(patientId),
      );
      expect(restored!.diagnoses.single.id, isNot(saved.diagnoses.single.id));

      // A concurrent save never silently overwrites another user's changes.
      await click(
        tester,
        find.widgetWithText(OutlinedButton, 'Разблокировать'),
      );
      await click(tester, find.widgetWithText(FilledButton, 'Разблокировать'));
      await click(tester, find.widgetWithText(Tab, 'Назначения'));
      await click(tester, find.text('Добавить назначение'));
      await tester.enterText(
        find.byType(TextFormField),
        'Высыпаться ежедневно',
      );
      await click(tester, find.text('Добавить в карту'));
      await tester.runAsync(
        () => api.saveRecord(restored, data['staff']!['id'] as String, [
          {'type': 'add_prescription', 'text': 'Параллельное назначение'},
        ]),
      );
      await click(tester, find.text('Сохранить изменения'));
      expect(find.text('Высыпаться ежедневно'), findsOneWidget);
      expect(find.textContaining('Черновик сохранён'), findsOneWidget);
      await capture(tester, previewKey, '10-version-conflict');
      await click(
        tester,
        find.widgetWithText(OutlinedButton, 'Обновить карту'),
      );
      await click(tester, find.widgetWithText(FilledButton, 'Обновить карту'));
      expect(find.text('Только просмотр'), findsOneWidget);
      expect(find.text('Высыпаться ежедневно'), findsNothing);
      expect(find.text('Параллельное назначение'), findsOneWidget);
      expect(tester.takeException(), isNull);
    },
    skip: !const bool.fromEnvironment('HIS_INTEGRATION'),
    timeout: const Timeout(Duration(minutes: 4)),
  );
}
