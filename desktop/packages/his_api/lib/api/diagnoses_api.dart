//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class DiagnosesApi {
  DiagnosesApi([ApiClient? apiClient])
      : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Создать запись
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [DiagnosisInput] diagnosisInput (required):
  Future<Response> createDiagnosisWithHttpInfo(
    DiagnosisInput diagnosisInput, {
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/diagnoses';

    // ignore: prefer_final_locals
    Object? postBody = diagnosisInput;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];

    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
      abortTrigger: abortTrigger,
    );
  }

  /// Создать запись
  ///
  /// Parameters:
  ///
  /// * [DiagnosisInput] diagnosisInput (required):
  Future<Diagnosis?> createDiagnosis(
    DiagnosisInput diagnosisInput, {
    Future<void>? abortTrigger,
  }) async {
    final response = await createDiagnosisWithHttpInfo(
      diagnosisInput,
      abortTrigger: abortTrigger,
    );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty &&
        response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(
        await _decodeBodyBytes(response),
        'Diagnosis',
      ) as Diagnosis;
    }
    return null;
  }

  /// Мягкое удаление
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  Future<Response> deleteDiagnosisWithHttpInfo(
    String id, {
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/diagnoses/{id}'.replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];

    return apiClient.invokeAPI(
      path,
      'DELETE',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
      abortTrigger: abortTrigger,
    );
  }

  /// Мягкое удаление
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  Future<void> deleteDiagnosis(
    String id, {
    Future<void>? abortTrigger,
  }) async {
    final response = await deleteDiagnosisWithHttpInfo(
      id,
      abortTrigger: abortTrigger,
    );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Карточка, включая мягко удалённые записи
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  Future<Response> getDiagnosisWithHttpInfo(
    String id, {
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/diagnoses/{id}'.replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];

    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
      abortTrigger: abortTrigger,
    );
  }

  /// Карточка, включая мягко удалённые записи
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  Future<Diagnosis?> getDiagnosis(
    String id, {
    Future<void>? abortTrigger,
  }) async {
    final response = await getDiagnosisWithHttpInfo(
      id,
      abortTrigger: abortTrigger,
    );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty &&
        response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(
        await _decodeBodyBytes(response),
        'Diagnosis',
      ) as Diagnosis;
    }
    return null;
  }

  /// Импортировать диагнозы из выбранного JSON-файла одной транзакцией
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [ImportInput] importInput (required):
  Future<Response> importDiagnosesWithHttpInfo(
    ImportInput importInput, {
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/diagnoses/import';

    // ignore: prefer_final_locals
    Object? postBody = importInput;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];

    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
      abortTrigger: abortTrigger,
    );
  }

  /// Импортировать диагнозы из выбранного JSON-файла одной транзакцией
  ///
  /// Parameters:
  ///
  /// * [ImportInput] importInput (required):
  Future<ImportResult?> importDiagnoses(
    ImportInput importInput, {
    Future<void>? abortTrigger,
  }) async {
    final response = await importDiagnosesWithHttpInfo(
      importInput,
      abortTrigger: abortTrigger,
    );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty &&
        response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(
        await _decodeBodyBytes(response),
        'ImportResult',
      ) as ImportResult;
    }
    return null;
  }

  /// Список активных записей
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [int] limit:
  ///
  /// * [int] offset:
  ///
  /// * [String] q:
  ///
  /// * [bool] includeDeleted:
  Future<Response> listDiagnosesWithHttpInfo({
    int? limit,
    int? offset,
    String? q,
    bool? includeDeleted,
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/diagnoses';

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    if (limit != null) {
      queryParams.addAll(_queryParams('', 'limit', limit));
    }
    if (offset != null) {
      queryParams.addAll(_queryParams('', 'offset', offset));
    }
    if (q != null) {
      queryParams.addAll(_queryParams('', 'q', q));
    }
    if (includeDeleted != null) {
      queryParams.addAll(_queryParams('', 'include_deleted', includeDeleted));
    }

    const contentTypes = <String>[];

    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
      abortTrigger: abortTrigger,
    );
  }

  /// Список активных записей
  ///
  /// Parameters:
  ///
  /// * [int] limit:
  ///
  /// * [int] offset:
  ///
  /// * [String] q:
  ///
  /// * [bool] includeDeleted:
  Future<DiagnosisPage?> listDiagnoses({
    int? limit,
    int? offset,
    String? q,
    bool? includeDeleted,
    Future<void>? abortTrigger,
  }) async {
    final response = await listDiagnosesWithHttpInfo(
      limit: limit,
      offset: offset,
      q: q,
      includeDeleted: includeDeleted,
      abortTrigger: abortTrigger,
    );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty &&
        response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(
        await _decodeBodyBytes(response),
        'DiagnosisPage',
      ) as DiagnosisPage;
    }
    return null;
  }

  /// Полностью обновить активную запись
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [DiagnosisInput] diagnosisInput (required):
  Future<Response> updateDiagnosisWithHttpInfo(
    String id,
    DiagnosisInput diagnosisInput, {
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/diagnoses/{id}'.replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody = diagnosisInput;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];

    return apiClient.invokeAPI(
      path,
      'PUT',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
      abortTrigger: abortTrigger,
    );
  }

  /// Полностью обновить активную запись
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [DiagnosisInput] diagnosisInput (required):
  Future<Diagnosis?> updateDiagnosis(
    String id,
    DiagnosisInput diagnosisInput, {
    Future<void>? abortTrigger,
  }) async {
    final response = await updateDiagnosisWithHttpInfo(
      id,
      diagnosisInput,
      abortTrigger: abortTrigger,
    );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty &&
        response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(
        await _decodeBodyBytes(response),
        'Diagnosis',
      ) as Diagnosis;
    }
    return null;
  }
}
