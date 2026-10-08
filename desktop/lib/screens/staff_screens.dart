import 'package:flutter/material.dart';
import '../core/models.dart';
import '../core/session.dart';
import '../core/theme.dart';
import '../widgets/common.dart';
import '../widgets/paged_table.dart';

class StaffListScreen extends StatefulWidget {
  const StaffListScreen({super.key, required this.session});
  final AppSession session;
  @override
  State<StaffListScreen> createState() => _StaffListScreenState();
}

class _StaffListScreenState extends State<StaffListScreen> {
  final _table = GlobalKey<PagedTableState<Staff>>();
  bool _includeDeleted = true;
  Future<void> _edit([Staff? staff]) async {
    await Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => StaffFormScreen(session: widget.session, staff: staff),
      ),
    );
    _table.currentState?.reload();
  }

  Future<void> _delete(Staff staff) async {
    if (!await confirm(
      context,
      title: 'Пометить сотрудника удалённым?',
      message:
          '${staff.name} останется в истории карт, но не сможет вносить новые изменения.',
      action: 'Пометить удалённым',
    )) {
      return;
    }
    try {
      await widget.session.api.request('DELETE', '/api/v1/staff/${staff.id}');
      if (mounted) {
        notify(context, 'Сотрудник помечен удалённым');
        _table.currentState?.reload();
      }
    } catch (error) {
      if (mounted) notify(context, errorMessage(error), error: true);
    }
  }

  @override
  Widget build(BuildContext context) => PageFrame(
    title: 'Медработники',
    subtitle: 'Сотрудники лечебного учреждения',
    actions: [
      FilledButton.icon(
        onPressed: () => _edit(),
        icon: const Icon(Icons.add_rounded, size: 20),
        label: const Text('Добавить медработника'),
      ),
    ],
    child: PagedTable<Staff>(
      key: _table,
      fetch: (offset, query) => widget.session.api.page(
        '/api/v1/staff',
        Staff.fromJson,
        offset: offset,
        query: query,
        includeDeleted: _includeDeleted,
      ),
      searchHint: 'Поиск по имени или должности',
      emptyTitle: 'Добавьте первого медработника',
      emptySubtitle: 'Сотрудники появятся здесь после создания карточки.',
      extraFilter: FilterChip(
        label: const Text('Показывать удалённых'),
        selected: _includeDeleted,
        onSelected: (value) {
          setState(() => _includeDeleted = value);
          _table.currentState?.reload(reset: true);
        },
      ),
      columns: const [
        DataColumn(label: Text('ФАМИЛИЯ И ИМЯ')),
        DataColumn(label: Text('ДОЛЖНОСТЬ')),
        DataColumn(label: Text('СТАТУС')),
        DataColumn(label: Text('ДЕЙСТВИЯ')),
      ],
      row: (staff) => DataRow(
        onSelectChanged: (_) => _edit(staff),
        cells: [
          DataCell(
            Row(
              children: [
                CircleAvatar(
                  radius: 17,
                  backgroundColor: MedicalColors.mint,
                  child: Text(
                    staff.firstName.characters.first,
                    style: const TextStyle(
                      fontSize: 13,
                      color: MedicalColors.teal,
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Text(
                  staff.name,
                  style: const TextStyle(fontWeight: FontWeight.w600),
                ),
              ],
            ),
          ),
          DataCell(Text(staff.position)),
          DataCell(
            StatusPill(
              staff.active ? 'Активен' : 'Удалён',
              tone: staff.active ? MedicalColors.green : MedicalColors.red,
            ),
          ),
          DataCell(
            Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                IconButton(
                  tooltip: staff.active
                      ? 'Редактировать сотрудника'
                      : 'Посмотреть карточку',
                  onPressed: () => _edit(staff),
                  icon: Icon(
                    staff.active
                        ? Icons.edit_outlined
                        : Icons.visibility_outlined,
                    size: 19,
                  ),
                ),
                IconButton(
                  tooltip: 'Пометить удалённым',
                  onPressed: staff.active ? () => _delete(staff) : null,
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

class StaffFormScreen extends StatefulWidget {
  const StaffFormScreen({super.key, required this.session, this.staff});
  final AppSession session;
  final Staff? staff;
  @override
  State<StaffFormScreen> createState() => _StaffFormScreenState();
}

class _StaffFormScreenState extends State<StaffFormScreen> {
  final _form = GlobalKey<FormState>();
  late final TextEditingController _first, _last, _position;
  bool _saving = false, _delete = false;
  Object? _error;
  bool get _archived => widget.staff?.active == false;
  @override
  void initState() {
    super.initState();
    _first = TextEditingController(text: widget.staff?.firstName);
    _last = TextEditingController(text: widget.staff?.lastName);
    _position = TextEditingController(text: widget.staff?.position);
  }

  @override
  void dispose() {
    _first.dispose();
    _last.dispose();
    _position.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (!_form.currentState!.validate()) return;
    if (_delete &&
        !await confirm(
          context,
          title: 'Сохранить и удалить сотрудника?',
          message: 'Сотрудник больше не будет доступен для редактирования ЭМК.',
          action: 'Сохранить и удалить',
        )) {
      return;
    }
    setState(() {
      _saving = true;
      _error = null;
    });
    try {
      final staff = Staff.fromJson(
        await widget.session.api.request(
          widget.staff == null ? 'POST' : 'PUT',
          '/api/v1/staff${widget.staff == null ? '' : '/${widget.staff!.id}'}',
          body: {
            'first_name': _first.text.trim(),
            'last_name': _last.text.trim(),
            'position': _position.text.trim(),
          },
        ),
      );
      if (_delete) {
        await widget.session.api.request('DELETE', '/api/v1/staff/${staff.id}');
      }
      if (!mounted) return;
      notify(
        context,
        _delete
            ? 'Сотрудник помечен удалённым'
            : 'Карточка сотрудника сохранена',
      );
      Navigator.pop(context);
    } catch (error) {
      if (mounted) setState(() => _error = error);
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  @override
  Widget build(BuildContext context) => PageFrame(
    title: widget.staff == null ? 'Новый медработник' : 'Карточка медработника',
    subtitle: 'Личные данные и должность',
    back: _saving ? null : () => Navigator.pop(context),
    backLabel: 'К списку медработников',
    child: ListView(
      children: [
        Align(
          alignment: Alignment.topLeft,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 900),
            child: Panel(
              child: Form(
                key: _form,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(
                          Icons.badge_outlined,
                          color: MedicalColors.teal,
                        ),
                        const SizedBox(width: 12),
                        const Expanded(
                          child: Text(
                            'Данные сотрудника',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ),
                        StatusPill(
                          _archived ? 'Удалён' : 'Активен',
                          tone: _archived
                              ? MedicalColors.red
                              : MedicalColors.green,
                        ),
                      ],
                    ),
                    const SizedBox(height: 26),
                    FormError(_error),
                    FormFields(
                      children: [
                        TextFormField(
                          controller: _last,
                          enabled: !_archived && !_saving,
                          decoration: const InputDecoration(
                            labelText: 'Фамилия',
                          ),
                          validator: (v) => validateText(v, 3, 24),
                        ),
                        TextFormField(
                          controller: _first,
                          enabled: !_archived && !_saving,
                          decoration: const InputDecoration(labelText: 'Имя'),
                          validator: (v) => validateText(v, 3, 24),
                        ),
                        TextFormField(
                          controller: _position,
                          enabled: !_archived && !_saving,
                          decoration: InputDecoration(
                            labelText: 'Должность',
                            suffixIcon: PopupMenuButton<String>(
                              tooltip: 'Выбрать должность',
                              enabled: !_archived && !_saving,
                              icon: const Icon(Icons.expand_more_rounded),
                              onSelected: (value) => _position.text = value,
                              itemBuilder: (_) =>
                                  [
                                        'Терапевт',
                                        'Хирург',
                                        'Педиатр',
                                        'Невролог',
                                        'Медсестра',
                                      ]
                                      .map(
                                        (p) => PopupMenuItem(
                                          value: p,
                                          child: Text(p),
                                        ),
                                      )
                                      .toList(),
                            ),
                          ),
                          validator: (v) => validateText(v, 3, 24),
                        ),
                      ],
                    ),
                    const SizedBox(height: 22),
                    if (!_archived && widget.staff != null)
                      CheckboxListTile(
                        contentPadding: EdgeInsets.zero,
                        title: const Text('Пометить как удалённого'),
                        subtitle: const Text(
                          'Сотрудник сохранится в истории медицинских карт',
                        ),
                        value: _delete,
                        onChanged: _saving
                            ? null
                            : (value) =>
                                  setState(() => _delete = value ?? false),
                        controlAffinity: ListTileControlAffinity.leading,
                      ),
                    const SizedBox(height: 22),
                    Row(
                      children: [
                        if (!_archived)
                          FilledButton.icon(
                            onPressed: _saving ? null : _save,
                            icon: Icon(
                              _saving
                                  ? Icons.hourglass_top_rounded
                                  : Icons.check_rounded,
                              size: 18,
                            ),
                            label: Text(_saving ? 'Сохранение…' : 'Сохранить'),
                          ),
                        const SizedBox(width: 12),
                        OutlinedButton(
                          onPressed: _saving
                              ? null
                              : () => Navigator.pop(context),
                          child: Text(
                            _archived ? 'К списку медработников' : 'Отмена',
                          ),
                        ),
                      ],
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
