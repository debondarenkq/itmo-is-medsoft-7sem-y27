import HISAPI
import SwiftUI

struct PatientCard: View {
  let gateway: Gateway
  let worker: Staff
  @State var patient: Patient
  @Environment(\.dismiss) private var dismiss
  @State private var editing = false
  @State private var showingRecord = false
  @State private var deleting = false
  @State private var busy = false
  @State private var message: String?

  var body: some View {
    VStack(spacing: 0) {
      Form {
        Section("Пациент") {
          LabeledContent("Фамилия", value: patient.lastName)
          LabeledContent("Имя", value: patient.firstName)
          LabeledContent("Отчество", value: patient.middleName ?? "—")
          LabeledContent(
            "Дата рождения",
            value: patient.birthDate.calendarDate.formatted(date: .numeric, time: .omitted))
          LabeledContent("Пол", value: patient.administrativeSex.label)
          LabeledContent("Комментарий", value: patient.comment)
          LabeledContent("ЭМК", value: patient.hasRecord ? "Создана" : "Нет")
        }
        Section {
          Button(
            patient.hasRecord ? "Открыть ЭМК" : "Создать ЭМК…", systemImage: "heart.text.clipboard"
          ) {
            showingRecord = true
          }
          Button("Изменить данные", systemImage: "square.and.pencil") { editing = true }
          Button("Удалить пациента", role: .destructive) { deleting = true }
        }
      }
      .formStyle(.grouped)
      .textSelection(.enabled)
      ErrorNotice(message: message)
      HStack {
        Spacer()
        Button("К списку пациентов") { dismiss() }.keyboardShortcut(.cancelAction)
      }.padding()
    }
    .frame(width: 560)
    .disabled(busy)
    .overlay { if busy { ProgressView() } }
    .sheet(isPresented: $editing) {
      PatientEditor(gateway: gateway, patient: patient) { Task { await reload() } }
    }
    .sheet(isPresented: $showingRecord, onDismiss: { Task { await reload() } }) {
      RecordView(gateway: gateway, worker: worker, patient: patient)
    }
    .confirmationDialog("Удалить пациента?", isPresented: $deleting, titleVisibility: .visible) {
      Button("Удалить", role: .destructive) {
        Task {
          busy = true
          defer { busy = false }
          do {
            try await gateway.patients.deletePatient(id: patient.id)
            dismiss()
          } catch { message = errorMessage(error) }
        }
      }
    }
  }

  private func reload() async {
    busy = true
    defer { busy = false }
    do {
      patient = try await gateway.patients.getPatient(id: patient.id)
      message = nil
    } catch { message = errorMessage(error) }
  }
}
