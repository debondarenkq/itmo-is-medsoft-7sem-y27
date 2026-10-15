import HISAPI
import SwiftUI

struct DirectoryPicker<Item: Identifiable & Sendable>: View where Item.ID == UUID {
  let title: String
  let label: (Item) -> String
  let fetch: (String, Int) async throws -> (items: [Item], total: Int)
  let select: (Item) -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var search = ""
  @State private var offset = 0
  @State private var items: [Item] = []
  @State private var total = 0
  @State private var selection: UUID?
  @State private var isBusy = false
  @State private var errorText: String?

  var body: some View {
    VStack {
      Text(title).font(.headline).padding(.top)
      TextField("Поиск", text: $search).textFieldStyle(.roundedBorder).padding(.horizontal)
      ErrorNotice(message: errorText)
      List(items, selection: $selection) { Text(label($0)).tag($0.id) }
        .overlay { if isBusy { ProgressView() } }
      PageControls(offset: $offset, total: total, count: items.count)
      HStack {
        Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
        Spacer()
        Button("Отмена", role: .cancel) { dismiss() }.keyboardShortcut(.cancelAction)
        Button("Выбрать") {
          if let item = items.first(where: { $0.id == selection }) {
            select(item)
            dismiss()
          }
        }
        .keyboardShortcut(.defaultAction)
        .disabled(selection == nil || isBusy)
      }.padding()
    }
    .frame(width: 520, height: 460)
    .task(id: "\(search)|\(offset)") { await load() }
    .onChange(of: search) {
      offset = 0
      selection = nil
    }
    .onChange(of: offset) { selection = nil }
  }

  private func load() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let result = try await fetch(search, offset)
      try Task.checkCancellation()
      items = result.items
      total = result.total
      errorText = nil
    } catch {
      if !Task.isCancelled {
        errorText = errorMessage(error)
      }
    }
  }
}
