import HISAPI
import SwiftUI

struct StaffView: View {
  let gateway: Gateway
  @State private var items: [Staff] = []
  @State private var selection: UUID?
  @State private var search = ""
  @State private var includeDeleted = false
  @State private var offset = 0
  @State private var total = 0
  @State private var busy = false
  @State private var message: String?
  @State private var editor: StaffEditorSelection?
  @State private var deleting = false
  private var selected: Staff? { items.first { $0.id == selection } }

  var body: some View {
    VStack(spacing: 0) {
      ErrorNotice(message: message)
      Table(items, selection: $selection) {
        TableColumn("Фамилия и имя", value: \.fullName)
        TableColumn("Должность", value: \.position)
        TableColumn("Статус") { staff in
          Text(staff.deletedAt == nil ? "Активен" : "Удалён")
            .foregroundStyle(staff.deletedAt == nil ? .green : .red)
        }
      }
      .contextMenu(forSelectionType: UUID.self) { _ in
      } primaryAction: { ids in
        if let value = items.first(where: { ids.contains($0.id) }) { edit(value) }
      }
      PageControls(offset: $offset, total: total, count: items.count)
    }
    .navigationTitle("Медработники")
    .searchable(text: $search, prompt: "Фамилия, имя или должность")
    .toolbar {
      Toggle("Показать удалённых", isOn: $includeDeleted)
      Button("Добавить", systemImage: "plus") { edit(nil) }
      Button("Изменить", systemImage: "square.and.pencil") { edit(selected) }.disabled(
        selected == nil)
      Button("Удалить", systemImage: "trash") { deleting = true }.disabled(selected == nil)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .disabled(busy)
    .overlay { if busy { ProgressView() } }
    .task(id: "\(search)|\(offset)|\(includeDeleted)") { await load() }
    .onChange(of: search) { offset = 0 }
    .onChange(of: includeDeleted) { offset = 0 }
    .sheet(item: $editor) { selection in
      StaffEditor(gateway: gateway, staff: selection.value) { Task { await load() } }
    }
    .confirmationDialog("Удалить медработника?", isPresented: $deleting, titleVisibility: .visible)
    {
      Button("Удалить", role: .destructive) { Task { await delete() } }
    } message: {
      Text(selected?.fullName ?? "")
    }
  }

  private func edit(_ staff: Staff?) {
    editor = StaffEditorSelection(value: staff)
  }

  private func load() async {
    busy = true
    defer { busy = false }
    do {
      let page = try await gateway.staff.listStaff(
        limit: 50, offset: offset, search: search, includeDeleted: includeDeleted)
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
      try await gateway.staff.deleteStaff(id: selected.id)
      await load()
    } catch { message = errorMessage(error) }
    busy = false
  }
}

private struct StaffEditorSelection: Identifiable {
  let id = UUID()
  let value: Staff?
}

private struct StaffEditor: View {
  let gateway: Gateway
  let staff: Staff?
  let saved: () -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var firstName = ""
  @State private var lastName = ""
  @State private var position = ""
  @State private var busy = false
  @State private var message: String?

  init(gateway: Gateway, staff: Staff?, saved: @escaping () -> Void) {
    self.gateway = gateway
    self.staff = staff
    self.saved = saved
    _firstName = State(initialValue: staff?.firstName ?? "")
    _lastName = State(initialValue: staff?.lastName ?? "")
    _position = State(initialValue: staff?.position ?? "")
  }

  var body: some View {
    VStack(spacing: 0) {
      Form {
        Section(staff == nil ? "Новый медработник" : "Медработник") {
          TextField("Фамилия", text: $lastName)
          TextField("Имя", text: $firstName)
          TextField("Должность", text: $position)
          if let staff {
            LabeledContent("Статус", value: staff.deletedAt == nil ? "Активен" : "Удалён")
          }
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
      let input = StaffInput(firstName: firstName, lastName: lastName, position: position)
      if let staff {
        _ = try await gateway.staff.updateStaff(id: staff.id, staffInput: input)
      } else {
        _ = try await gateway.staff.createStaff(staffInput: input)
      }
      saved()
      dismiss()
    } catch { message = errorMessage(error) }
  }
}
