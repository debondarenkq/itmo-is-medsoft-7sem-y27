import 'package:flutter/foundation.dart';
import 'api_client.dart';
import 'models.dart';

class AppSession extends ChangeNotifier {
  AppSession(this.api);
  final ApiClient api;
  Staff? worker;
  Future<bool> Function()? beforeNavigation;
  void selectWorker(Staff? value) {
    worker = value;
    notifyListeners();
  }

  Future<bool> canNavigate() async => await beforeNavigation?.call() ?? true;
}
