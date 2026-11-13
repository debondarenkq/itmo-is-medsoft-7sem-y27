import AppKit

@MainActor
func configureApplicationIcon() {
  guard let url = Bundle.main.url(forResource: "AppIcon", withExtension: "icns"),
    let icon = NSImage(contentsOf: url) else { return }
  NSApplication.shared.applicationIconImage = icon
}
