import HISAPI
import SwiftUI

struct StaffView: View {
  let gateway: Gateway
  @State private var staffMembers: [Staff] = []
  @State private var selectedStaffID: UUID?
  @State private var search = ""
  @State private var includeDeleted = false
  @State private var offset = 0
  @State private var total = 0
  @State private var isBusy = false
  @State private var errorText: String?
  @State private var editor: StaffEditorSelection?
  @State private var deleting = false
  private var selectedStaff: Staff? { staffMembers.first { $0.id == selectedStaffID } }

  var body: some View {
    VStack(spacing: 0) {
      ErrorNotice(message: errorText)
      Table(staffMembers, selection: $selectedStaffID) {
        TableColumn("Фамилия и имя", value: \.fullName)
        TableColumn("Должность", value: \.position)
        TableColumn("Статус") { staff in
          Text(staff.deletedAt == nil ? "Активен" : "Удалён")
            .foregroundStyle(staff.deletedAt == nil ? .green : .red)
        }
      }
      .contextMenu(forSelectionType: UUID.self) { _ in
      } primaryAction: { ids in
        if let value = staffMembers.first(where: { ids.contains($0.id) }) { edit(value) }
      }
      PageControls(offset: $offset, total: total, count: staffMembers.count)
    }
    .navigationTitle("Медработники")
    .searchable(text: $search, prompt: "Фамилия, имя или должность")
    .toolbar {
      Toggle("Показать удалённых", isOn: $includeDeleted)
      Button("Добавить", systemImage: "plus") { edit(nil) }
      Button("Изменить", systemImage: "square.and.pencil") { edit(selectedStaff) }.disabled(
        selectedStaff == nil)
      Button("Удалить", systemImage: "trash") { deleting = true }.disabled(selectedStaff == nil)
      Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
    }
    .disabled(isBusy)
    .overlay { if isBusy { ProgressView() } }
    .task(id: "\(search)|\(offset)|\(includeDeleted)") { await load() }
    .onChange(of: search) { offset = 0 }
    .onChange(of: includeDeleted) { offset = 0 }
    .sheet(item: $editor) { selection in
      StaffEditorView(gateway: gateway, staff: selection.value) { Task { await load() } }
    }
    .confirmationDialog("Удалить медработника?", isPresented: $deleting, titleVisibility: .visible)
    {
      Button("Удалить", role: .destructive) { Task { await delete() } }
    } message: {
      Text(selectedStaff?.fullName ?? "")
    }
  }

  private func edit(_ staff: Staff?) {
    editor = StaffEditorSelection(value: staff)
  }

  private func load() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let page = try await gateway.staff.listStaff(
        limit: 50, offset: offset, search: search, includeDeleted: includeDeleted)
      try Task.checkCancellation()
      staffMembers = page.items
      total = page.total
      errorText = nil
    } catch {
      if !Task.isCancelled {
        errorText = errorMessage(error)
      }
    }
  }

  private func delete() async {
    guard let selectedStaff else { return }
    isBusy = true
    do {
      try await gateway.staff.deleteStaff(id: selectedStaff.id)
      await load()
    } catch {
      errorText = errorMessage(error)
    }
    isBusy = false
  }
}

private struct StaffEditorSelection: Identifiable {
  let id = UUID()
  let value: Staff?
}
