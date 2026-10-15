import HISAPI
import SwiftUI

struct PatientEditorView: View {
  let gateway: Gateway
  let patient: Patient?
  let saved: () -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var firstName: String
  @State private var lastName: String
  @State private var middleName: String
  @State private var birthDate: Date
  @State private var sex: AdministrativeSex
  @State private var comment: String
  @State private var isBusy = false
  @State private var errorText: String?

  init(gateway: Gateway, patient: Patient?, saved: @escaping () -> Void) {
    self.gateway = gateway
    self.patient = patient
    self.saved = saved
    _firstName = State(initialValue: patient?.firstName ?? "")
    _lastName = State(initialValue: patient?.lastName ?? "")
    _middleName = State(initialValue: patient?.middleName ?? "")
    _birthDate = State(initialValue: patient?.birthDate.calendarDate ?? Date())
    _sex = State(initialValue: patient?.administrativeSex ?? .unknown)
    _comment = State(initialValue: patient?.comment ?? "")
  }

  var body: some View {
    VStack(spacing: 0) {
      Form {
        Section(patient == nil ? "Новый пациент" : "Пациент") {
          TextField("Фамилия", text: $lastName)
          TextField("Имя", text: $firstName)
          TextField("Отчество", text: $middleName)
          DatePicker("Дата рождения", selection: $birthDate, displayedComponents: .date)
          Picker("Пол", selection: $sex) {
            ForEach(AdministrativeSex.allCases, id: \.self) { Text($0.label).tag($0) }
          }
          TextField("Комментарий", text: $comment, axis: .vertical)
        }
      }.formStyle(.grouped).disabled(isBusy)
      ErrorNotice(message: errorText)
      FormActions(isBusy: isBusy, cancel: { dismiss() }, save: { Task { await save() } })
    }
    .frame(width: 520)
    .interactiveDismissDisabled()
  }

  private func save() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let input = PatientInput(
        administrativeSex: sex,
        birthDate: OpenAPIDateWithoutTime(wrappedDate: birthDate, timezone: .current),
        comment: comment, firstName: firstName, lastName: lastName,
        middleName: middleName.isEmpty ? nil : middleName)
      if let patient {
        _ = try await gateway.patients.updatePatient(id: patient.id, patientInput: input)
      } else {
        _ = try await gateway.patients.createPatient(patientInput: input)
      }
      saved()
      dismiss()
    } catch {
      errorText = errorMessage(error)
    }
  }
}
