import HISAPI
import SwiftUI
import UniformTypeIdentifiers

struct DiagnosesView: View {
  let gateway: Gateway
  @State private var items: [Diagnosis] = []
  @State private var selection: UUID?
  @State private var search = ""
  @State private var offset = 0
  @State private var total = 0
  @State private var busy = false
  @State private var message: String?
  @State private var importResult: String?
  @State private var importing = false
  @State private var editor: DiagnosisEditorSelection?
  @State private var deleting = false
  private var selected: Diagnosis? { items.first { $0.id == selection } }

  var body: some View {
    VStack(spacing: 0) {
      ErrorNotice(message: message)
      if let importResult { Text(importResult).foregroundStyle(.secondary).padding() }
      Table(items, selection: $selection) {
        TableColumn("Код", value: \.code)
        TableColumn("Наименование", value: \.name)
      }
      .contextMenu(forSelectionType: UUID.self) { _ in
      } primaryAction: { ids in
        if let value = items.first(where: { ids.contains($0.id) }) { edit(value) }
      }
      PageControls(offset: $offset, total: total, count: items.count)
    }
    .navigationTitle("Диагнозы")
    .searchable(text: $search, prompt: "Код или наименование")
    .toolbar {
      Button("Импортировать JSON", systemImage: "square.and.arrow.down") { importing = true }
      Button("Добавить", systemImage: "plus") { edit(nil) }
      Button("Изменить", systemImage: "square.and.pencil") { edit(selected) }.disabled(
        selected == nil)
      Button("Удалить", systemImage: "trash") { deleting = true }.disabled(selected == nil)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .disabled(busy)
    .overlay { if busy { ProgressView() } }
    .task(id: "\(search)|\(offset)") { await load() }
    .onChange(of: search) { offset = 0 }
    .fileImporter(isPresented: $importing, allowedContentTypes: [.json]) { result in
      Task {
        busy = true
        defer { busy = false }
        do {
          let count = try await gateway.importDiagnoses(from: result.get())
          importResult = "Импортировано записей: \(count)"
          await load()
        } catch { message = errorMessage(error) }
      }
    }
    .sheet(item: $editor) { selection in
      DiagnosisEditor(gateway: gateway, diagnosis: selection.value) { Task { await load() } }
    }
    .confirmationDialog("Удалить диагноз?", isPresented: $deleting, titleVisibility: .visible) {
      Button("Удалить", role: .destructive) { Task { await delete() } }
    } message: {
      Text(selected?.name ?? "")
    }
  }

  private func edit(_ diagnosis: Diagnosis?) {
    editor = DiagnosisEditorSelection(value: diagnosis)
  }

  private func load() async {
    busy = true
    defer { busy = false }
    do {
      let page = try await gateway.diagnoses.listDiagnoses(
        limit: 50, offset: offset, search: search)
      try Task.checkCancellation()
      items = page.items
      total = page.total
      message = nil
    } catch { if !Task.isCancelled { message = errorMessage(error) } }
  }

  private func delete() async {
    guard let selected else { return }
    busy = true
    do {
      try await gateway.diagnoses.deleteDiagnosis(id: selected.id)
      await load()
    } catch { message = errorMessage(error) }
    busy = false
  }
}

private struct DiagnosisEditorSelection: Identifiable {
  let id = UUID()
  let value: Diagnosis?
}

private struct DiagnosisEditor: View {
  let gateway: Gateway
  let diagnosis: Diagnosis?
  let saved: () -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var code = ""
  @State private var name = ""
  @State private var busy = false
  @State private var message: String?

  init(gateway: Gateway, diagnosis: Diagnosis?, saved: @escaping () -> Void) {
    self.gateway = gateway
    self.diagnosis = diagnosis
    self.saved = saved
    _code = State(initialValue: diagnosis?.code ?? "")
    _name = State(initialValue: diagnosis?.name ?? "")
  }

  var body: some View {
    VStack(spacing: 0) {
      Form {
        Section(diagnosis == nil ? "Новый диагноз" : "Диагноз") {
          TextField("Код", text: $code)
          TextField("Наименование", text: $name)
        }
      }.formStyle(.grouped).disabled(busy)
      ErrorNotice(message: message)
      FormActions(busy: busy, cancel: { dismiss() }, save: { Task { await save() } })
    }
    .frame(width: 480)
    .interactiveDismissDisabled()
  }

  private func save() async {
    busy = true
    defer { busy = false }
    do {
      let input = DiagnosisInput(code: code, name: name)
      if let diagnosis {
        _ = try await gateway.diagnoses.updateDiagnosis(id: diagnosis.id, diagnosisInput: input)
      } else {
        _ = try await gateway.diagnoses.createDiagnosis(diagnosisInput: input)
      }
      saved()
      dismiss()
    } catch { message = errorMessage(error) }
  }
}
