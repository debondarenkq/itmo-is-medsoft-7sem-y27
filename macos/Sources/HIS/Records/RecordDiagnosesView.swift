import HISAPI
import SwiftUI

struct RecordDiagnosesView: View {
  let gateway: Gateway
  let draft: RecordDraft
  let isEditing: Bool
  let addDiagnosis: () -> Void
  let removeDiagnosis: (UUID) -> Void
  let restoreDiagnosis: (RecordDiagnosis) -> Void
  @State private var events: [Event] = []
  @State private var removed: [RecordDiagnosis] = []
  @State private var selectedDiagnosisID: UUID?
  @State private var errorText: String?
  @State private var isLoading = false

  private var rows: [RecordDiagnosis] {
    var result = draft.diagnoses
    for diagnosis in draft.base.state.diagnoses + removed {
      let alreadyListed = result.contains { row in row.id == diagnosis.id }
      if !alreadyListed { result.append(diagnosis) }
    }
    return result
  }

  private var selectedDiagnosis: RecordDiagnosis? { rows.first { $0.id == selectedDiagnosisID } }

  private var canRemoveDiagnosis: Bool {
    guard let selectedDiagnosis else { return false }
    return isPresent(selectedDiagnosis)
  }

  private var canRestoreDiagnosis: Bool {
    guard let selectedDiagnosis else { return false }
    return !isPresent(selectedDiagnosis)
  }

  var body: some View {
    VStack {
      ErrorNotice(message: errorText)
      Table(rows, selection: $selectedDiagnosisID) {
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
        if isEditing {
          Button("Добавить диагноз…", systemImage: "plus", action: addDiagnosis)
          Button("Удалить", systemImage: "minus") {
            if let selectedDiagnosis { removeDiagnosis(selectedDiagnosis.id) }
          }.disabled(!canRemoveDiagnosis)
          Button("Восстановить", systemImage: "arrow.uturn.backward") {
            if let selectedDiagnosis {
              restoreDiagnosis(selectedDiagnosis)
              selectedDiagnosisID = nil
            }
          }.disabled(!canRestoreDiagnosis)
        }
        Spacer()
        if isLoading { ProgressView().controlSize(.small) }
        Button("Обновить историю диагнозов", systemImage: "arrow.clockwise") {
          Task { await load() }
        }
        .disabled(isLoading)
      }.padding()
    }
    .task(id: draft.base.version) { await load() }
  }

  private func isPresent(_ diagnosis: RecordDiagnosis) -> Bool {
    draft.diagnoses.contains { $0.id == diagnosis.id }
  }

  private func load() async {
    isLoading = true
    defer { isLoading = false }
    do {
      var history: [Event] = []
      var offset = 0
      while true {
        let page = try await gateway.records.getRecordHistory(
          id: draft.base.id, limit: 200, offset: offset)
        try Task.checkCancellation()
        for event in page.items {
          if event.version <= draft.base.version { history.append(event) }
        }
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
      errorText = nil
    } catch {
      if !Task.isCancelled {
        errorText = errorMessage(error)
      }
    }
  }
}
