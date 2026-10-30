import HISAPI
import SwiftUI

private enum AuditPeriod: String, CaseIterable {
  case all = "За всё время", today = "Сегодня", week = "Последние 7 дней", custom = "Выбрать период"
}

struct AuditView: View {
  let gateway: Gateway
  @State private var events: [AuditEvent] = []
  @State private var eventType: AuditEventKind?
  @State private var period = AuditPeriod.all
  @State private var fromDate = Calendar.current.startOfDay(for: Date())
  @State private var toDate = Date()
  @State private var offset = 0
  @State private var total = 0
  @State private var followLatest = true
  @State private var errorText: String?

  private var startDate: Date? {
    switch period {
    case .all: return nil
    case .today: return Calendar.current.startOfDay(for: Date())
    case .week: return Calendar.current.date(byAdding: .day, value: -7, to: Date())
    case .custom: return fromDate
    }
  }
  private var endDate: Date? { period == .custom ? toDate : nil }

  var body: some View {
    VStack(spacing: 0) {
      VStack(alignment: .leading, spacing: 12) {
        HStack {
          Picker("Событие", selection: $eventType) {
            Text("Все события").tag(Optional<AuditEventKind>.none)
            ForEach(AuditEventKind.allCases, id: \.self) { Text($0.label).tag(Optional.some($0)) }
          }.frame(maxWidth: 380)
          Picker("Период", selection: $period) {
            ForEach(AuditPeriod.allCases, id: \.self) { Text($0.rawValue).tag($0) }
          }.frame(maxWidth: 300)
          Spacer()
        }
        if period == .custom {
          HStack {
            DatePicker("С", selection: $fromDate)
            DatePicker("По", selection: $toDate)
          }
        }
        HStack {
          Label("Обновление каждую секунду", systemImage: "arrow.triangle.2.circlepath")
            .font(.caption).foregroundStyle(.secondary)
          Spacer()
          Toggle("Последние события", isOn: $followLatest)
        }
      }.padding()
      Divider()
      ErrorNotice(message: errorText)
      Table(events) {
        TableColumn("Время") { Text($0.eventTime, format: .dateTime.day().month().year().hour().minute().second()) }
        TableColumn("Событие") { event in
          Label(event.eventType.label, systemImage: event.eventType.icon).foregroundStyle(event.message.hasPrefix("Отклонено:") ? Color.red : (event.message.hasPrefix("Запрошено:") ? Color.secondary : event.eventType.color))
        }
        TableColumn("Автор") { Text($0.actorName).help($0.actor) }
        TableColumn("Описание", value: \.message)
      }.textSelection(.enabled)
      if followLatest {
        HStack { Text("Последние \(events.count) из \(total) событий").foregroundStyle(.secondary); Spacer() }.padding()
      } else { PageControls(offset: $offset, total: total, count: events.count) }
    }
    .navigationTitle("Аудит")
    .task(id: "\(eventType?.rawValue ?? "all")|\(period)|\(fromDate)|\(toDate)|\(followLatest)|\(followLatest ? 0 : offset)") {
      while !Task.isCancelled {
        await load()
        try? await Task.sleep(for: .seconds(1))
      }
    }
    .onChange(of: eventType) { offset = 0 }
    .onChange(of: period) { offset = 0 }
    .onChange(of: fromDate) { offset = 0 }
    .onChange(of: toDate) { offset = 0 }
  }

  private func load() async {
    do {
      var pageOffset = offset
      if followLatest {
        let count = try await gateway.security.listAudit(limit: 1, offset: 0, eventType: eventType, from: startDate, to: endDate)
        pageOffset = max(0, count.total - 50)
      }
      let page = try await gateway.security.listAudit(limit: 50, offset: pageOffset, eventType: eventType, from: startDate, to: endDate)
      try Task.checkCancellation()
      events = page.items; total = page.total; errorText = nil
    } catch { if !Task.isCancelled { errorText = errorMessage(error) } }
  }
}

extension AuditEvent: Identifiable { public var id: UUID { eventId } }
