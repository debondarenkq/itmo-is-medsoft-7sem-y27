import HISAPI
import SwiftUI

struct LabTestsView: View {
  let gateway: Gateway
  @State private var tests: [LabTest] = []
  @State private var selectedID: UUID?
  @State private var editing: LabTest?
  @State private var creating = false
  @State private var errorText: String?
  @State private var busy = false
  private var selected: LabTest? { tests.first { $0.id == selectedID } }

  var body: some View {
    VStack {
      ErrorNotice(message: errorText)
      Table(tests, selection: $selectedID) {
        TableColumn("LOINC", value: \.code)
        TableColumn("Анализ", value: \.name)
        TableColumn("Статус") { Text($0.active ? "Активен" : "Отключён") }
      }
      Text("Единицы измерения не используются. Отключённые анализы сохраняются в истории назначений.")
        .font(.caption).foregroundStyle(.secondary).padding()
    }
    .navigationTitle("Лабораторные анализы")
    .toolbar {
      Button("Создать", systemImage: "plus") { creating = true }
      Button("Изменить", systemImage: "square.and.pencil") { editing = selected }.disabled(selected == nil)
      Button("Отключить", systemImage: "minus.circle") { Task { await deactivate() } }.disabled(selected?.active != true)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .disabled(busy)
    .task { await load() }
    .sheet(isPresented: $creating, onDismiss: { Task { await load() } }) { LabTestEditor(gateway: gateway, test: nil) }
    .sheet(item: $editing, onDismiss: { Task { await load() } }) { LabTestEditor(gateway: gateway, test: $0) }
  }

  private func load() async {
    do { tests = try await gateway.laboratory.listLabTests(includeInactive: true); errorText = nil }
    catch { errorText = errorMessage(error) }
  }

  private func deactivate() async {
    guard let selected else { return }
    busy = true
    defer { busy = false }
    do { _ = try await gateway.laboratory.deactivateLabTest(id: selected.id); await load() }
    catch { errorText = errorMessage(error) }
  }
}

private struct LabTestEditor: View {
  let gateway: Gateway
  let test: LabTest?
  @Environment(\.dismiss) private var dismiss
  @State private var name = ""
  @State private var code = ""
  @State private var errorText: String?
  @State private var busy = false

  var body: some View {
    VStack(alignment: .leading, spacing: 16) {
      Text(test == nil ? "Новый анализ" : "Изменить анализ").font(.title2)
      Form {
        TextField("Название", text: $name)
        TextField("Код LOINC", text: $code)
      }
      ErrorNotice(message: errorText)
      HStack {
        Button("Отмена") { dismiss() }.keyboardShortcut(.cancelAction)
        Spacer()
        Button("Сохранить") { Task { await save() } }.buttonStyle(.borderedProminent)
      }
    }.padding(24).frame(width: 480).disabled(busy)
    .onAppear { name = test?.name ?? ""; code = test?.code ?? "" }
  }

  private func save() async {
    busy = true
    defer { busy = false }
    do {
      let input = LabTestInput(code: code, name: name)
      if let test { _ = try await gateway.laboratory.updateLabTest(id: test.id, labTestInput: input) }
      else { _ = try await gateway.laboratory.createLabTest(labTestInput: input) }
      dismiss()
    } catch { errorText = errorMessage(error) }
  }
}
