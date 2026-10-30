import HISAPI
import SwiftUI

struct StaffEditorView: View {
  let gateway: Gateway
  let staff: Staff?
  let saved: () -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var firstName: String
  @State private var lastName: String
  @State private var role: StaffRole
  @State private var position: String
  @State private var isBusy = false
  @State private var errorText: String?

  init(gateway: Gateway, staff: Staff?, saved: @escaping () -> Void) {
    self.gateway = gateway
    self.staff = staff
    self.saved = saved
    _firstName = State(initialValue: staff?.firstName ?? "")
    _lastName = State(initialValue: staff?.lastName ?? "")
    _position = State(initialValue: staff?.position ?? "")
    _role = State(initialValue: staff?.role ?? .other)
  }

  var body: some View {
    VStack(spacing: 0) {
      Form {
        Section(staff == nil ? "Новый медработник" : "Медработник") {
          TextField("Фамилия", text: $lastName)
          TextField("Имя", text: $firstName)
          TextField("Должность", text: $position)
          Picker("Роль", selection: $role) {
            ForEach(StaffRole.allCases, id: \.self) { Text($0.label).tag($0) }
          }
          if let staff {
            LabeledContent("Статус", value: staff.deletedAt == nil ? "Активен" : "Удалён")
          }
        }
      }.formStyle(.grouped).disabled(isBusy)
      ErrorNotice(message: errorText)
      FormActions(isBusy: isBusy, cancel: { dismiss() }, save: { Task { await save() } })
    }
    .frame(width: 480)
    .interactiveDismissDisabled()
  }

  private func save() async {
    isBusy = true
    defer { isBusy = false }
    do {
      let input = StaffInput(firstName: firstName, lastName: lastName, position: position, role: role)
      if let staff {
        _ = try await gateway.staff.updateStaff(id: staff.id, staffInput: input)
      } else {
        _ = try await gateway.staff.createStaff(staffInput: input)
      }
      saved()
      dismiss()
    } catch {
      errorText = errorMessage(error)
    }
  }
}
