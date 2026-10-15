import HISAPI
import SwiftUI

struct DirectoryPicker<Item: Identifiable & Sendable>: View where Item.ID == UUID {
  let title: String
  let label: (Item) -> String
  let fetch: (String, Int) async throws -> ([Item], Int)
  let select: (Item) -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var search = ""
  @State private var offset = 0
  @State private var items: [Item] = []
  @State private var total = 0
  @State private var selection: UUID?
  @State private var busy = false
  @State private var message: String?

  var body: some View {
    VStack {
      Text(title).font(.headline).padding(.top)
      TextField("Поиск", text: $search).textFieldStyle(.roundedBorder).padding(.horizontal)
      ErrorNotice(message: message)
      List(items, selection: $selection) { Text(label($0)).tag($0.id) }
        .overlay { if busy { ProgressView() } }
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
        .disabled(selection == nil || busy)
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
    busy = true
    defer { busy = false }
    do {
      let result = try await fetch(search, offset)
      try Task.checkCancellation()
      items = result.0
      total = result.1
      message = nil
    } catch { if !Task.isCancelled { message = errorMessage(error) } }
  }
}
