import AppKit
import CoreImage.CIFilterBuiltins

func totpQRCode(secret: String) -> NSImage? {
  var url = URLComponents()
  url.scheme = "otpauth"
  url.host = "totp"
  url.path = "/HIS:admin"
  url.queryItems = [
    URLQueryItem(name: "secret", value: secret), URLQueryItem(name: "issuer", value: "HIS"),
    URLQueryItem(name: "algorithm", value: "SHA1"), URLQueryItem(name: "digits", value: "6"),
    URLQueryItem(name: "period", value: "30")
  ]
  guard let text = url.string else { return nil }
  let filter = CIFilter.qrCodeGenerator()
  filter.message = Data(text.utf8)
  guard let output = filter.outputImage?.transformed(by: CGAffineTransform(scaleX: 8, y: 8)),
    let image = CIContext().createCGImage(output, from: output.extent) else { return nil }
  return NSImage(cgImage: image, size: NSSize(width: image.width, height: image.height))
}
