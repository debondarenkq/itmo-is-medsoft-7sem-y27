import 'dart:ui' show AppExitResponse;
import 'package:flutter/material.dart';
import '../core/session.dart';
import '../core/theme.dart';
import '../widgets/common.dart';
import '../widgets/staff_picker.dart';
import 'diagnosis_screen.dart';
import 'patient_screens.dart';
import 'staff_screens.dart';

enum Section { staff, diagnoses, patients }

class ApplicationShell extends StatefulWidget {
  const ApplicationShell({super.key, required this.session});
  final AppSession session;
  @override
  State<ApplicationShell> createState() => _ApplicationShellState();
}

class _ApplicationShellState extends State<ApplicationShell> {
  final _navigator = GlobalKey<NavigatorState>();
  Section _section = Section.staff;
  int _roleRevision = 0;
  late final AppLifecycleListener _lifecycle;
  bool get _admin => _section != Section.patients;
  @override
  void initState() {
    super.initState();
    _lifecycle = AppLifecycleListener(
      onExitRequested: () async => await widget.session.canNavigate()
          ? AppExitResponse.exit
          : AppExitResponse.cancel,
    );
  }

  @override
  void dispose() {
    _lifecycle.dispose();
    super.dispose();
  }

  Widget _screen(Section section) => switch (section) {
    Section.staff => StaffListScreen(session: widget.session),
    Section.diagnoses => DiagnosisScreen(session: widget.session),
    Section.patients => PatientListScreen(session: widget.session),
  };
  Future<void> _select(Section section) async {
    if (!await widget.session.canNavigate()) {
      if (mounted) setState(() => _roleRevision++);
      return;
    }
    if (!mounted) return;
    setState(() => _section = section);
    _navigator.currentState!.pushAndRemoveUntil(
      MaterialPageRoute<void>(builder: (_) => _screen(section)),
      (_) => false,
    );
  }

  Future<void> _connection() async {
    if (!await widget.session.canNavigate() || !mounted) return;
    final address = await showDialog<String>(
      context: context,
      builder: (_) => _ConnectionDialog(address: widget.session.api.baseUrl),
    );
    if (address == null || !mounted) return;
    widget.session.api.configure(address);
    widget.session.selectWorker(null);
    _navigator.currentState!.pushAndRemoveUntil(
      MaterialPageRoute<void>(builder: (_) => _screen(_section)),
      (_) => false,
    );
    notify(context, 'Адрес сервера обновлён');
  }

  Widget _destination(Section section, IconData icon, String text) {
    final selected = section == _section;
    return Padding(
      padding: const EdgeInsets.only(bottom: 6),
      child: Material(
        color: selected ? MedicalColors.mint : Colors.transparent,
        borderRadius: BorderRadius.circular(10),
        child: InkWell(
          borderRadius: BorderRadius.circular(10),
          onTap: () => _select(section),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 15),
            child: Row(
              children: [
                Icon(
                  icon,
                  color: selected ? MedicalColors.teal : MedicalColors.muted,
                  size: 21,
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    text,
                    style: TextStyle(
                      color: selected ? MedicalColors.teal : MedicalColors.ink,
                      fontWeight: selected ? FontWeight.w600 : FontWeight.w400,
                      fontSize: 14,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) => Scaffold(
    body: Row(
      children: [
        Container(
          width: 236,
          decoration: const BoxDecoration(
            color: Colors.white,
            border: Border(right: BorderSide(color: MedicalColors.line)),
          ),
          child: SafeArea(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(18, 28, 18, 18),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Container(
                        width: 42,
                        height: 42,
                        decoration: BoxDecoration(
                          color: MedicalColors.teal,
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: const Icon(
                          Icons.add_rounded,
                          color: Colors.white,
                          size: 33,
                        ),
                      ),
                      const SizedBox(width: 12),
                      const Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'МедКарта',
                              style: TextStyle(
                                fontSize: 21,
                                fontWeight: FontWeight.w700,
                                letterSpacing: -.5,
                              ),
                            ),
                            Text(
                              'Госпитальная система',
                              style: TextStyle(
                                fontSize: 10,
                                color: MedicalColors.muted,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 34),
                  const Padding(
                    padding: EdgeInsets.only(left: 12, bottom: 12),
                    child: Text(
                      'РАБОЧЕЕ МЕСТО',
                      style: TextStyle(
                        fontSize: 10,
                        fontWeight: FontWeight.w600,
                        letterSpacing: 1.4,
                        color: MedicalColors.muted,
                      ),
                    ),
                  ),
                  DropdownButtonFormField<bool>(
                    key: ValueKey((_admin, _roleRevision)),
                    initialValue: _admin,
                    isExpanded: true,
                    decoration: const InputDecoration(
                      isDense: true,
                      contentPadding: EdgeInsets.symmetric(
                        horizontal: 12,
                        vertical: 12,
                      ),
                    ),
                    items: const [
                      DropdownMenuItem(
                        value: true,
                        child: Text(
                          'Администратор',
                          style: TextStyle(fontSize: 13),
                        ),
                      ),
                      DropdownMenuItem(
                        value: false,
                        child: Text(
                          'Медработник',
                          style: TextStyle(fontSize: 13),
                        ),
                      ),
                    ],
                    onChanged: (admin) {
                      if (admin != null) {
                        _select(admin ? Section.staff : Section.patients);
                      }
                    },
                  ),
                  const SizedBox(height: 30),
                  if (_admin) ...[
                    _destination(
                      Section.staff,
                      Icons.badge_outlined,
                      'Медработники',
                    ),
                    const Padding(
                      padding: EdgeInsets.fromLTRB(14, 20, 0, 12),
                      child: Text(
                        'СПРАВОЧНИКИ',
                        style: TextStyle(
                          fontSize: 10,
                          fontWeight: FontWeight.w600,
                          letterSpacing: 1.4,
                          color: MedicalColors.muted,
                        ),
                      ),
                    ),
                    _destination(
                      Section.diagnoses,
                      Icons.medical_information_outlined,
                      'Диагнозы',
                    ),
                  ] else
                    _destination(
                      Section.patients,
                      Icons.people_outline_rounded,
                      'Пациенты',
                    ),
                  const Spacer(),
                  const Divider(),
                  const SizedBox(height: 14),
                  ListenableBuilder(
                    listenable: widget.session,
                    builder: (context, _) => Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            CircleAvatar(
                              radius: 18,
                              backgroundColor: MedicalColors.background,
                              child: Icon(
                                _admin
                                    ? Icons.admin_panel_settings_outlined
                                    : Icons.person_outline_rounded,
                                color: MedicalColors.teal,
                                size: 21,
                              ),
                            ),
                            const SizedBox(width: 10),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    _admin
                                        ? 'Администратор'
                                        : widget.session.worker?.name ??
                                              'Медработник',
                                    maxLines: 2,
                                    style: const TextStyle(
                                      fontSize: 12,
                                      fontWeight: FontWeight.w600,
                                    ),
                                  ),
                                  Text(
                                    _admin
                                        ? 'Управление учреждением'
                                        : widget.session.worker?.position ??
                                              'Сотрудник не выбран',
                                    style: const TextStyle(
                                      fontSize: 10,
                                      color: MedicalColors.muted,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                        if (!_admin)
                          TextButton(
                            onPressed: () =>
                                selectStaff(context, widget.session),
                            child: const Text(
                              'Выбрать сотрудника',
                              style: TextStyle(fontSize: 12),
                            ),
                          ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextButton.icon(
                    onPressed: _connection,
                    icon: const Icon(Icons.settings_outlined, size: 17),
                    label: const Text(
                      'Подключение',
                      style: TextStyle(fontSize: 12),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        Expanded(
          child: Navigator(
            key: _navigator,
            onGenerateRoute: (_) =>
                MaterialPageRoute<void>(builder: (_) => _screen(_section)),
          ),
        ),
      ],
    ),
  );
}

class _ConnectionDialog extends StatefulWidget {
  const _ConnectionDialog({required this.address});
  final String address;
  @override
  State<_ConnectionDialog> createState() => _ConnectionDialogState();
}

class _ConnectionDialogState extends State<_ConnectionDialog> {
  final _form = GlobalKey<FormState>();
  late final TextEditingController _address;
  @override
  void initState() {
    super.initState();
    _address = TextEditingController(text: widget.address);
  }

  @override
  void dispose() {
    _address.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => AlertDialog(
    title: const Text('Подключение к серверу'),
    content: SizedBox(
      width: 460,
      child: Form(
        key: _form,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Адрес единой точки входа госпитальной системы.'),
            const SizedBox(height: 20),
            TextFormField(
              controller: _address,
              decoration: const InputDecoration(
                labelText: 'Адрес сервера',
                hintText: 'http://localhost:8080',
              ),
              validator: (value) {
                final uri = Uri.tryParse(value?.trim() ?? '');
                return uri == null ||
                        !['http', 'https'].contains(uri.scheme) ||
                        uri.host.isEmpty ||
                        uri.userInfo.isNotEmpty ||
                        uri.hasQuery ||
                        uri.hasFragment
                    ? 'Укажите полный адрес HTTP или HTTPS'
                    : null;
              },
            ),
          ],
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
            Navigator.pop(context, _address.text.trim());
          }
        },
        child: const Text('Подключиться'),
      ),
    ],
  );
}
