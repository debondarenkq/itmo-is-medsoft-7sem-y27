import SwiftUI

struct PageControls: View {
  @Binding var offset: Int
  let total: Int
  let count: Int
  var limit = 50

  var body: some View {
    HStack {
      Text(total == 0 ? "Нет записей" : (count == 0 ? "На этой странице нет записей · всего \(total)" : "\(offset + 1)–\(offset + count) из \(total)"))
        .foregroundStyle(.secondary)
      Spacer()
      Button("Предыдущая", systemImage: "chevron.left") { offset = max(0, offset - limit) }
        .disabled(offset == 0)
      Button("Следующая", systemImage: "chevron.right") { offset += limit }
        .disabled(offset + count >= total)
    }
    .padding()
  }
}
