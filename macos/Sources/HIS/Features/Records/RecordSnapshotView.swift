import HISAPI
import SwiftUI

struct RecordSnapshotView: View {
  let record: Record
  @Environment(\.dismiss) private var dismiss

  var body: some View {
    VStack {
      HStack {
        Label(
          "Историческое состояние · версия \(record.version)", systemImage: "clock.arrow.circlepath"
        )
        .font(.headline)
        Spacer()
        Button("Закрыть") { dismiss() }.keyboardShortcut(.cancelAction)
      }.padding()
      Text(record.updatedAt, format: .dateTime.day().month().year().hour().minute().second())
        .foregroundStyle(.secondary)
      TabView {
        Table(record.state.diagnoses) {
          TableColumn("Код", value: \.code)
          TableColumn("Диагноз", value: \.name)
        }.tabItem { Text("Диагнозы") }
        Table(record.state.prescriptions) {
          TableColumn("Назначение", value: \.text)
          TableColumn("Статус") { Text($0.status == .active ? "Активно" : "Отменено") }
        }.tabItem { Text("Назначения") }
      }.padding()
    }
    .frame(width: 760, height: 520)
  }
}
