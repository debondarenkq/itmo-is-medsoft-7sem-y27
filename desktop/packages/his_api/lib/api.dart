//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

library openapi.api;

import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:collection/collection.dart';
import 'package:http/http.dart';
import 'package:intl/intl.dart';
import 'package:meta/meta.dart';

part 'api_client.dart';
part 'api_helper.dart';
part 'api_exception.dart';
part 'auth/authentication.dart';
part 'auth/api_key_auth.dart';
part 'auth/oauth.dart';
part 'auth/http_basic_auth.dart';
part 'auth/http_bearer_auth.dart';
part 'optional.dart';

part 'api/diagnoses_api.dart';
part 'api/patients_api.dart';
part 'api/records_api.dart';
part 'api/staff_api.dart';

part 'model/actor.dart';
part 'model/administrative_sex.dart';
part 'model/changes_input.dart';
part 'model/command.dart';
part 'model/command_type.dart';
part 'model/create_record_input.dart';
part 'model/diagnosis.dart';
part 'model/diagnosis_input.dart';
part 'model/diagnosis_page.dart';
part 'model/error.dart';
part 'model/error_detail.dart';
part 'model/event.dart';
part 'model/event_page.dart';
part 'model/event_type.dart';
part 'model/import_input.dart';
part 'model/import_result.dart';
part 'model/patient.dart';
part 'model/patient_input.dart';
part 'model/patient_page.dart';
part 'model/prescription.dart';
part 'model/prescription_status.dart';
part 'model/record.dart';
part 'model/record_diagnosis.dart';
part 'model/staff.dart';
part 'model/staff_input.dart';
part 'model/staff_page.dart';
part 'model/state.dart';

/// An [ApiClient] instance that uses the default values obtained from
/// the OpenAPI specification file.
var defaultApiClient = ApiClient();

const _delimiters = {'csv': ',', 'ssv': ' ', 'tsv': '\t', 'pipes': '|'};
const _dateEpochMarker = 'epoch';
const _deepEquality = DeepCollectionEquality();
final _dateFormatter = DateFormat('yyyy-MM-dd');
final _regList = RegExp(r'^List<(.*)>$');
final _regSet = RegExp(r'^Set<(.*)>$');
final _regMap = RegExp(r'^Map<String,(.*)>$');

bool _isEpochMarker(String? pattern) =>
    pattern == _dateEpochMarker || pattern == '/$_dateEpochMarker/';
