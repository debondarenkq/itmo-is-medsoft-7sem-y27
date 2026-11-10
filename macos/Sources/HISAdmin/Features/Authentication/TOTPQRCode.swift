import AppKit
import CoreImage.CIFilterBuiltins

func totpQRCode(provisioningURI: String) -> NSImage? {
  let filter = CIFilter.qrCodeGenerator()
  filter.message = Data(provisioningURI.utf8)
  guard let output = filter.outputImage?.transformed(by: CGAffineTransform(scaleX: 8, y: 8)),
    let image = CIContext().createCGImage(output, from: output.extent) else { return nil }
  return NSImage(cgImage: image, size: NSSize(width: image.width, height: image.height))
}
