import 'package:flutter/material.dart';
import '../core/api_client.dart';
import '../core/models.dart';
import '../core/record_draft.dart';
import '../core/session.dart';
import '../core/theme.dart';
import '../widgets/common.dart';
import '../widgets/staff_picker.dart';

class RecordScreen extends StatefulWidget {
  const RecordScreen({
    super.key,
    required this.session,
    required this.patient,
    required this.record,
  });
  final AppSession session;
  final Patient patient;
  final MedicalRecord record;
  @override
  State<RecordScreen> createState() => _RecordScreenState();
}

class _RecordScreenState extends State<RecordScreen> {
  late MedicalRecord _record;
  late RecordDraft _draft;
  late final Future<bool> Function() _navigationGuard;
  List<RecordEvent> _events = [];
  MedicalRecord? _historical;
  Staff? _actor;
  Object? _error, _historyError;
  bool _editing = false,
      _saving = false,
      _historyLoading = true,
      _allowPop = false,
      _conflict = false;
  bool get _canEdit => _editing && _historical == null && !_saving;
  @override
  void initState() {
    super.initState();
    _record = widget.record;
    _draft = RecordDraft(_record);
    _navigationGuard = _canLeave;
    widget.session.beforeNavigation = _navigationGuard;
    _loadHistory();
  }

  @override
  void dispose() {
    if (identical(widget.session.beforeNavigation, _navigationGuard)) {
      widget.session.beforeNavigation = null;
    }
    _draft.dispose();
    super.dispose();
  }

  Future<bool> _canLeave() async {
    if (_saving) return false;
    if (!_draft.dirty) return true;
    return confirm(
      context,
      title: 'Отменить несохранённые изменения?',
      message: 'Изменения медицинской карты ещё не сохранены.',
      action: 'Отменить изменения',
    );
  }

  Future<void> _leave() async {
    if (!await _canLeave() || !mounted) return;
    setState(() => _allowPop = true);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) Navigator.pop(context);
    });
  }

  Future<void> _loadHistory() async {
    setState(() {
      _historyLoading = true;
      _historyError = null;
    });
    try {
      final events = await widget.session.api.history(_record.id);
      if (mounted) setState(() => _events = events);
    } catch (error) {
      if (mounted) setState(() => _historyError = error);
    } finally {
      if (mounted) setState(() => _historyLoading = false);
    }
  }

  Future<void> _unlock() async {
    final actor = await selectStaff(
      context,
      widget.session,
      title: 'Разблокировать ЭМК',
      action: 'Разблокировать',
    );
    if (actor != null && mounted) {
      setState(() {
        _actor = actor;
        _editing = true;
        _error = null;
      });
    }
  }

  void _replaceRecord(MedicalRecord value) {
    _draft.dispose();
    _record = value;
    _draft = RecordDraft(value);
    _historical = null;
    _editing = false;
    _actor = null;
    _conflict = false;
    _error = null;
  }

  Future<void> _save() async {
    final commands = _draft.commands;
    if (commands.isEmpty || _actor == null || _saving) return;
    if (commands.length > 100) {
      notify(
        context,
        'Слишком много изменений. Сохраните не более 100 действий за один раз.',
        error: true,
      );
      return;
    }
    setState(() {
      _saving = true;
      _error = null;
    });
    try {
      final record = await widget.session.api.saveRecord(
        _record,
        _actor!.id,
        commands,
      );
      if (!mounted) return;
      setState(() => _replaceRecord(record));
      notify(
        context,
        'Изменения сохранены. Карта снова доступна только для просмотра.',
      );
      _loadHistory();
    } catch (error) {
      if (mounted) {
        setState(() {
          _error = error;
          _conflict = error is ApiException && error.code == 'VERSION_CONFLICT';
        });
      }
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _cancel() async {
    if (!await _canLeave() || !mounted) return;
    setState(() => _replaceRecord(_record));
  }

  Future<void> _reload() async {
    if (_draft.dirty &&
        !await confirm(
          context,
          title: 'Обновить карту?',
          message:
              'Карта будет загружена заново. Текущий черновик не будет сохранён.',
          action: 'Обновить карту',
        )) {
      return;
    }
    setState(() => _saving = true);
    try {
      final record = await widget.session.api.patientRecord(widget.patient.id);
      if (!mounted) return;
      setState(() => _replaceRecord(record));
      _loadHistory();
    } catch (error) {
      if (mounted) setState(() => _error = error);
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _showState({int? version, DateTime? at}) async {
    try {
      final record = await widget.session.api.historicalRecord(
        _record.id,
        version: version,
        at: at,
      );
      if (mounted) {
        setState(() {
          _historical = record;
          _error = null;
        });
      }
    } catch (error) {
      if (mounted) notify(context, errorMessage(error), error: true);
    }
  }

  Future<void> _chooseMoment() async {
    final at = await showDialog<DateTime>(
      context: context,
      builder: (_) => const _MomentDialog(),
    );
    if (at != null) await _showState(at: at);
  }

  Future<void> _addDiagnosis() async {
    final diagnosis = await showDialog<Diagnosis>(
      context: context,
      builder: (_) => _DiagnosisPicker(
        session: widget.session,
        excluded: _draft.diagnoses.map((d) => d.diagnosisId).toSet(),
      ),
    );
    if (diagnosis != null && mounted) _draft.addDiagnosis(diagnosis);
  }

  Future<void> _restoreDiagnosis(RecordDiagnosis value) async {
    if (_draft.base.diagnoses.any((d) => d.id == value.id)) {
      _draft.undoDiagnosisRemoval(value);
      return;
    }
    try {
      final diagnoses = await widget.session.api.diagnoses();
      final active = diagnoses
          .where((d) => d.id == value.diagnosisId)
          .firstOrNull;
      if (!mounted) return;
      if (active == null) {
        notify(
          context,
          'Диагноз больше не доступен в справочнике.',
          error: true,
        );
        return;
      }
      _draft.addDiagnosis(active);
    } catch (error) {
      if (mounted) notify(context, errorMessage(error), error: true);
    }
  }

  Future<void> _prescription([Prescription? value]) async {
    final text = await showDialog<String>(
      context: context,
      builder: (_) => _PrescriptionDialog(text: value?.text),
    );
    if (text != null && mounted) {
      if (value == null) {
        _draft.addPrescription(text);
      } else {
        _draft.editPrescription(value.id, text);
      }
    }
  }

  Future<void> _cancelPrescription(Prescription value) async {
    final isNew = _draft.isNewPrescription(value.id);
    if (!await confirm(
      context,
      title: isNew ? 'Убрать назначение из черновика?' : 'Отменить назначение?',
      message: value.text,
      action: isNew ? 'Убрать' : 'Отменить назначение',
    )) {
      return;
    }
    if (mounted) _draft.cancelPrescription(value.id);
  }

  DateTime? _addedAt(String id) => _events
      .where(
        (e) =>
            e.entityId == id &&
            (e.type == 'diagnosis_added' || e.type == 'prescription_added'),
      )
      .firstOrNull
      ?.at;
  List<RecordDiagnosis> _removedDiagnoses(List<RecordDiagnosis> active) {
    final result = <String, RecordDiagnosis>{};
    for (final event in _events) {
      if (event.type == 'diagnosis_removed' && event.before != null) {
        final value = RecordDiagnosis.fromJson(event.before!);
        result[value.diagnosisId] = value;
      }
    }
    for (final diagnosis in _record.diagnoses) {
      if (!active.any((d) => d.id == diagnosis.id)) {
        result[diagnosis.diagnosisId] = diagnosis;
      }
    }
    result.removeWhere((id, _) => active.any((d) => d.diagnosisId == id));
    return result.values.toList();
  }

  Widget _diagnosesTab() {
    final active = _historical?.diagnoses ?? _draft.diagnoses;
    final removed = _historical == null
        ? _removedDiagnoses(active)
        : <RecordDiagnosis>[];
    final all = [...active, ...removed];
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(20),
          child: Row(
            children: [
              const Expanded(
                child: Text(
                  'Диагнозы пациента',
                  style: TextStyle(fontSize: 17, fontWeight: FontWeight.w600),
                ),
              ),
              FilledButton.icon(
                onPressed: _canEdit ? _addDiagnosis : null,
                icon: const Icon(Icons.add_rounded, size: 18),
                label: const Text('Добавить диагноз'),
              ),
            ],
          ),
        ),
        const Divider(height: 1),
        Expanded(
          child: all.isEmpty
              ? const EmptyPanel(
                  title: 'Диагнозы не добавлены',
                  subtitle: 'Разблокируйте карту, чтобы добавить диагноз.',
                  icon: Icons.medical_information_outlined,
                )
              : _table(
                  columns: const [
                    DataColumn(label: Text('ДАТА')),
                    DataColumn(label: Text('ДИАГНОЗ')),
                    DataColumn(label: Text('СТАТУС')),
                    DataColumn(label: Text('ДЕЙСТВИЯ')),
                  ],
                  rows: all.map((value) {
                    final isActive = active.contains(value);
                    final date = _addedAt(value.id);
                    return DataRow(
                      cells: [
                        DataCell(
                          Text(
                            date == null
                                ? 'В черновике'
                                : displayDate(date.toLocal()),
                            style: const TextStyle(fontSize: 12),
                          ),
                        ),
                        DataCell(
                          Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                value.name,
                                style: const TextStyle(
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                value.code,
                                style: const TextStyle(
                                  fontSize: 12,
                                  color: MedicalColors.muted,
                                ),
                              ),
                            ],
                          ),
                        ),
                        DataCell(
                          StatusPill(
                            isActive ? 'Активен' : 'Удалён',
                            tone: isActive
                                ? MedicalColors.green
                                : MedicalColors.red,
                          ),
                        ),
                        DataCell(
                          IconButton(
                            tooltip: isActive
                                ? 'Удалить диагноз из ЭМК'
                                : 'Восстановить диагноз',
                            onPressed: !_canEdit
                                ? null
                                : isActive
                                ? () => _draft.removeDiagnosis(value.id)
                                : () => _restoreDiagnosis(value),
                            icon: Icon(
                              isActive
                                  ? Icons.delete_outline_rounded
                                  : Icons.restore_rounded,
                              size: 20,
                              color: isActive
                                  ? MedicalColors.red
                                  : MedicalColors.teal,
                            ),
                          ),
                        ),
                      ],
                    );
                  }).toList(),
                ),
        ),
      ],
    );
  }

  Widget _prescriptionsTab() {
    final prescriptions = _historical?.prescriptions ?? _draft.prescriptions;
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(20),
          child: Row(
            children: [
              const Expanded(
                child: Text(
                  'Назначения',
                  style: TextStyle(fontSize: 17, fontWeight: FontWeight.w600),
                ),
              ),
              FilledButton.icon(
                onPressed: _canEdit ? () => _prescription() : null,
                icon: const Icon(Icons.add_rounded, size: 18),
                label: const Text('Добавить назначение'),
              ),
            ],
          ),
        ),
        const Divider(height: 1),
        Expanded(
          child: prescriptions.isEmpty
              ? const EmptyPanel(
                  title: 'Назначений пока нет',
                  subtitle: 'Добавьте назначение в режиме редактирования.',
                  icon: Icons.assignment_outlined,
                )
              : ListView.separated(
                  padding: const EdgeInsets.all(20),
                  itemCount: prescriptions.length,
                  separatorBuilder: (_, index) => const SizedBox(height: 12),
                  itemBuilder: (context, index) {
                    final value = prescriptions[index];
                    final fresh =
                        _historical == null &&
                        _draft.isNewPrescription(value.id);
                    final edited =
                        _historical == null &&
                        _draft.isEditedPrescription(value.id);
                    return Container(
                      padding: const EdgeInsets.all(18),
                      decoration: BoxDecoration(
                        border: Border.all(color: MedicalColors.line),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Row(
                        children: [
                          Container(
                            padding: const EdgeInsets.all(10),
                            decoration: BoxDecoration(
                              color: MedicalColors.background,
                              borderRadius: BorderRadius.circular(10),
                            ),
                            child: const Icon(
                              Icons.description_outlined,
                              color: MedicalColors.teal,
                              size: 22,
                            ),
                          ),
                          const SizedBox(width: 14),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  value.text,
                                  style: const TextStyle(
                                    fontSize: 14,
                                    fontWeight: FontWeight.w500,
                                  ),
                                ),
                                const SizedBox(height: 8),
                                StatusPill(
                                  fresh
                                      ? 'Новое'
                                      : edited
                                      ? 'Изменено'
                                      : value.active
                                      ? 'Активно'
                                      : 'Отменено',
                                  tone: value.active
                                      ? MedicalColors.green
                                      : MedicalColors.muted,
                                ),
                              ],
                            ),
                          ),
                          IconButton(
                            tooltip: 'Редактировать назначение',
                            onPressed: _canEdit && value.active
                                ? () => _prescription(value)
                                : null,
                            icon: const Icon(Icons.edit_outlined, size: 19),
                          ),
                          IconButton(
                            tooltip: fresh
                                ? 'Убрать из черновика'
                                : edited
                                ? 'Сначала сохраните изменение назначения'
                                : 'Отменить назначение',
                            onPressed: _canEdit && value.active && !edited
                                ? () => _cancelPrescription(value)
                                : null,
                            icon: Icon(
                              fresh
                                  ? Icons.delete_outline_rounded
                                  : Icons.block_outlined,
                              color: MedicalColors.red,
                              size: 19,
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
        ),
      ],
    );
  }

  Widget _historyTab() => Column(
    children: [
      Padding(
        padding: const EdgeInsets.all(20),
        child: Wrap(
          spacing: 12,
          runSpacing: 12,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            const Text(
              'История изменений',
              style: TextStyle(fontSize: 17, fontWeight: FontWeight.w600),
            ),
            OutlinedButton.icon(
              onPressed: _chooseMoment,
              icon: const Icon(Icons.history_rounded, size: 18),
              label: const Text('Состояние на дату и время'),
            ),
            IconButton(
              tooltip: 'Обновить историю',
              onPressed: _loadHistory,
              icon: const Icon(Icons.refresh_rounded),
            ),
          ],
        ),
      ),
      if (_historyLoading) const LinearProgressIndicator(minHeight: 2),
      Expanded(
        child: _historyError != null
            ? ErrorPanel(_historyError!, retry: _loadHistory)
            : ListView.separated(
                padding: const EdgeInsets.fromLTRB(20, 0, 20, 20),
                itemCount: _events.length,
                separatorBuilder: (_, index) => const SizedBox(height: 10),
                itemBuilder: (context, index) {
                  final event = _events[_events.length - index - 1];
                  return DecoratedBox(
                    decoration: BoxDecoration(
                      border: Border.all(color: MedicalColors.line),
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: ExpansionTile(
                      shape: const Border(),
                      collapsedShape: const Border(),
                      leading: const Icon(
                        Icons.receipt_long_outlined,
                        color: MedicalColors.teal,
                      ),
                      title: Text(
                        event.label,
                        style: const TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      subtitle: Padding(
                        padding: const EdgeInsets.only(top: 5),
                        child: Text(
                          '${displayTime(event.at)}\n${event.actor.name} · ${event.actor.position}${event.description.isEmpty ? '' : '\n${event.description}'}',
                          style: const TextStyle(fontSize: 12, height: 1.5),
                        ),
                      ),
                      childrenPadding: const EdgeInsets.all(18),
                      expandedCrossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        if (event.type != 'record_created') ...[
                          const Text(
                            'Данные до изменения',
                            style: TextStyle(fontWeight: FontWeight.w600),
                          ),
                          const SizedBox(height: 6),
                          SelectableText(_snapshot(event.before)),
                          const SizedBox(height: 14),
                          const Text(
                            'Данные после изменения',
                            style: TextStyle(fontWeight: FontWeight.w600),
                          ),
                          const SizedBox(height: 6),
                          SelectableText(_snapshot(event.after)),
                          const SizedBox(height: 14),
                        ],
                        OutlinedButton.icon(
                          onPressed: () => _showState(version: event.version),
                          icon: const Icon(Icons.visibility_outlined, size: 18),
                          label: Text(
                            'Посмотреть карту · версия ${event.version}',
                          ),
                        ),
                      ],
                    ),
                  );
                },
              ),
      ),
    ],
  );
  String _snapshot(Json? value) {
    if (value == null) return 'Запись отсутствует';
    return value.entries
        .where((e) => !['id', 'diagnosis_id', 'patient_id'].contains(e.key))
        .map(
          (e) =>
              '${switch (e.key) {
                'text' => 'Назначение',
                'name' => 'Наименование',
                'code' => 'Код',
                'status' => 'Статус',
                _ => e.key,
              }}: ${e.value == 'active'
                  ? 'Активно'
                  : e.value == 'cancelled'
                  ? 'Отменено'
                  : e.value}',
        )
        .join('\n');
  }

  Widget _table({
    required List<DataColumn> columns,
    required List<DataRow> rows,
  }) => LayoutBuilder(
    builder: (context, constraints) => SingleChildScrollView(
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: ConstrainedBox(
          constraints: BoxConstraints(minWidth: constraints.maxWidth),
          child: DataTable(columns: columns, rows: rows),
        ),
      ),
    ),
  );
  @override
  Widget build(BuildContext context) => ListenableBuilder(
    listenable: _draft,
    builder: (context, _) => PopScope(
      canPop: !_draft.dirty || _allowPop,
      onPopInvokedWithResult: (didPop, result) {
        if (!didPop) _leave();
      },
      child: PageFrame(
        title: 'Электронная медицинская карта',
        subtitle:
            '${widget.patient.name}  ·  ${displayDate(widget.patient.birthDate)}  ·  ${widget.patient.sexLabel}',
        back: _saving ? null : _leave,
        backLabel: 'К карточке пациента',
        actions: [
          IconButton(
            tooltip: 'Обновить карту',
            onPressed: _saving ? null : _reload,
            icon: const Icon(Icons.refresh_rounded),
          ),
        ],
        child: DefaultTabController(
          length: 3,
          child: Column(
            children: [
              Panel(
                padding: const EdgeInsets.symmetric(
                  horizontal: 20,
                  vertical: 16,
                ),
                child: Row(
                  children: [
                    Icon(
                      _historical != null
                          ? Icons.history_rounded
                          : _editing
                          ? Icons.lock_open_rounded
                          : Icons.lock_outline_rounded,
                      color: MedicalColors.teal,
                      size: 25,
                    ),
                    const SizedBox(width: 14),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            _historical != null
                                ? 'Историческое состояние · версия ${_historical!.version}'
                                : _editing
                                ? 'Редактирование карты'
                                : 'Только просмотр',
                            style: const TextStyle(fontWeight: FontWeight.w600),
                          ),
                          const SizedBox(height: 4),
                          Text(
                            _historical != null
                                ? displayTime(_historical!.updatedAt)
                                : _editing
                                ? 'Автор изменений: ${_actor!.name}'
                                : 'Для внесения изменений выберите медработника.',
                            style: const TextStyle(
                              fontSize: 12,
                              color: MedicalColors.muted,
                            ),
                          ),
                        ],
                      ),
                    ),
                    if (_historical != null)
                      OutlinedButton(
                        onPressed: () => setState(() => _historical = null),
                        child: const Text('К текущей карте'),
                      )
                    else if (!_editing)
                      OutlinedButton.icon(
                        onPressed: _saving ? null : _unlock,
                        icon: const Icon(Icons.lock_open_rounded, size: 18),
                        label: const Text('Разблокировать'),
                      ),
                  ],
                ),
              ),
              const SizedBox(height: 16),
              if (_error != null) FormError(_error),
              if (_conflict)
                Padding(
                  padding: const EdgeInsets.only(bottom: 12),
                  child: Row(
                    children: [
                      const Expanded(
                        child: Text(
                          'Черновик сохранён. Обновите карту, чтобы увидеть актуальные данные.',
                        ),
                      ),
                      OutlinedButton(
                        onPressed: _reload,
                        child: const Text('Обновить карту'),
                      ),
                    ],
                  ),
                ),
              Expanded(
                child: Panel(
                  padding: EdgeInsets.zero,
                  child: Column(
                    children: [
                      const TabBar(
                        labelStyle: TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                        tabs: [
                          Tab(text: 'Диагнозы'),
                          Tab(text: 'Назначения'),
                          Tab(text: 'История изменений'),
                        ],
                      ),
                      Expanded(
                        child: TabBarView(
                          children: [
                            _diagnosesTab(),
                            _prescriptionsTab(),
                            _historyTab(),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              if (_editing && _historical == null)
                Padding(
                  padding: const EdgeInsets.only(top: 16),
                  child: Row(
                    children: [
                      Expanded(
                        child: Text(
                          _draft.dirty
                              ? 'Несохранённых действий: ${_draft.commands.length}'
                              : 'Изменений пока нет',
                          style: const TextStyle(color: MedicalColors.muted),
                        ),
                      ),
                      OutlinedButton(
                        onPressed: _saving ? null : _cancel,
                        child: const Text('Отмена'),
                      ),
                      const SizedBox(width: 12),
                      FilledButton.icon(
                        onPressed: _saving || !_draft.dirty ? null : _save,
                        icon: const Icon(Icons.check_rounded, size: 18),
                        label: Text(
                          _saving ? 'Сохранение…' : 'Сохранить изменения',
                        ),
                      ),
                    ],
                  ),
                ),
            ],
          ),
        ),
      ),
    ),
  );
}

class _PrescriptionDialog extends StatefulWidget {
  const _PrescriptionDialog({this.text});
  final String? text;
  @override
  State<_PrescriptionDialog> createState() => _PrescriptionDialogState();
}

class _PrescriptionDialogState extends State<_PrescriptionDialog> {
  final _form = GlobalKey<FormState>();
  late final TextEditingController _text;
  @override
  void initState() {
    super.initState();
    _text = TextEditingController(text: widget.text);
  }

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => AlertDialog(
    title: Text(
      widget.text == null ? 'Добавить назначение' : 'Изменить назначение',
    ),
    content: SizedBox(
      width: 480,
      child: Form(
        key: _form,
        child: TextFormField(
          controller: _text,
          autofocus: true,
          minLines: 3,
          maxLines: 5,
          decoration: const InputDecoration(
            labelText: 'Назначение',
            helperText: 'От 3 до 128 символов',
          ),
          validator: (v) => validateText(v, 3, 128),
        ),
      ),
    ),
    actions: [
      TextButton(
        onPressed: () => Navigator.pop(context),
        child: const Text('Отмена'),
      ),
      FilledButton(
        onPressed: () {
          if (_form.currentState!.validate()) {
            Navigator.pop(context, _text.text.trim());
          }
        },
        child: const Text('Добавить в карту'),
      ),
    ],
  );
}

class _DiagnosisPicker extends StatefulWidget {
  const _DiagnosisPicker({required this.session, required this.excluded});
  final AppSession session;
  final Set<String> excluded;
  @override
  State<_DiagnosisPicker> createState() => _DiagnosisPickerState();
}

class _DiagnosisPickerState extends State<_DiagnosisPicker> {
  final _search = TextEditingController();
  List<Diagnosis>? _values;
  Object? _error;
  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  void dispose() {
    _search.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final values = await widget.session.api.diagnoses();
      if (mounted) {
        setState(() {
          _values = values
              .where((d) => !widget.excluded.contains(d.id))
              .toList();
          _error = null;
        });
      }
    } catch (error) {
      if (mounted) setState(() => _error = error);
    }
  }

  @override
  Widget build(BuildContext context) {
    final query = _search.text.toLowerCase();
    final values = _values
        ?.where((d) => '${d.code} ${d.name}'.toLowerCase().contains(query))
        .toList();
    return AlertDialog(
      title: const Text('Добавить диагноз в ЭМК'),
      content: SizedBox(
        width: 500,
        height: 380,
        child: Column(
          children: [
            TextField(
              controller: _search,
              decoration: const InputDecoration(
                hintText: 'Код или наименование',
                prefixIcon: Icon(Icons.search_rounded),
              ),
              onChanged: (_) => setState(() {}),
            ),
            const SizedBox(height: 14),
            Expanded(
              child: _error != null
                  ? ErrorPanel(_error!, retry: _load)
                  : values == null
                  ? const Center(child: CircularProgressIndicator())
                  : values.isEmpty
                  ? const EmptyPanel(
                      title: 'Нет доступных диагнозов',
                      subtitle: 'Загрузите справочник или измените поиск.',
                    )
                  : ListView.builder(
                      itemCount: values.length,
                      itemBuilder: (context, index) {
                        final d = values[index];
                        return ListTile(
                          title: Text(d.name),
                          subtitle: Text(d.code),
                          trailing: const Icon(
                            Icons.add_circle_outline_rounded,
                            color: MedicalColors.teal,
                          ),
                          onTap: () => Navigator.pop(context, d),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Отмена'),
        ),
      ],
    );
  }
}

class _MomentDialog extends StatefulWidget {
  const _MomentDialog();
  @override
  State<_MomentDialog> createState() => _MomentDialogState();
}

class _MomentDialogState extends State<_MomentDialog> {
  final _form = GlobalKey<FormState>();
  late final TextEditingController _text;
  @override
  void initState() {
    super.initState();
    _text = TextEditingController(
      text: displayTime(DateTime.now()).replaceAll(' · ', ' '),
    );
  }

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  DateTime? _parse() {
    final match = RegExp(
      r'^(\d{2})\.(\d{2})\.(\d{4}) (\d{2}):(\d{2}):(\d{2})$',
    ).firstMatch(_text.text.trim());
    if (match == null) return null;
    final numbers = List.generate(6, (i) => int.parse(match.group(i + 1)!));
    final date = DateTime(
      numbers[2],
      numbers[1],
      numbers[0],
      numbers[3],
      numbers[4],
      numbers[5],
    );
    return date.year == numbers[2] &&
            date.month == numbers[1] &&
            date.day == numbers[0] &&
            date.hour == numbers[3] &&
            date.minute == numbers[4] &&
            date.second == numbers[5]
        ? date
        : null;
  }

  @override
  Widget build(BuildContext context) => AlertDialog(
    title: const Text('Карта на выбранный момент'),
    content: SizedBox(
      width: 420,
      child: Form(
        key: _form,
        child: TextFormField(
          controller: _text,
          autofocus: true,
          decoration: const InputDecoration(
            labelText: 'Дата и время',
            helperText: 'ДД.ММ.ГГГГ ЧЧ:ММ:СС, местное время',
          ),
          validator: (_) =>
              _parse() == null ? 'Введите корректные дату и время' : null,
        ),
      ),
    ),
    actions: [
      TextButton(
        onPressed: () => Navigator.pop(context),
        child: const Text('Отмена'),
      ),
      FilledButton(
        onPressed: () {
          if (_form.currentState!.validate()) {
            Navigator.pop(context, _parse());
          }
        },
        child: const Text('Посмотреть'),
      ),
    ],
  );
}
