//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class RecordsApi {
  RecordsApi([ApiClient? apiClient])
      : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Создать ЭМК с обязательным автором
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [CreateRecordInput] createRecordInput (required):
  Future<Response> createRecordWithHttpInfo(
    String id,
    CreateRecordInput createRecordInput, {
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/patients/{id}/record'.replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody = createRecordInput;

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

  /// Создать ЭМК с обязательным автором
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [CreateRecordInput] createRecordInput (required):
  Future<Record?> createRecord(
    String id,
    CreateRecordInput createRecordInput, {
    Future<void>? abortTrigger,
  }) async {
    final response = await createRecordWithHttpInfo(
      id,
      createRecordInput,
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
        'Record',
      ) as Record;
    }
    return null;
  }

  /// Текущее состояние карты
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  Future<Response> getRecordWithHttpInfo(
    String id, {
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/records/{id}'.replaceAll('{id}', id);

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

  /// Текущее состояние карты
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  Future<Record?> getRecord(
    String id, {
    Future<void>? abortTrigger,
  }) async {
    final response = await getRecordWithHttpInfo(
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
        'Record',
      ) as Record;
    }
    return null;
  }

  /// История по порядку событий
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [int] limit:
  ///
  /// * [int] offset:
  Future<Response> getRecordHistoryWithHttpInfo(
    String id, {
    int? limit,
    int? offset,
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/records/{id}/history'.replaceAll('{id}', id);

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

  /// История по порядку событий
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [int] limit:
  ///
  /// * [int] offset:
  Future<EventPage?> getRecordHistory(
    String id, {
    int? limit,
    int? offset,
    Future<void>? abortTrigger,
  }) async {
    final response = await getRecordHistoryWithHttpInfo(
      id,
      limit: limit,
      offset: offset,
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
        'EventPage',
      ) as EventPage;
    }
    return null;
  }

  /// Восстановить карту; требуется ровно один из at или version
  ///
  /// Требуется ровно один параметр at или version. Несуществующая версия и момент до создания карты возвращают 404.
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [DateTime] at:
  ///
  /// * [int] version:
  Future<Response> getRecordStateWithHttpInfo(
    String id, {
    DateTime? at,
    int? version,
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/records/{id}/state'.replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    if (at != null) {
      queryParams.addAll(_queryParams('', 'at', at));
    }
    if (version != null) {
      queryParams.addAll(_queryParams('', 'version', version));
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

  /// Восстановить карту; требуется ровно один из at или version
  ///
  /// Требуется ровно один параметр at или version. Несуществующая версия и момент до создания карты возвращают 404.
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [DateTime] at:
  ///
  /// * [int] version:
  Future<Record?> getRecordState(
    String id, {
    DateTime? at,
    int? version,
    Future<void>? abortTrigger,
  }) async {
    final response = await getRecordStateWithHttpInfo(
      id,
      at: at,
      version: version,
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
        'Record',
      ) as Record;
    }
    return null;
  }

  /// Атомарно сохранить команды и историю
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [ChangesInput] changesInput (required):
  Future<Response> saveRecordWithHttpInfo(
    String id,
    ChangesInput changesInput, {
    Future<void>? abortTrigger,
  }) async {
    // ignore: prefer_const_declarations
    final path = r'/api/v1/records/{id}/changes'.replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody = changesInput;

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

  /// Атомарно сохранить команды и историю
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///
  /// * [ChangesInput] changesInput (required):
  Future<Record?> saveRecord(
    String id,
    ChangesInput changesInput, {
    Future<void>? abortTrigger,
  }) async {
    final response = await saveRecordWithHttpInfo(
      id,
      changesInput,
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
        'Record',
      ) as Record;
    }
    return null;
  }
}
