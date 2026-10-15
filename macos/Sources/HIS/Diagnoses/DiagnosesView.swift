import HISAPI
import SwiftUI
import UniformTypeIdentifiers

struct DiagnosesView: View {
  let gateway: Gateway
  @State private var diagnoses: [Diagnosis] = []
  @State private var selectedDiagnosisID: UUID?
  @State private var search = ""
  @State private var offset = 0
  @State private var total = 0
  @State private var isBusy = false
  @State private var errorText: String?
  @State private var importResult: String?
  @State private var importing = false
  @State private var editor: DiagnosisEditorSelection?
  @State private var deleting = false
  private var selectedDiagnosis: Diagnosis? { diagnoses.first { $0.id == selectedDiagnosisID } }

  var body: some View {
    VStack(spacing: 0) {
      ErrorNotice(message: errorText)
      if let importResult { Text(importResult).foregroundStyle(.secondary).padding() }
      Table(diagnoses, selection: $selectedDiagnosisID) {
        TableColumn("Код", value: \.code)
        TableColumn("Наименование", value: \.name)
      }
      .contextMenu(forSelectionType: UUID.self) { _ in
      } primaryAction: { ids in
        if let value = diagnoses.first(where: { ids.contains($0.id) }) { edit(value) }
      }
      PageControls(offset: $offset, total: total, count: diagnoses.count)
    }
    .navigationTitle("Диагнозы")
    .searchable(text: $search, prompt: "Код или наименование")
    .toolbar {
      Button("Импортировать JSON", systemImage: "square.and.arrow.down") { importing = true }
      Button("Добавить", systemImage: "plus") { edit(nil) }
      Button("Изменить", systemImage: "square.and.pencil") { edit(selectedDiagnosis) }.disabled(
        selectedDiagnosis == nil)
      Button("Удалить", systemImage: "trash") { deleting = true }.disabled(selectedDiagnosis == nil)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .disabled(isBusy)
    .overlay { if isBusy { ProgressView() } }
    .task(id: "\(search)|\(offset)") { await load() }
    .onChange(of: search) { offset = 0 }
    .fileImporter(isPresented: $importing, allowedContentTypes: [.json]) { result in
      Task {
        do {
          let file = try result.get()
          await importDiagnoses(from: file)
        } catch {
          errorText = errorMessage(error)
        }
      }
    }
    .sheet(item: $editor) { selection in
      DiagnosisEditor(gateway: gateway, diagnosis: selection.value) { Task { await load() } }
    }
    .confirmationDialog("Удалить диагноз?", isPresented: $deleting, titleVisibility: .visible) {
      Button("Удалить", role: .destructive) { Task { await delete() } }
    } message: {
      Text(selectedDiagnosis?.name ?? "")
    }
  }

  private func importDiagnoses(from file: URL) async {
    isBusy = true
    defer { isBusy = false }

    let accessed = file.startAccessingSecurityScopedResource()
    defer {
      if accessed { file.stopAccessingSecurityScopedResource() }
    }

    do {
      let contents = try Data(contentsOf: file)
      let result = try await gateway.diagnoses.importDiagnoses(contents: contents)
      importResult = "Импортировано записей: \(result.items.count)"
      await load()
    } catch {
      errorText = errorMessage(error)
    }
  }

  private func edit(_ diagnosis: Diagnosis?) {
    editor = DiagnosisEditorSelection(value: diagnosis)
  }

  private func load() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let page = try await gateway.diagnoses.listDiagnoses(
        limit: 50, offset: offset, search: search)
      try Task.checkCancellation()
      diagnoses = page.items
      total = page.total
      errorText = nil
    } catch {
      if !Task.isCancelled {
        errorText = errorMessage(error)
      }
    }
  }

  private func delete() async {
    guard let selectedDiagnosis else { return }
    isBusy = true
    do {
      try await gateway.diagnoses.deleteDiagnosis(id: selectedDiagnosis.id)
      await load()
    } catch {
      errorText = errorMessage(error)
    }
    isBusy = false
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
  @State private var code: String
  @State private var name: String
  @State private var isBusy = false
  @State private var errorText: String?

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
      }.formStyle(.grouped).disabled(isBusy)
      ErrorNotice(message: errorText)
      FormActions(isBusy: isBusy, cancel: { dismiss() }, save: { Task { await save() } })
    }
    .frame(width: 480)
    .interactiveDismissDisabled()
  }

  private func save() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let input = DiagnosisInput(code: code, name: name)
      if let diagnosis {
        _ = try await gateway.diagnoses.updateDiagnosis(id: diagnosis.id, diagnosisInput: input)
      } else {
        _ = try await gateway.diagnoses.createDiagnosis(diagnosisInput: input)
      }
      saved()
      dismiss()
    } catch {
      errorText = errorMessage(error)
    }
  }
}
