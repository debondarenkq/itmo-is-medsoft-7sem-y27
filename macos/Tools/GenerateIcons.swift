import AppKit

struct IconTheme {
  let top: NSColor
  let bottom: NSColor
  let ink: NSColor
}

func color(_ red: CGFloat, _ green: CGFloat, _ blue: CGFloat, _ alpha: CGFloat = 1) -> NSColor {
  NSColor(srgbRed: red, green: green, blue: blue, alpha: alpha)
}

func rounded(_ rect: NSRect, _ radius: CGFloat) -> NSBezierPath {
  NSBezierPath(roundedRect: rect, xRadius: radius, yRadius: radius)
}

func gradient(_ path: NSBezierPath, top: NSColor, bottom: NSColor) {
  NSGradient(starting: bottom, ending: top)!.draw(in: path, angle: 90)
}

func shadowed(_ draw: () -> Void, blur: CGFloat = 20, offset: CGFloat = -12, opacity: CGFloat = 0.2) {
  NSGraphicsContext.saveGraphicsState()
  let shadow = NSShadow()
  shadow.shadowColor = NSColor.black.withAlphaComponent(opacity)
  shadow.shadowBlurRadius = blur
  shadow.shadowOffset = NSSize(width: 0, height: offset)
  shadow.set()
  draw()
  NSGraphicsContext.restoreGraphicsState()
}

func tile(_ theme: IconTheme) {
  let path = rounded(NSRect(x: 100, y: 100, width: 824, height: 824), 184)
  shadowed({ gradient(path, top: theme.top, bottom: theme.bottom) }, blur: 30, offset: -15, opacity: 0.22)
  NSGraphicsContext.saveGraphicsState()
  path.addClip()
  let glow = rounded(NSRect(x: 102, y: 104, width: 820, height: 816), 182)
  NSColor.white.withAlphaComponent(0.28).setStroke()
  glow.lineWidth = 3
  glow.stroke()
  NSGraphicsContext.restoreGraphicsState()
}

func paper(_ path: NSBezierPath) {
  shadowed {
    gradient(path, top: color(1, 1, 1), bottom: color(0.89, 0.94, 1))
  }
  NSColor.white.withAlphaComponent(0.85).setStroke()
  path.lineWidth = 4
  path.stroke()
}

func medicalRecord(_ theme: IconTheme) {
  paper(rounded(NSRect(x: 296, y: 248, width: 432, height: 508), 46))
  shadowed({
    gradient(rounded(NSRect(x: 428, y: 714, width: 168, height: 64), 19),
      top: color(0.99, 1, 1), bottom: color(0.64, 0.76, 0.91))
    gradient(NSBezierPath(ovalIn: NSRect(x: 485, y: 746, width: 54, height: 54)),
      top: color(0.99, 1, 1), bottom: color(0.71, 0.82, 0.95))
  }, blur: 9, offset: -5, opacity: 0.14)
  color(0.42, 0.56, 0.73).setFill()
  NSBezierPath(ovalIn: NSRect(x: 503, y: 764, width: 18, height: 18)).fill()
  theme.ink.setFill()
  rounded(NSRect(x: 490, y: 516, width: 44, height: 140), 9).fill()
  rounded(NSRect(x: 442, y: 564, width: 140, height: 44), 9).fill()
  theme.ink.withAlphaComponent(0.76).setFill()
  rounded(NSRect(x: 369, y: 421, width: 286, height: 27), 13).fill()
  theme.ink.withAlphaComponent(0.42).setFill()
  rounded(NSRect(x: 369, y: 357, width: 224, height: 27), 13).fill()
}

func administrator(_ theme: IconTheme) {
  let shield = NSBezierPath()
  shield.move(to: NSPoint(x: 512, y: 238))
  shield.curve(to: NSPoint(x: 300, y: 500), controlPoint1: NSPoint(x: 363, y: 310), controlPoint2: NSPoint(x: 300, y: 390))
  shield.line(to: NSPoint(x: 300, y: 690))
  shield.curve(to: NSPoint(x: 512, y: 770), controlPoint1: NSPoint(x: 384, y: 706), controlPoint2: NSPoint(x: 463, y: 740))
  shield.curve(to: NSPoint(x: 724, y: 690), controlPoint1: NSPoint(x: 561, y: 740), controlPoint2: NSPoint(x: 640, y: 706))
  shield.line(to: NSPoint(x: 724, y: 500))
  shield.curve(to: NSPoint(x: 512, y: 238), controlPoint1: NSPoint(x: 724, y: 390), controlPoint2: NSPoint(x: 661, y: 310))
  shield.close()
  paper(shield)
  theme.ink.setStroke()
  let ring = NSBezierPath(ovalIn: NSRect(x: 455, y: 552, width: 114, height: 114))
  ring.lineWidth = 30
  ring.stroke()
  theme.ink.setFill()
  rounded(NSRect(x: 497, y: 370, width: 30, height: 196), 8).fill()
  rounded(NSRect(x: 512, y: 394, width: 80, height: 30), 7).fill()
  rounded(NSRect(x: 512, y: 446, width: 62, height: 30), 7).fill()
}

func laboratory(_ theme: IconTheme) {
  let flask = NSBezierPath()
  flask.move(to: NSPoint(x: 452, y: 738))
  flask.line(to: NSPoint(x: 572, y: 738))
  flask.line(to: NSPoint(x: 572, y: 564))
  flask.line(to: NSPoint(x: 710, y: 332))
  flask.curve(to: NSPoint(x: 678, y: 268), controlPoint1: NSPoint(x: 733, y: 292), controlPoint2: NSPoint(x: 719, y: 268))
  flask.line(to: NSPoint(x: 346, y: 268))
  flask.curve(to: NSPoint(x: 314, y: 332), controlPoint1: NSPoint(x: 305, y: 268), controlPoint2: NSPoint(x: 291, y: 292))
  flask.line(to: NSPoint(x: 452, y: 564))
  flask.close()
  paper(flask)
  NSGraphicsContext.saveGraphicsState()
  flask.addClip()
  let liquid = NSBezierPath()
  liquid.move(to: NSPoint(x: 280, y: 418))
  liquid.curve(to: NSPoint(x: 744, y: 418), controlPoint1: NSPoint(x: 438, y: 470), controlPoint2: NSPoint(x: 591, y: 380))
  liquid.line(to: NSPoint(x: 744, y: 250))
  liquid.line(to: NSPoint(x: 280, y: 250))
  liquid.close()
  gradient(liquid, top: color(0.22, 0.86, 0.73), bottom: color(0.02, 0.48, 0.43))
  NSColor.white.withAlphaComponent(0.7).setFill()
  NSBezierPath(ovalIn: NSRect(x: 513, y: 333, width: 30, height: 30)).fill()
  NSColor.white.withAlphaComponent(0.46).setFill()
  NSBezierPath(ovalIn: NSRect(x: 578, y: 365, width: 19, height: 19)).fill()
  NSGraphicsContext.restoreGraphicsState()
  NSColor.white.withAlphaComponent(0.9).setStroke()
  flask.lineWidth = 9
  flask.stroke()
  gradient(rounded(NSRect(x: 434, y: 712, width: 156, height: 36), 12),
    top: NSColor.white, bottom: color(0.78, 0.95, 0.9))
  let reflection = NSBezierPath()
  reflection.move(to: NSPoint(x: 474, y: 680))
  reflection.line(to: NSPoint(x: 474, y: 558))
  reflection.line(to: NSPoint(x: 405, y: 442))
  NSColor.white.withAlphaComponent(0.88).setStroke()
  reflection.lineWidth = 13
  reflection.lineCapStyle = .round
  reflection.lineJoinStyle = .round
  reflection.stroke()
}

func render(_ name: String, size: Int) -> NSBitmapImageRep {
  let bitmap = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: size, pixelsHigh: size,
    bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false,
    colorSpaceName: .deviceRGB, bytesPerRow: 0, bitsPerPixel: 0)!
  NSGraphicsContext.saveGraphicsState()
  NSGraphicsContext.current = NSGraphicsContext(bitmapImageRep: bitmap)
  let context = NSGraphicsContext.current!.cgContext
  context.clear(CGRect(x: 0, y: 0, width: size, height: size))
  context.scaleBy(x: CGFloat(size) / 1024, y: CGFloat(size) / 1024)
  context.setShouldAntialias(true)
  let theme: IconTheme
  switch name {
  case "HIS": theme = IconTheme(top: color(0.24, 0.61, 1), bottom: color(0.06, 0.27, 0.79), ink: color(0.10, 0.40, 0.87))
  case "HISAdmin": theme = IconTheme(top: color(1, 0.71, 0.30), bottom: color(0.88, 0.36, 0.08), ink: color(0.86, 0.39, 0.09))
  default: theme = IconTheme(top: color(0.20, 0.80, 0.64), bottom: color(0.02, 0.43, 0.36), ink: color(0.02, 0.48, 0.40))
  }
  tile(theme)
  switch name {
  case "HIS": medicalRecord(theme)
  case "HISAdmin": administrator(theme)
  default: laboratory(theme)
  }
  NSGraphicsContext.restoreGraphicsState()
  return bitmap
}

let root = URL(fileURLWithPath: FileManager.default.currentDirectoryPath)
let sources = root.appendingPathComponent("macos/Resources/Icons")
let temporary = root.appendingPathComponent("macos/build/Icons")
try FileManager.default.createDirectory(at: sources, withIntermediateDirectories: true)
try FileManager.default.createDirectory(at: temporary, withIntermediateDirectories: true)
let sizes = [("16x16", 16), ("16x16@2x", 32), ("32x32", 32), ("32x32@2x", 64),
  ("128x128", 128), ("128x128@2x", 256), ("256x256", 256), ("256x256@2x", 512),
  ("512x512", 512), ("512x512@2x", 1024)]
for name in ["HIS", "HISAdmin", "LIS"] {
  let iconset = temporary.appendingPathComponent("\(name).iconset")
  try FileManager.default.createDirectory(at: iconset, withIntermediateDirectories: true)
  for (suffix, size) in sizes {
    let data = render(name, size: size).representation(using: .png, properties: [:])!
    try data.write(to: iconset.appendingPathComponent("icon_\(suffix).png"))
    if size == 1024 { try data.write(to: sources.appendingPathComponent("\(name).png")) }
  }
  let destination = root.appendingPathComponent("macos/Resources/\(name)")
  try FileManager.default.createDirectory(at: destination, withIntermediateDirectories: true)
  let process = Process()
  process.executableURL = URL(fileURLWithPath: "/usr/bin/iconutil")
  process.arguments = ["-c", "icns", iconset.path, "-o", destination.appendingPathComponent("AppIcon.icns").path]
  try process.run()
  process.waitUntilExit()
  guard process.terminationStatus == 0 else { fatalError("iconutil failed for \(name)") }
  print("Created \(name).icns")
}
