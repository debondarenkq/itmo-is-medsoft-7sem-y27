import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'core/api_client.dart';
import 'core/session.dart';
import 'core/theme.dart';
import 'screens/shell.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const HisApp());
}

class HisApp extends StatefulWidget {
  const HisApp({super.key, this.api});
  final ApiClient? api;
  @override
  State<HisApp> createState() => _HisAppState();
}

class _HisAppState extends State<HisApp> {
  late final AppSession _session;
  @override
  void initState() {
    super.initState();
    _session = AppSession(widget.api ?? ApiClient());
  }

  @override
  void dispose() {
    _session.api.close();
    _session.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => MaterialApp(
    title: 'МедКарта — HIS',
    debugShowCheckedModeBanner: false,
    theme: medicalTheme(),
    locale: const Locale('ru'),
    supportedLocales: const [Locale('ru')],
    localizationsDelegates: GlobalMaterialLocalizations.delegates,
    home: ApplicationShell(session: _session),
  );
}
