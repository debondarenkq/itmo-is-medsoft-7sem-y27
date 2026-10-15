import HISAPI
import SwiftUI

struct RecordHistoryView: View {
  let gateway: Gateway
  let recordID: UUID
  let version: Int64
  @State private var events: [Event] = []
  @State private var offset = 0
  @State private var total = 0
  @State private var moment = Date()
  @State private var snapshot: Record?
  @State private var errorText: String?
  @State private var isBusy = false

  var body: some View {
    VStack {
      HStack {
        DatePicker(
          "Состояние на", selection: $moment, displayedComponents: [.date, .hourAndMinute])
        Button("Посмотреть") { Task { await showState(at: moment) } }
        Button("Обновить", systemImage: "arrow.clockwise") { Task { await load() } }
      }.padding()
      ErrorNotice(message: errorText)
      List(events) { event in
        DisclosureGroup {
          LabeledContent("До изменения", value: describe(event.before))
          LabeledContent("После изменения", value: describe(event.after))
          Button("Открыть версию \(event.version)") {
            Task { await showState(version: event.version) }
          }
        } label: {
          VStack(alignment: .leading) {
            Text(event.type.label)
            Text("\(event.actor.lastName) \(event.actor.firstName) · \(event.actor.position)")
            Text(event.occurredAt, format: .dateTime.day().month().year().hour().minute().second())
              .foregroundStyle(.secondary)
          }
        }
        .textSelection(.enabled)
      }
      PageControls(offset: $offset, total: total, count: events.count)
    }
    .disabled(isBusy)
    .overlay { if isBusy { ProgressView() } }
    .task(id: "\(version)|\(offset)") { await load() }
    .sheet(item: $snapshot) { RecordSnapshotView(record: $0) }
  }

  private func load() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let page = try await gateway.records.getRecordHistory(id: recordID, limit: 50, offset: offset)
      try Task.checkCancellation()
      events = page.items
      total = page.total
      errorText = nil
    } catch {
      if !Task.isCancelled {
        errorText = errorMessage(error)
      }
    }
  }

  private func showState(at: Date? = nil, version: Int64? = nil) async {
    isBusy = true
    defer { isBusy = false }
    do {
      snapshot = try await gateway.records.getRecordState(id: recordID, at: at, version: version)
      errorText = nil
    } catch {
      errorText = errorMessage(error)
    }
  }

  private func describe(_ snapshot: [String: JSONValue]?) -> String {
    guard let snapshot else { return "Запись отсутствует" }
    var lines: [String] = []
    for key in snapshot.keys.sorted() {
      guard let text = snapshot[key]?.stringValue else { continue }
      let label: String
      switch key {
      case "text": label = "Назначение"
      case "name": label = "Наименование"
      case "code": label = "Код"
      case "status": label = "Статус"
      default: continue
      }

      var displayedValue = text
      if key == "status" {
        displayedValue = text == "active" ? "Активно" : "Отменено"
      }
      lines.append("\(label): \(displayedValue)")
    }
    return lines.joined(separator: "\n")
  }
}

extension EventType {
  fileprivate var label: String {
    switch self {
    case .recordCreated: return "Создана медицинская карта"
    case .diagnosisAdded: return "Добавлен диагноз"
    case .diagnosisRemoved: return "Удалён диагноз"
    case .prescriptionAdded: return "Добавлено назначение"
    case .prescriptionCancelled: return "Отменено назначение"
    case .prescriptionRemoved: return "Удалено назначение"
    }
  }
}

private struct RecordSnapshotView: View {
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
