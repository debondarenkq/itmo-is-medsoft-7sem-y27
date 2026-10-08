import Cocoa
import FlutterMacOS

class MainFlutterWindow: NSWindow {
  override func awakeFromNib() {
    let flutterViewController = FlutterViewController()
    self.contentViewController = flutterViewController
    self.title = "МедКарта — HIS"
    self.minSize = NSSize(width: 1080, height: 760)
    self.setContentSize(NSSize(width: 1280, height: 840))
    self.center()

    RegisterGeneratedPlugins(registry: flutterViewController)

    super.awakeFromNib()
  }
}
