import 'dart:async';
import 'dart:convert';
import 'package:his_api/api.dart' as contract;
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

// One generated gateway client serves every public resource group.
class ApiClient {
  ApiClient({
    String baseUrl = const String.fromEnvironment(
      'API_BASE_URL',
      defaultValue: 'http://localhost:8080',
    ),
    http.Client? client,
  }) : _http = client ?? http.Client() {
    configure(baseUrl);
  }
  final http.Client _http;
  late contract.ApiClient _gateway;
  late contract.StaffApi _staff;
  late contract.DiagnosesApi _diagnoses;
  late contract.PatientsApi _patients;
  late contract.RecordsApi _records;
  String get baseUrl => _gateway.basePath;
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
    final address = uri
        .replace(path: uri.path.replaceFirst(RegExp(r'/+$'), ''))
        .toString();
    _gateway = contract.ApiClient(basePath: address);
    _gateway.client.close();
    _gateway.client = _http;
    _staff = contract.StaffApi(_gateway);
    _diagnoses = contract.DiagnosesApi(_gateway);
    _patients = contract.PatientsApi(_gateway);
    _records = contract.RecordsApi(_gateway);
  }

  void close() => _http.close();

  Future<T> _request<T>(Future<T> Function() operation) async {
    try {
      return await operation().timeout(const Duration(seconds: 20));
    } on TimeoutException {
      throw const ApiException(
        'Сервер не ответил вовремя. Попробуйте ещё раз.',
      );
    } on contract.ApiException catch (error) {
      if (error.code == 0 ||
          error.code == 400 && error.innerException != null) {
        throw const ApiException(
          'Нет соединения с сервером. Проверьте подключение и адрес сервера.',
        );
      }
      try {
        final detail = contract.Error.fromJson(
          jsonDecode(error.message ?? ''),
        )?.error;
        if (detail != null) {
          throw ApiException(
            detail.message,
            status: error.code,
            code: detail.code,
            fields: detail.fields.isPresent ? detail.fields.value ?? {} : {},
          );
        }
      } on FormatException {
        /* A proxy may return an HTML error instead of JSON. */
      }
      throw ApiException(
        'Сервер временно недоступен. Повторите запрос.',
        status: error.code,
      );
    } on http.ClientException {
      throw const ApiException('Не удалось подключиться к серверу.');
    }
  }

  T _present<T>(T? value) {
    if (value == null) {
      throw const ApiException('Сервер вернул неполный ответ.');
    }
    return value;
  }

  Future<PageResult<Staff>> staffPage({
    int offset = 0,
    String query = '',
    bool includeDeleted = false,
    int limit = 50,
  }) async {
    final page = _present(
      await _request(
        () => _staff.listStaff(
          limit: limit,
          offset: offset,
          q: query,
          includeDeleted: includeDeleted,
        ),
      ),
    );
    return PageResult(page.items.map(Staff.fromContract).toList(), page.total);
  }

  Future<PageResult<Diagnosis>> diagnosisPage({
    int offset = 0,
    String query = '',
    int limit = 50,
  }) async {
    final page = _present(
      await _request(
        () => _diagnoses.listDiagnoses(limit: limit, offset: offset, q: query),
      ),
    );
    return PageResult(
      page.items.map(Diagnosis.fromContract).toList(),
      page.total,
    );
  }

  Future<PageResult<Patient>> patientPage({
    int offset = 0,
    String query = '',
    int limit = 50,
  }) async {
    final page = _present(
      await _request(
        () => _patients.listPatients(limit: limit, offset: offset, q: query),
      ),
    );
    return PageResult(
      page.items.map(Patient.fromContract).toList(),
      page.total,
    );
  }

  Future<List<T>> _all<T>(
    Future<PageResult<T>> Function(int offset) fetch,
  ) async {
    final items = <T>[];
    while (true) {
      final page = await fetch(items.length);
      items.addAll(page.items);
      if (page.items.isEmpty || items.length >= page.total) return items;
    }
  }

  Future<List<Staff>> staff() =>
      _all((offset) => staffPage(offset: offset, limit: 200));
  Future<List<Diagnosis>> diagnoses() =>
      _all((offset) => diagnosisPage(offset: offset, limit: 200));
  Future<Staff> createStaff(contract.StaffInput input) async =>
      Staff.fromContract(
        _present(await _request(() => _staff.createStaff(input))),
      );
  Future<Staff> updateStaff(String id, contract.StaffInput input) async =>
      Staff.fromContract(
        _present(await _request(() => _staff.updateStaff(id, input))),
      );
  Future<void> deleteStaff(String id) => _request(() => _staff.deleteStaff(id));
  Future<Diagnosis> createDiagnosis(contract.DiagnosisInput input) async =>
      Diagnosis.fromContract(
        _present(await _request(() => _diagnoses.createDiagnosis(input))),
      );
  Future<Diagnosis> updateDiagnosis(
    String id,
    contract.DiagnosisInput input,
  ) async => Diagnosis.fromContract(
    _present(await _request(() => _diagnoses.updateDiagnosis(id, input))),
  );
  Future<void> deleteDiagnosis(String id) =>
      _request(() => _diagnoses.deleteDiagnosis(id));
  Future<int> importDiagnoses(String contents) async {
    dynamic document;
    try {
      document = jsonDecode(contents.replaceFirst(RegExp('^\uFEFF'), ''));
    } on FormatException {
      throw const FormatException('Файл должен содержать корректный JSON.');
    }
    if (document is! Json) {
      throw const FormatException(
        'В файле ожидается JSON-объект со списком diagnoses.',
      );
    }
    final result = _present(
      await _request(
        () => _diagnoses.importDiagnoses(_ImportDocument(document)),
      ),
    );
    return result.items.length;
  }

  Future<Patient> createPatient(contract.PatientInput input) async =>
      Patient.fromContract(
        _present(await _request(() => _patients.createPatient(input))),
      );
  Future<Patient> updatePatient(String id, contract.PatientInput input) async =>
      Patient.fromContract(
        _present(await _request(() => _patients.updatePatient(id, input))),
      );
  Future<void> deletePatient(String id) =>
      _request(() => _patients.deletePatient(id));
  Future<Patient> patient(String id) async => Patient.fromContract(
    _present(await _request(() => _patients.getPatient(id))),
  );
  Future<MedicalRecord> patientRecord(String id) async =>
      MedicalRecord.fromContract(
        _present(await _request(() => _patients.getPatientRecord(id))),
      );
  Future<MedicalRecord> createRecord(String id, String staffId) async =>
      MedicalRecord.fromContract(
        _present(
          await _request(
            () => _records.createRecord(
              id,
              contract.CreateRecordInput(staffId: staffId),
            ),
          ),
        ),
      );
  Future<List<RecordEvent>> history(String id) => _all((offset) async {
    final page = _present(
      await _request(
        () => _records.getRecordHistory(id, limit: 200, offset: offset),
      ),
    );
    return PageResult(
      page.items.map(RecordEvent.fromContract).toList(),
      page.total,
    );
  });
  Future<MedicalRecord> saveRecord(
    MedicalRecord base,
    String staffId,
    List<contract.Command> commands,
  ) async => MedicalRecord.fromContract(
    _present(
      await _request(
        () => _records.saveRecord(
          base.id,
          contract.ChangesInput(
            staffId: staffId,
            expectedVersion: base.version,
            commands: commands,
          ),
        ),
      ),
    ),
  );
  Future<MedicalRecord> historicalRecord(
    String id, {
    int? version,
    DateTime? at,
  }) async => MedicalRecord.fromContract(
    _present(
      await _request(
        () => _records.getRecordState(id, version: version, at: at?.toUtc()),
      ),
    ),
  );
}

// Preserve all file fields for server-side validation, including unknown fields.
class _ImportDocument extends contract.ImportInput {
  _ImportDocument(this.document) : super(diagnoses: []);
  final Json document;
  @override
  Json toJson() => document;
}
