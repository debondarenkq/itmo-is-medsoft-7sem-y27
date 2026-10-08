import 'package:flutter/material.dart';
import '../core/api_client.dart';
import '../core/models.dart';
import '../core/session.dart';
import '../core/theme.dart';
import '../widgets/common.dart';
import '../widgets/paged_table.dart';
import '../widgets/staff_picker.dart';
import 'record_screen.dart';

class PatientListScreen extends StatefulWidget {
  const PatientListScreen({super.key, required this.session});
  final AppSession session;
  @override
  State<PatientListScreen> createState() => _PatientListScreenState();
}

class _PatientListScreenState extends State<PatientListScreen> {
  final _table = GlobalKey<PagedTableState<Patient>>();
  Future<void> _open({Patient? patient, bool edit = false}) async {
    await Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => PatientFormScreen(
          session: widget.session,
          patient: patient,
          edit: edit,
        ),
      ),
    );
    _table.currentState?.reload();
  }

  Future<void> _delete(Patient patient) async {
    if (!await confirm(
      context,
      title: 'Удалить пациента?',
      message: '${patient.name} будет полностью удалён из системы.',
    )) {
      return;
    }
    try {
      await widget.session.api.request(
        'DELETE',
        '/api/v1/patients/${patient.id}',
      );
      if (mounted) {
        notify(context, 'Пациент удалён');
        _table.currentState?.reload();
      }
    } catch (error) {
      if (mounted) notify(context, errorMessage(error), error: true);
    }
  }

  @override
  Widget build(BuildContext context) => PageFrame(
    title: 'Пациенты',
    subtitle: 'Карточки пациентов и электронные медицинские карты',
    actions: [
      FilledButton.icon(
        onPressed: () => _open(edit: true),
        icon: const Icon(Icons.add_rounded, size: 20),
        label: const Text('Добавить пациента'),
      ),
    ],
    child: PagedTable<Patient>(
      key: _table,
      fetch: (offset, query) => widget.session.api.page(
        '/api/v1/patients',
        Patient.fromJson,
        offset: offset,
        query: query,
      ),
      searchHint: 'Поиск по фамилии, имени или отчеству',
      emptyTitle: 'Пациентов пока нет',
      emptySubtitle: 'Создайте карточку пациента, чтобы начать работу с ЭМК.',
      columns: const [
        DataColumn(label: Text('ФИО')),
        DataColumn(label: Text('ДАТА РОЖДЕНИЯ')),
        DataColumn(label: Text('ПОЛ')),
        DataColumn(label: Text('МЕДИЦИНСКАЯ КАРТА')),
        DataColumn(label: Text('ДЕЙСТВИЯ')),
      ],
      row: (p) => DataRow(
        onSelectChanged: (_) => _open(patient: p),
        cells: [
          DataCell(
            Text(p.name, style: const TextStyle(fontWeight: FontWeight.w600)),
          ),
          DataCell(Text(displayDate(p.birthDate))),
          DataCell(Text(p.sexLabel)),
          DataCell(
            StatusPill(
              p.hasRecord ? 'ЭМК есть' : 'Нет ЭМК',
              tone: p.hasRecord ? MedicalColors.teal : MedicalColors.muted,
            ),
          ),
          DataCell(
            Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                IconButton(
                  tooltip: 'Посмотреть пациента',
                  onPressed: () => _open(patient: p),
                  icon: const Icon(Icons.visibility_outlined, size: 19),
                ),
                IconButton(
                  tooltip: 'Редактировать пациента',
                  onPressed: () => _open(patient: p, edit: true),
                  icon: const Icon(Icons.edit_outlined, size: 19),
                ),
                IconButton(
                  tooltip: p.hasRecord
                      ? 'Удаление запрещено: у пациента есть ЭМК'
                      : 'Удалить пациента',
                  onPressed: p.hasRecord ? null : () => _delete(p),
                  icon: Icon(
                    Icons.delete_outline_rounded,
                    size: 19,
                    color: p.hasRecord
                        ? MedicalColors.muted
                        : MedicalColors.red,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    ),
  );
}

class PatientFormScreen extends StatefulWidget {
  const PatientFormScreen({
    super.key,
    required this.session,
    this.patient,
    this.edit = false,
  });
  final AppSession session;
  final Patient? patient;
  final bool edit;
  @override
  State<PatientFormScreen> createState() => _PatientFormScreenState();
}

class _PatientFormScreenState extends State<PatientFormScreen> {
  final _form = GlobalKey<FormState>();
  late final TextEditingController _first, _last, _middle, _comment, _birth;
  Patient? _patient;
  DateTime? _birthDate;
  String _sex = 'UNKNOWN';
  bool _readOnly = false, _saving = false, _opening = false;
  Object? _error;
  bool get _busy => _saving || _opening;
  @override
  void initState() {
    super.initState();
    _patient = widget.patient;
    _readOnly = _patient != null && !widget.edit;
    _first = TextEditingController(text: _patient?.firstName);
    _last = TextEditingController(text: _patient?.lastName);
    _middle = TextEditingController(text: _patient?.middleName);
    _comment = TextEditingController(text: _patient?.comment);
    _birthDate = _patient?.birthDate;
    _birth = TextEditingController(
      text: _birthDate == null ? '' : displayDate(_birthDate!),
    );
    _sex = _patient?.sex ?? 'UNKNOWN';
  }

  @override
  void dispose() {
    _first.dispose();
    _last.dispose();
    _middle.dispose();
    _comment.dispose();
    _birth.dispose();
    super.dispose();
  }

  Future<void> _pickDate() async {
    final now = DateTime.now();
    final value = await showDatePicker(
      context: context,
      initialDate: _birthDate ?? DateTime(2000, 1, 1),
      firstDate: DateTime(1),
      lastDate: DateTime(now.year, now.month, now.day),
      helpText: 'Дата рождения',
    );
    if (value != null && mounted) {
      setState(() {
        _birthDate = value;
        _birth.text = displayDate(value);
      });
    }
  }

  Future<void> _save() async {
    if (!_form.currentState!.validate()) return;
    setState(() {
      _saving = true;
      _error = null;
    });
    try {
      final data = await widget.session.api.request(
        _patient == null ? 'POST' : 'PUT',
        '/api/v1/patients${_patient == null ? '' : '/${_patient!.id}'}',
        body: {
          'first_name': _first.text.trim(),
          'last_name': _last.text.trim(),
          'middle_name': _middle.text.trim().isEmpty
              ? null
              : _middle.text.trim(),
          'birth_date': calendarDate(_birthDate!),
          'administrative_sex': _sex,
          'comment': _comment.text.trim(),
        },
      );
      if (!mounted) return;
      setState(() {
        _patient = Patient.fromJson(data);
        _readOnly = true;
      });
      notify(context, 'Карточка пациента сохранена');
    } catch (error) {
      if (mounted) setState(() => _error = error);
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _delete() async {
    final patient = _patient!;
    if (!await confirm(
      context,
      title: 'Удалить пациента?',
      message: '${patient.name} будет полностью удалён из системы.',
    )) {
      return;
    }
    setState(() => _saving = true);
    try {
      await widget.session.api.request(
        'DELETE',
        '/api/v1/patients/${patient.id}',
      );
      if (mounted) {
        notify(context, 'Пациент удалён');
        Navigator.pop(context);
      }
    } catch (error) {
      if (mounted) setState(() => _error = error);
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _openRecord() async {
    setState(() {
      _opening = true;
      _error = null;
    });
    try {
      final patient = await widget.session.api.patient(_patient!.id);
      MedicalRecord record;
      if (patient.hasRecord) {
        record = await widget.session.api.patientRecord(patient.id);
      } else {
        if (!mounted) return;
        final actor = await selectStaff(
          context,
          widget.session,
          title: 'Создать медицинскую карту',
          action: 'Создать ЭМК',
        );
        if (actor == null) return;
        try {
          record = await widget.session.api.createRecord(patient.id, actor.id);
        } on ApiException catch (error) {
          if (error.code != 'RECORD_EXISTS') rethrow;
          record = await widget.session.api.patientRecord(patient.id);
        }
      }
      if (!mounted) return;
      await Navigator.of(context).push(
        MaterialPageRoute<void>(
          builder: (_) => RecordScreen(
            session: widget.session,
            patient: patient,
            record: record,
          ),
        ),
      );
      final fresh = await widget.session.api.patient(patient.id);
      if (mounted) setState(() => _patient = fresh);
    } catch (error) {
      if (mounted) setState(() => _error = error);
    } finally {
      if (mounted) setState(() => _opening = false);
    }
  }

  void _cancelEdit() {
    if (_patient == null) {
      Navigator.pop(context);
      return;
    }
    final p = _patient!;
    setState(() {
      _first.text = p.firstName;
      _last.text = p.lastName;
      _middle.text = p.middleName ?? '';
      _comment.text = p.comment;
      _birthDate = p.birthDate;
      _birth.text = displayDate(p.birthDate);
      _sex = p.sex;
      _readOnly = true;
      _error = null;
    });
  }

  @override
  Widget build(BuildContext context) => PageFrame(
    title: _patient == null ? 'Новый пациент' : 'Карточка пациента',
    subtitle: _patient?.name ?? 'Заполните данные для создания карточки',
    back: _busy ? null : () => Navigator.pop(context),
    backLabel: 'К списку пациентов',
    actions: [
      if (_patient != null)
        FilledButton.icon(
          onPressed: _busy || !_readOnly ? null : _openRecord,
          icon: Icon(
            _opening
                ? Icons.hourglass_top_rounded
                : Icons.folder_shared_outlined,
            size: 20,
          ),
          label: Text(
            _opening
                ? 'Открытие…'
                : _patient!.hasRecord
                ? 'Открыть ЭМК'
                : 'Создать ЭМК',
          ),
        ),
    ],
    child: ListView(
      children: [
        Align(
          alignment: Alignment.topLeft,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 950),
            child: Panel(
              child: Form(
                key: _form,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(
                          Icons.person_outline_rounded,
                          color: MedicalColors.teal,
                        ),
                        const SizedBox(width: 12),
                        const Expanded(
                          child: Text(
                            'Персональные данные',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ),
                        if (_patient != null)
                          StatusPill(
                            _patient!.hasRecord ? 'ЭМК есть' : 'Нет ЭМК',
                            tone: _patient!.hasRecord
                                ? MedicalColors.teal
                                : MedicalColors.muted,
                          ),
                      ],
                    ),
                    const SizedBox(height: 26),
                    FormError(_error),
                    FormFields(
                      children: [
                        TextFormField(
                          controller: _last,
                          enabled: !_readOnly && !_busy,
                          decoration: const InputDecoration(
                            labelText: 'Фамилия',
                          ),
                          validator: (v) => validateText(v, 1, 64),
                        ),
                        TextFormField(
                          controller: _first,
                          enabled: !_readOnly && !_busy,
                          decoration: const InputDecoration(labelText: 'Имя'),
                          validator: (v) => validateText(v, 1, 64),
                        ),
                        TextFormField(
                          controller: _middle,
                          enabled: !_readOnly && !_busy,
                          decoration: const InputDecoration(
                            labelText: 'Отчество',
                            hintText: 'При наличии',
                          ),
                          validator: (v) =>
                              validateText(v, 1, 64, optional: true),
                        ),
                        TextFormField(
                          controller: _birth,
                          readOnly: true,
                          enabled: !_readOnly && !_busy,
                          onTap: _pickDate,
                          decoration: const InputDecoration(
                            labelText: 'Дата рождения',
                            suffixIcon: Icon(
                              Icons.calendar_today_outlined,
                              size: 20,
                            ),
                          ),
                          validator: (_) => _birthDate == null
                              ? 'Выберите дату рождения'
                              : null,
                        ),
                        DropdownButtonFormField<String>(
                          key: ValueKey(_sex),
                          initialValue: _sex,
                          decoration: const InputDecoration(
                            labelText: 'Административный пол',
                          ),
                          items: const [
                            DropdownMenuItem(
                              value: 'M',
                              child: Text('М — мужской'),
                            ),
                            DropdownMenuItem(
                              value: 'F',
                              child: Text('Ж — женский'),
                            ),
                            DropdownMenuItem(
                              value: 'UNKNOWN',
                              child: Text('НУ — не установлен'),
                            ),
                          ],
                          onChanged: _readOnly || _busy
                              ? null
                              : (value) => setState(() => _sex = value!),
                        ),
                        TextFormField(
                          controller: _comment,
                          enabled: !_readOnly && !_busy,
                          decoration: const InputDecoration(
                            labelText: 'Комментарий',
                            helperText: 'До 32 символов',
                          ),
                          validator: (v) =>
                              validateText(v, 0, 32, optional: true),
                        ),
                      ],
                    ),
                    const SizedBox(height: 30),
                    Wrap(
                      spacing: 12,
                      runSpacing: 12,
                      children: [
                        if (_readOnly)
                          FilledButton.icon(
                            onPressed: _busy
                                ? null
                                : () => setState(() => _readOnly = false),
                            icon: const Icon(Icons.edit_outlined, size: 18),
                            label: const Text('Редактировать'),
                          )
                        else
                          FilledButton.icon(
                            onPressed: _busy ? null : _save,
                            icon: const Icon(Icons.check_rounded, size: 18),
                            label: Text(_saving ? 'Сохранение…' : 'Сохранить'),
                          ),
                        if (!_readOnly)
                          OutlinedButton(
                            onPressed: _busy ? null : _cancelEdit,
                            child: const Text('Отмена'),
                          ),
                        if (_patient != null)
                          TextButton.icon(
                            onPressed: _patient!.hasRecord || _busy
                                ? null
                                : _delete,
                            icon: const Icon(
                              Icons.delete_outline_rounded,
                              size: 18,
                            ),
                            label: const Text('Удалить пациента'),
                            style: TextButton.styleFrom(
                              foregroundColor: MedicalColors.red,
                            ),
                          ),
                      ],
                    ),
                    if (_patient?.hasRecord == true)
                      const Padding(
                        padding: EdgeInsets.only(top: 18),
                        child: Text(
                          'Удаление пациента с медицинской картой запрещено.',
                          style: TextStyle(
                            fontSize: 12,
                            color: MedicalColors.muted,
                          ),
                        ),
                      ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    ),
  );
}
