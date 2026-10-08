import 'package:file_selector/file_selector.dart';
import 'package:flutter/material.dart';
import '../core/models.dart';
import '../core/session.dart';
import '../core/theme.dart';
import '../widgets/common.dart';
import '../widgets/paged_table.dart';

class DiagnosisScreen extends StatefulWidget {
  const DiagnosisScreen({super.key, required this.session});
  final AppSession session;
  @override
  State<DiagnosisScreen> createState() => _DiagnosisScreenState();
}

class _DiagnosisScreenState extends State<DiagnosisScreen> {
  final _table = GlobalKey<PagedTableState<Diagnosis>>();
  bool _importing = false;
  Future<void> _edit([Diagnosis? diagnosis]) async {
    final changed = await showDialog<bool>(
      context: context,
      builder: (_) =>
          _DiagnosisForm(session: widget.session, diagnosis: diagnosis),
    );
    if (changed == true) _table.currentState?.reload();
  }

  Future<void> _delete(Diagnosis value) async {
    if (!await confirm(
      context,
      title: 'Удалить диагноз из справочника?',
      message:
          '${value.code} — ${value.name}. В ранее созданных медицинских картах он сохранится.',
    )) {
      return;
    }
    try {
      await widget.session.api.request(
        'DELETE',
        '/api/v1/diagnoses/${value.id}',
      );
      if (mounted) {
        notify(context, 'Диагноз удалён из справочника');
        _table.currentState?.reload();
      }
    } catch (error) {
      if (mounted) notify(context, errorMessage(error), error: true);
    }
  }

  Future<void> _import() async {
    setState(() => _importing = true);
    try {
      final file = await openFile(
        acceptedTypeGroups: [
          const XTypeGroup(
            label: 'Справочник диагнозов (JSON)',
            extensions: ['json'],
            uniformTypeIdentifiers: ['public.json'],
          ),
        ],
      );
      if (file == null) return;
      if (await file.length() > 1024 * 1024) {
        throw const FormatException('Размер файла не должен превышать 1 МБ.');
      }
      final data = await widget.session.api.request(
        'POST',
        '/api/v1/diagnoses/import',
        fileContents: await file.readAsString(),
      );
      if (!mounted) return;
      notify(
        context,
        'Импорт завершён. Загружено записей: ${(data['items'] as List).length}',
      );
      _table.currentState?.reload(reset: true);
    } catch (error) {
      if (mounted) {
        await showDialog<void>(
          context: context,
          builder: (context) => AlertDialog(
            title: const Text('Файл не импортирован'),
            content: SizedBox(
              width: 460,
              child: SelectableText(errorMessage(error)),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: const Text('Понятно'),
              ),
            ],
          ),
        );
      }
    } finally {
      if (mounted) setState(() => _importing = false);
    }
  }

  @override
  Widget build(BuildContext context) => PageFrame(
    title: 'Справочник диагнозов',
    subtitle: 'Единый реестр для электронных медицинских карт',
    actions: [
      OutlinedButton.icon(
        onPressed: _importing ? null : _import,
        icon: Icon(
          _importing ? Icons.hourglass_top_rounded : Icons.upload_file_rounded,
          size: 20,
        ),
        label: Text(_importing ? 'Импорт…' : 'Загрузить из файла'),
      ),
      FilledButton.icon(
        onPressed: () => _edit(),
        icon: const Icon(Icons.add_rounded, size: 20),
        label: const Text('Добавить диагноз'),
      ),
    ],
    child: PagedTable<Diagnosis>(
      key: _table,
      fetch: (offset, query) => widget.session.api.page(
        '/api/v1/diagnoses',
        Diagnosis.fromJson,
        offset: offset,
        query: query,
      ),
      searchHint: 'Поиск по коду или наименованию',
      emptyTitle: 'Справочник пока пуст',
      emptySubtitle:
          'Добавьте диагноз вручную или загрузите JSON-файл.\nКод диагноза — 6 символов, наименование — от 3 до 24.',
      columns: const [
        DataColumn(label: Text('КОД')),
        DataColumn(label: Text('НАИМЕНОВАНИЕ')),
        DataColumn(label: Text('ДЕЙСТВИЯ')),
      ],
      row: (d) => DataRow(
        onSelectChanged: (_) => _edit(d),
        cells: [
          DataCell(
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
              decoration: BoxDecoration(
                color: MedicalColors.background,
                borderRadius: BorderRadius.circular(6),
              ),
              child: Text(
                d.code,
                style: const TextStyle(
                  fontWeight: FontWeight.w600,
                  letterSpacing: .5,
                ),
              ),
            ),
          ),
          DataCell(Text(d.name)),
          DataCell(
            Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                IconButton(
                  tooltip: 'Редактировать диагноз',
                  onPressed: () => _edit(d),
                  icon: const Icon(Icons.edit_outlined, size: 19),
                ),
                IconButton(
                  tooltip: 'Удалить диагноз',
                  onPressed: () => _delete(d),
                  icon: const Icon(
                    Icons.delete_outline_rounded,
                    size: 19,
                    color: MedicalColors.red,
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

class _DiagnosisForm extends StatefulWidget {
  const _DiagnosisForm({required this.session, this.diagnosis});
  final AppSession session;
  final Diagnosis? diagnosis;
  @override
  State<_DiagnosisForm> createState() => _DiagnosisFormState();
}

class _DiagnosisFormState extends State<_DiagnosisForm> {
  final _form = GlobalKey<FormState>();
  late final TextEditingController _code, _name;
  bool _saving = false;
  Object? _error;
  @override
  void initState() {
    super.initState();
    _code = TextEditingController(text: widget.diagnosis?.code);
    _name = TextEditingController(text: widget.diagnosis?.name);
  }

  @override
  void dispose() {
    _code.dispose();
    _name.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (!_form.currentState!.validate()) return;
    setState(() {
      _saving = true;
      _error = null;
    });
    try {
      await widget.session.api.request(
        widget.diagnosis == null ? 'POST' : 'PUT',
        '/api/v1/diagnoses${widget.diagnosis == null ? '' : '/${widget.diagnosis!.id}'}',
        body: {'code': _code.text.trim(), 'name': _name.text.trim()},
      );
      if (mounted) Navigator.pop(context, true);
    } catch (error) {
      if (mounted) setState(() => _error = error);
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  @override
  Widget build(BuildContext context) => AlertDialog(
    title: Text(
      widget.diagnosis == null ? 'Добавить диагноз' : 'Редактировать диагноз',
    ),
    content: SizedBox(
      width: 460,
      child: Form(
        key: _form,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            FormError(_error),
            TextFormField(
              controller: _code,
              enabled: !_saving,
              autofocus: true,
              decoration: const InputDecoration(
                labelText: 'Код диагноза',
                helperText: 'Ровно 6 символов, например E11.00',
              ),
              validator: (v) => validateText(v, 6, 6),
            ),
            const SizedBox(height: 20),
            TextFormField(
              controller: _name,
              enabled: !_saving,
              decoration: const InputDecoration(labelText: 'Наименование'),
              validator: (v) => validateText(v, 3, 24),
            ),
          ],
        ),
      ),
    ),
    actions: [
      TextButton(
        onPressed: _saving ? null : () => Navigator.pop(context),
        child: const Text('Отмена'),
      ),
      FilledButton(
        onPressed: _saving ? null : _save,
        child: Text(_saving ? 'Сохранение…' : 'Сохранить'),
      ),
    ],
  );
}
