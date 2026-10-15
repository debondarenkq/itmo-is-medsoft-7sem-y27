import HISAPI
import SwiftUI

struct RecordDiagnosesView: View {
  let gateway: Gateway
  @Binding var draft: RecordDraft
  let editing: Bool
  let add: () -> Void
  @State private var events: [Event] = []
  @State private var removed: [RecordDiagnosis] = []
  @State private var selection: UUID?
  @State private var message: String?
  @State private var loading = false

  private var rows: [RecordDiagnosis] {
    var result = draft.diagnoses
    for diagnosis in draft.base.state.diagnoses + removed
    where !result.contains(where: { $0.id == diagnosis.id }) {
      result.append(diagnosis)
    }
    return result
  }

  private var selected: RecordDiagnosis? { rows.first { $0.id == selection } }

  var body: some View {
    VStack {
      ErrorNotice(message: message)
      Table(rows, selection: $selection) {
        TableColumn("Дата") { diagnosis in
          if let event = events.first(where: {
            $0.entityId == diagnosis.id && $0.type == .diagnosisAdded
          }) {
            Text(event.occurredAt, format: .dateTime.day().month().year())
          } else {
            Text(
              draft.base.state.diagnoses.contains(where: { $0.id == diagnosis.id })
                ? "—" : "Черновик"
            ).foregroundStyle(.secondary)
          }
        }
        TableColumn("Код", value: \.code)
        TableColumn("Диагноз", value: \.name)
        TableColumn("Статус") { diagnosis in
          Text(isPresent(diagnosis) ? "Активен" : "Удалён")
            .foregroundStyle(isPresent(diagnosis) ? .primary : .secondary)
        }
      }
      HStack {
        if editing {
          Button("Добавить диагноз…", systemImage: "plus", action: add)
          Button("Удалить", systemImage: "minus") {
            draft.diagnoses.removeAll { $0.id == selection }
          }.disabled(selected.map { !isPresent($0) } ?? true)
          Button("Восстановить", systemImage: "arrow.uturn.backward") {
            guard let selected else { return }
            if draft.base.state.diagnoses.contains(where: { $0.id == selected.id }) {
              draft.diagnoses.append(selected)
            } else {
              draft.diagnoses.append(
                RecordDiagnosis(
                  code: selected.code, diagnosisId: selected.diagnosisId,
                  id: UUID(), name: selected.name))
            }
            selection = nil
          }.disabled(selected.map { isPresent($0) } ?? true)
        }
        Spacer()
        if loading { ProgressView().controlSize(.small) }
        Button("Обновить историю диагнозов", systemImage: "arrow.clockwise") {
          Task { await load() }
        }
        .disabled(loading)
      }.padding()
    }
    .task(id: draft.base.version) { await load() }
  }

  private func isPresent(_ diagnosis: RecordDiagnosis) -> Bool {
    draft.diagnoses.contains { $0.id == diagnosis.id }
  }

  private func load() async {
    loading = true
    defer { loading = false }
    do {
      var history: [Event] = []
      var offset = 0
      while true {
        let page = try await gateway.records.getRecordHistory(
          id: draft.base.id, limit: 200, offset: offset)
        try Task.checkCancellation()
        history += page.items.filter { $0.version <= draft.base.version }
        offset += page.items.count
        if page.items.isEmpty || page.offset + page.items.count >= page.total { break }
      }
      var snapshots: [RecordDiagnosis] = []
      for event in history where event.type == .diagnosisRemoved {
        if let before = event.before {
          let data = try JSONEncoder().encode(before)
          snapshots.append(try JSONDecoder().decode(RecordDiagnosis.self, from: data))
        }
      }
      events = history
      removed = snapshots
      message = nil
    } catch { if !Task.isCancelled { message = errorMessage(error) } }
  }
}
