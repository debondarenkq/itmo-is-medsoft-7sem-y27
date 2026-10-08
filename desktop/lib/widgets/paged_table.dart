import 'dart:async';
import 'package:flutter/material.dart';
import '../core/models.dart';
import '../core/theme.dart';
import 'common.dart';

class PagedTable<T> extends StatefulWidget {
  const PagedTable({
    super.key,
    required this.fetch,
    required this.columns,
    required this.row,
    required this.emptyTitle,
    required this.emptySubtitle,
    this.searchHint = 'Поиск',
    this.extraFilter,
  });
  final Future<PageResult<T>> Function(int offset, String query) fetch;
  final List<DataColumn> columns;
  final DataRow Function(T) row;
  final String emptyTitle, emptySubtitle, searchHint;
  final Widget? extraFilter;
  @override
  State<PagedTable<T>> createState() => PagedTableState<T>();
}

class PagedTableState<T> extends State<PagedTable<T>> {
  final _search = TextEditingController();
  final _vertical = ScrollController(), _horizontal = ScrollController();
  Timer? _debounce;
  PageResult<T>? _page;
  Object? _error;
  int _offset = 0, _generation = 0;
  bool _loading = true;
  @override
  void initState() {
    super.initState();
    reload();
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _search.dispose();
    _vertical.dispose();
    _horizontal.dispose();
    super.dispose();
  }

  Future<void> reload({bool reset = false}) async {
    if (reset) _offset = 0;
    final generation = ++_generation;
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final result = await widget.fetch(_offset, _search.text.trim());
      if (!mounted || generation != _generation) return;
      if (result.items.isEmpty && _offset > 0 && result.total <= _offset) {
        _offset = ((result.total - 1).clamp(0, result.total) ~/ 50) * 50;
        await reload();
        return;
      }
      setState(() => _page = result);
    } catch (error) {
      if (mounted && generation == _generation) setState(() => _error = error);
    } finally {
      if (mounted && generation == _generation) {
        setState(() => _loading = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) => Panel(
    padding: EdgeInsets.zero,
    child: Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(20),
          child: Row(
            children: [
              Expanded(
                child: TextField(
                  controller: _search,
                  decoration: InputDecoration(
                    hintText: widget.searchHint,
                    prefixIcon: const Icon(Icons.search_rounded, size: 21),
                    isDense: true,
                    suffixIcon: _search.text.isEmpty
                        ? null
                        : IconButton(
                            tooltip: 'Очистить поиск',
                            onPressed: () {
                              _search.clear();
                              reload(reset: true);
                            },
                            icon: const Icon(Icons.close_rounded, size: 18),
                          ),
                  ),
                  onChanged: (_) {
                    setState(() {});
                    _debounce?.cancel();
                    _debounce = Timer(
                      const Duration(milliseconds: 350),
                      () => reload(reset: true),
                    );
                  },
                ),
              ),
              if (widget.extraFilter != null) ...[
                const SizedBox(width: 16),
                widget.extraFilter!,
              ],
              const SizedBox(width: 10),
              IconButton(
                onPressed: _loading ? null : () => reload(),
                tooltip: 'Обновить список',
                icon: const Icon(Icons.refresh_rounded),
              ),
            ],
          ),
        ),
        if (_loading)
          const LinearProgressIndicator(minHeight: 2)
        else
          const Divider(height: 1),
        Expanded(
          child: _error != null
              ? ErrorPanel(_error!, retry: reload)
              : _page == null
              ? const Center(child: CircularProgressIndicator())
              : _page!.items.isEmpty
              ? EmptyPanel(
                  title: _search.text.isEmpty
                      ? widget.emptyTitle
                      : 'Ничего не найдено',
                  subtitle: _search.text.isEmpty
                      ? widget.emptySubtitle
                      : 'Попробуйте изменить поисковый запрос.',
                )
              : LayoutBuilder(
                  builder: (context, constraints) => Scrollbar(
                    controller: _vertical,
                    thumbVisibility: true,
                    child: SingleChildScrollView(
                      controller: _vertical,
                      child: Scrollbar(
                        controller: _horizontal,
                        thumbVisibility: true,
                        notificationPredicate: (notification) =>
                            notification.metrics.axis == Axis.horizontal,
                        child: SingleChildScrollView(
                          controller: _horizontal,
                          scrollDirection: Axis.horizontal,
                          child: ConstrainedBox(
                            constraints: BoxConstraints(
                              minWidth: constraints.maxWidth,
                            ),
                            child: DataTable(
                              showCheckboxColumn: false,
                              columns: widget.columns,
                              rows: _page!.items.map(widget.row).toList(),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ),
        ),
        const Divider(height: 1),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 22, vertical: 10),
          child: Row(
            children: [
              Text(
                _page == null ? 'Загрузка…' : 'Всего записей: ${_page!.total}',
                style: const TextStyle(
                  fontSize: 12,
                  color: MedicalColors.muted,
                ),
              ),
              const Spacer(),
              Text(
                '${_page == null || _page!.total == 0 ? 0 : _offset + 1}–${_offset + (_page?.items.length ?? 0)}',
                style: const TextStyle(fontSize: 12),
              ),
              const SizedBox(width: 12),
              IconButton(
                tooltip: 'Предыдущая страница',
                onPressed: _offset == 0 || _loading
                    ? null
                    : () {
                        _offset -= 50;
                        reload();
                      },
                icon: const Icon(Icons.chevron_left_rounded),
              ),
              IconButton(
                tooltip: 'Следующая страница',
                onPressed: _loading || _offset + 50 >= (_page?.total ?? 0)
                    ? null
                    : () {
                        _offset += 50;
                        reload();
                      },
                icon: const Icon(Icons.chevron_right_rounded),
              ),
            ],
          ),
        ),
      ],
    ),
  );
}
