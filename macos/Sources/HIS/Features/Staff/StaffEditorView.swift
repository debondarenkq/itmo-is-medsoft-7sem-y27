import HISAPI
import SwiftUI

struct StaffEditorView: View {
  let gateway: Gateway
  let staff: Staff?
  let saved: () -> Void
  @Environment(\.dismiss) private var dismiss
  @State private var firstName: String
  @State private var lastName: String
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
  }

  var body: some View {
    VStack(spacing: 0) {
      Form {
        Section(staff == nil ? "Новый медработник" : "Медработник") {
          TextField("Фамилия", text: $lastName)
          TextField("Имя", text: $firstName)
          TextField("Должность", text: $position)
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
      let input = StaffInput(firstName: firstName, lastName: lastName, position: position)
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
