import HISAPI
import SwiftUI

struct WorkerPermissionsView: View {
  let worker: Staff

  var body: some View {
    GroupBox("Ваши права") {
      VStack(alignment: .leading, spacing: 12) {
        Text(worker.role.label).font(.headline)
        Text(worker.fullName).foregroundStyle(.secondary)
        Divider()
        Label("Просмотр ЭМК и истории", systemImage: "checkmark.circle")
        if worker.role.canCreate {
          Label("Создание пациентов и ЭМК", systemImage: "checkmark.circle")
        }
        if worker.role == .doctor {
          Label("Разблокировка ЭМК", systemImage: "checkmark.circle")
          Label("Изменение диагнозов и назначений", systemImage: "checkmark.circle")
          Text("Чтобы внести изменения, сначала разблокируйте карту.")
            .font(.caption).foregroundStyle(.secondary)
        } else {
          Text("Изменение содержимого ЭМК доступно врачу.")
            .font(.caption).foregroundStyle(.secondary)
        }
      }.frame(maxWidth: .infinity, alignment: .leading).padding(8)
    }.frame(width: 210)
  }
}
