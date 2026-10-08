import 'package:flutter/material.dart';
import '../core/api_client.dart';
import '../core/theme.dart';

class PageFrame extends StatelessWidget {
  const PageFrame({
    super.key,
    required this.title,
    required this.subtitle,
    required this.child,
    this.actions = const [],
    this.back,
    this.backLabel = 'Назад к списку',
  });
  final String title, subtitle, backLabel;
  final Widget child;
  final List<Widget> actions;
  final VoidCallback? back;
  @override
  Widget build(BuildContext context) => Material(
    color: MedicalColors.background,
    child: SafeArea(
      child: Padding(
        padding: const EdgeInsets.all(28),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (back != null)
              Padding(
                padding: const EdgeInsets.only(bottom: 12),
                child: TextButton.icon(
                  onPressed: back,
                  icon: const Icon(Icons.arrow_back_rounded, size: 18),
                  label: Text(backLabel),
                ),
              ),
            LayoutBuilder(
              builder: (context, constraints) {
                final heading = Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: const TextStyle(
                        fontSize: 27,
                        fontWeight: FontWeight.w700,
                        letterSpacing: -.5,
                      ),
                    ),
                    const SizedBox(height: 7),
                    Text(
                      subtitle,
                      style: const TextStyle(
                        color: MedicalColors.muted,
                        fontSize: 13,
                      ),
                    ),
                  ],
                );
                if (constraints.maxWidth < 850 && actions.isNotEmpty) {
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      heading,
                      const SizedBox(height: 16),
                      Wrap(spacing: 10, runSpacing: 10, children: actions),
                    ],
                  );
                }
                return Row(
                  children: [
                    Expanded(child: heading),
                    if (actions.isNotEmpty)
                      Wrap(spacing: 10, runSpacing: 10, children: actions),
                  ],
                );
              },
            ),
            const SizedBox(height: 24),
            Expanded(child: child),
          ],
        ),
      ),
    ),
  );
}

class Panel extends StatelessWidget {
  const Panel({
    super.key,
    required this.child,
    this.padding = const EdgeInsets.all(24),
  });
  final Widget child;
  final EdgeInsetsGeometry padding;
  @override
  Widget build(BuildContext context) => Card(
    margin: EdgeInsets.zero,
    clipBehavior: Clip.antiAlias,
    child: Padding(padding: padding, child: child),
  );
}

class StatusPill extends StatelessWidget {
  const StatusPill(this.label, {super.key, this.tone = MedicalColors.green});
  final String label;
  final Color tone;
  @override
  Widget build(BuildContext context) => Container(
    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
    decoration: BoxDecoration(
      color: tone.withValues(alpha: .09),
      borderRadius: BorderRadius.circular(20),
    ),
    child: Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 6,
          height: 6,
          decoration: BoxDecoration(color: tone, shape: BoxShape.circle),
        ),
        const SizedBox(width: 7),
        Text(
          label,
          style: TextStyle(
            color: tone,
            fontSize: 12,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    ),
  );
}

class ErrorPanel extends StatelessWidget {
  const ErrorPanel(this.error, {super.key, required this.retry});
  final Object error;
  final VoidCallback retry;
  @override
  Widget build(BuildContext context) => Center(
    child: ConstrainedBox(
      constraints: const BoxConstraints(maxWidth: 520),
      child: Panel(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.cloud_off_rounded,
              size: 40,
              color: MedicalColors.muted,
            ),
            const SizedBox(height: 16),
            const Text(
              'Не удалось загрузить данные',
              style: TextStyle(fontWeight: FontWeight.w600, fontSize: 18),
            ),
            const SizedBox(height: 10),
            SelectableText(errorMessage(error), textAlign: TextAlign.center),
            const SizedBox(height: 20),
            OutlinedButton.icon(
              onPressed: retry,
              icon: const Icon(Icons.refresh_rounded),
              label: const Text('Попробовать снова'),
            ),
          ],
        ),
      ),
    ),
  );
}

class EmptyPanel extends StatelessWidget {
  const EmptyPanel({
    super.key,
    required this.title,
    required this.subtitle,
    this.icon = Icons.folder_open_rounded,
  });
  final String title, subtitle;
  final IconData icon;
  @override
  Widget build(BuildContext context) => Center(
    child: Padding(
      padding: const EdgeInsets.all(32),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            padding: const EdgeInsets.all(20),
            decoration: const BoxDecoration(
              color: MedicalColors.mint,
              shape: BoxShape.circle,
            ),
            child: Icon(icon, size: 34, color: MedicalColors.teal),
          ),
          const SizedBox(height: 18),
          Text(
            title,
            style: const TextStyle(fontSize: 19, fontWeight: FontWeight.w600),
          ),
          const SizedBox(height: 8),
          Text(
            subtitle,
            textAlign: TextAlign.center,
            style: const TextStyle(color: MedicalColors.muted),
          ),
        ],
      ),
    ),
  );
}

String errorMessage(Object error) => error is ApiException
    ? error.toString()
    : error is FormatException
    ? error.message
    : 'Не удалось выполнить действие. Попробуйте ещё раз.';
void notify(BuildContext context, String text, {bool error = false}) =>
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(text),
        backgroundColor: error ? MedicalColors.red : MedicalColors.ink,
        behavior: SnackBarBehavior.floating,
      ),
    );

Future<bool> confirm(
  BuildContext context, {
  required String title,
  required String message,
  String action = 'Удалить',
  bool destructive = true,
}) async =>
    await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(title),
        content: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 420),
          child: Text(message),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Отмена'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            style: destructive
                ? FilledButton.styleFrom(backgroundColor: MedicalColors.red)
                : null,
            child: Text(action),
          ),
        ],
      ),
    ) ??
    false;

String? validateText(String? value, int min, int max, {bool optional = false}) {
  final text = (value ?? '').trim();
  if (optional && text.isEmpty) return null;
  if (text.runes.length < min || text.runes.length > max) {
    return 'От $min до $max символов';
  }
  return null;
}

class FormFields extends StatelessWidget {
  const FormFields({super.key, required this.children});
  final List<Widget> children;
  @override
  Widget build(BuildContext context) => LayoutBuilder(
    builder: (context, constraints) {
      final width = constraints.maxWidth >= 620
          ? (constraints.maxWidth - 20) / 2
          : constraints.maxWidth;
      return Wrap(
        spacing: 20,
        runSpacing: 22,
        children: children
            .map((child) => SizedBox(width: width, child: child))
            .toList(),
      );
    },
  );
}

class FormError extends StatelessWidget {
  const FormError(this.error, {super.key});
  final Object? error;
  @override
  Widget build(BuildContext context) => error == null
      ? const SizedBox.shrink()
      : Container(
          width: double.infinity,
          margin: const EdgeInsets.only(bottom: 20),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: MedicalColors.red.withValues(alpha: .06),
            borderRadius: BorderRadius.circular(10),
          ),
          child: Text(
            errorMessage(error!),
            style: const TextStyle(color: MedicalColors.red),
          ),
        );
}
