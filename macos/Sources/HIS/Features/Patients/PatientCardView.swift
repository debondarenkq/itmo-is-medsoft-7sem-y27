import HISAPI
import SwiftUI

struct PatientCardView: View {
  let gateway: Gateway
  let worker: Staff
  @State var patient: Patient
  @Environment(\.dismiss) private var dismiss
  @State private var editing = false
  @State private var showingRecord = false
  @State private var deleting = false
  @State private var isBusy = false
  @State private var errorText: String?

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
      ErrorNotice(message: errorText)
      HStack {
        Spacer()
        Button("К списку пациентов") { dismiss() }.keyboardShortcut(.cancelAction)
      }.padding()
    }
    .frame(width: 560)
    .disabled(isBusy)
    .overlay { if isBusy { ProgressView() } }
    .sheet(isPresented: $editing) {
      PatientEditorView(gateway: gateway, patient: patient) { Task { await reload() } }
    }
    .sheet(isPresented: $showingRecord, onDismiss: { Task { await reload() } }) {
      RecordView(gateway: gateway, worker: worker, patient: patient)
    }
    .confirmationDialog("Удалить пациента?", isPresented: $deleting, titleVisibility: .visible) {
      Button("Удалить", role: .destructive) {
        Task {
          isBusy = true
          defer { isBusy = false }
          do {
            try await gateway.patients.deletePatient(id: patient.id)
            dismiss()
          } catch {
            errorText = errorMessage(error)
          }
        }
      }
    }
  }

  private func reload() async {
    isBusy = true
    defer { isBusy = false }
    do {
      patient = try await gateway.patients.getPatient(id: patient.id)
      errorText = nil
    } catch {
      errorText = errorMessage(error)
    }
  }
}
