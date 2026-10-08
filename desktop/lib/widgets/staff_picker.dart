import 'package:flutter/material.dart';
import '../core/models.dart';
import '../core/session.dart';
import 'common.dart';

Future<Staff?> selectStaff(
  BuildContext context,
  AppSession session, {
  String title = 'Выберите медработника',
  String action = 'Выбрать',
}) => showDialog<Staff>(
  context: context,
  builder: (context) =>
      _StaffPicker(session: session, title: title, action: action),
);

class _StaffPicker extends StatefulWidget {
  const _StaffPicker({
    required this.session,
    required this.title,
    required this.action,
  });
  final AppSession session;
  final String title, action;
  @override
  State<_StaffPicker> createState() => _StaffPickerState();
}

class _StaffPickerState extends State<_StaffPicker> {
  List<Staff>? _staff;
  Staff? _selected;
  Object? _error;
  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _error = null;
      _staff = null;
    });
    try {
      final staff = await widget.session.api.staff();
      if (!mounted) return;
      setState(() {
        _staff = staff;
        _selected = staff
            .where((s) => s.id == widget.session.worker?.id)
            .firstOrNull;
      });
    } catch (error) {
      if (mounted) setState(() => _error = error);
    }
  }

  @override
  Widget build(BuildContext context) => AlertDialog(
    title: Text(widget.title),
    content: SizedBox(
      width: 460,
      child: _error != null
          ? Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(errorMessage(_error!)),
                TextButton(onPressed: _load, child: const Text('Повторить')),
              ],
            )
          : _staff == null
          ? const SizedBox(
              height: 100,
              child: Center(child: CircularProgressIndicator()),
            )
          : _staff!.isEmpty
          ? const Text(
              'Нет активных медработников. Добавьте сотрудника в разделе администратора.',
            )
          : Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('Сотрудник будет указан автором изменений карты.'),
                const SizedBox(height: 20),
                DropdownButtonFormField<Staff>(
                  initialValue: _selected,
                  isExpanded: true,
                  decoration: const InputDecoration(labelText: 'Медработник'),
                  items: _staff!
                      .map(
                        (s) => DropdownMenuItem(
                          value: s,
                          child: Text(
                            '${s.name} · ${s.position}',
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                      )
                      .toList(),
                  onChanged: (value) => setState(() => _selected = value),
                ),
              ],
            ),
    ),
    actions: [
      TextButton(
        onPressed: () => Navigator.pop(context),
        child: const Text('Отмена'),
      ),
      FilledButton(
        onPressed: _selected == null
            ? null
            : () {
                widget.session.selectWorker(_selected);
                Navigator.pop(context, _selected);
              },
        child: Text(widget.action),
      ),
    ],
  );
}
