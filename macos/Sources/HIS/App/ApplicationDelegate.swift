import AppKit

@MainActor
final class ApplicationDelegate: NSObject, NSApplicationDelegate {
  func applicationShouldTerminate(_ sender: NSApplication) -> NSApplication.TerminateReply {
    guard sender.windows.contains(where: { $0.attachedSheet != nil }) else { return .terminateNow }
    let alert = NSAlert()
    alert.messageText = "Завершить работу HIS?"
    alert.informativeText = "Открыт редактор. Несохранённые изменения будут потеряны."
    alert.addButton(withTitle: "Продолжить работу")
    alert.addButton(withTitle: "Завершить")
    return alert.runModal() == .alertSecondButtonReturn ? .terminateNow : .terminateCancel
  }
}
