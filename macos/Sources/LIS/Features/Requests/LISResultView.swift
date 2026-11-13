import HISAPI
import SwiftUI

struct LISResultView: View {
  let gateway: Gateway
  let requestID: UUID
  var snapshot: LISHistory? = nil
  @Environment(\.dismiss) private var dismiss
  @State private var details: LISRequestDetails?
  @State private var result = ""
  @State private var errorText: String?
  @State private var busy = false

  private var readOnly: Bool { details?.history != nil || details?.request.status == .completed }

  var body: some View {
    VStack(alignment: .leading, spacing: 18) {
      Text(readOnly ? "Сведения об анализе · Только просмотр" : "Проведение анализа").font(.title2)
      ErrorNotice(message: errorText)
      if let details {
        ScrollView {
          VStack(alignment: .leading, spacing: 18) {
            GroupBox("Пациент") {
              VStack(alignment: .leading, spacing: 8) {
                Text(details.patient.fullName).font(.headline)
                if let birthDate = details.patient.birthDate {
                  Text(birthDate.calendarDate, format: .dateTime.day().month().year())
                }
                Text("ID: \(details.patient.id.uuidString)").font(.caption).textSelection(.enabled)
              }.frame(maxWidth: .infinity, alignment: .leading).padding(8)
            }
            GroupBox("Назначение") {
              Grid(alignment: .leading, horizontalSpacing: 18, verticalSpacing: 10) {
                GridRow { Text("Анализ").foregroundStyle(.secondary); Text(details.request.testName) }
                GridRow { Text("LOINC").foregroundStyle(.secondary); Text(details.request.loincCode) }
                GridRow { Text("Назначил").foregroundStyle(.secondary); Text(details.practitioner.fullName) }
                GridRow { Text("Медработник ID").foregroundStyle(.secondary); Text(details.practitioner.id.uuidString).font(.caption).textSelection(.enabled) }
                GridRow { Text("Дата назначения").foregroundStyle(.secondary); Text(details.request.authoredOn, format: .dateTime.day().month().year().hour().minute()) }
                GridRow { Text("Запрос ID").foregroundStyle(.secondary); Text(details.request.id.uuidString).font(.caption).textSelection(.enabled) }
                if let note = details.request.note?.map(\.text).joined(separator: "\n"), !note.isEmpty {
                  GridRow { Text("Комментарий").foregroundStyle(.secondary); Text(note) }
                }
              }.frame(maxWidth: .infinity, alignment: .leading).padding(8)
            }
            if let observation = details.history?.result ?? details.result {
              GroupBox("Зафиксированный результат") {
                VStack(alignment: .leading, spacing: 10) {
                  Text(observation.valueString).font(.title2).textSelection(.enabled)
                  Text(observation.effectiveDateTime, format: .dateTime.day().month().year().hour().minute().second())
                  if let history = details.history {
                    Label(history.delivered ? "Доставлен в HIS" : "Доставка в HIS не подтверждена", systemImage: history.delivered ? "checkmark.circle" : "exclamationmark.triangle")
                      .foregroundStyle(history.delivered ? Color.green : Color.orange)
                  } else {
                    Label("Результат получен из HIS", systemImage: "checkmark.circle").foregroundStyle(.green)
                  }
                  Text("Результат ID: \(observation.id.uuidString)").font(.caption).textSelection(.enabled)
                  Text("Повторная фиксация этого назначения запрещена.").font(.caption).foregroundStyle(.secondary)
                }.frame(maxWidth: .infinity, alignment: .leading).padding(8)
              }
            } else if !readOnly {
              TextField("Фактический результат", text: $result, axis: .vertical).lineLimit(2...4)
              Text("Дата проведения — момент фиксации результата в LIS.").font(.caption).foregroundStyle(.secondary)
            }
            Text("Единицы измерения не используются.").font(.caption).foregroundStyle(.secondary)
          }
        }
      } else if busy { ProgressView() }
      HStack {
        Button("Закрыть") { dismiss() }.keyboardShortcut(.cancelAction)
        Spacer()
        if !readOnly && details != nil {
          Button("Зафиксировать результат") { Task { await save() } }
            .buttonStyle(.borderedProminent).disabled(result.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty)
        }
      }
    }.padding(24).frame(width: 760, height: 700).disabled(busy)
    .task {
      if let snapshot {
        details = LISRequestDetails(history: snapshot, patient: snapshot.patient, practitioner: snapshot.practitioner, request: snapshot.request)
      } else { await load() }
    }
  }

  private func load() async {
    busy = true
    defer { busy = false }
    do { details = try await gateway.lis.getLISRequest(id: requestID) }
    catch { errorText = errorMessage(error) }
  }

  private func save() async {
    busy = true
    defer { busy = false }
    do {
      let history = try await gateway.lis.recordLISResult(lISResultInput: LISResultInput(requestId: requestID, value: result))
      details = LISRequestDetails(history: history, patient: history.patient, practitioner: history.practitioner, request: history.request)
      errorText = nil
    } catch {
      errorText = errorMessage(error)
      await load()
    }
  }
}
