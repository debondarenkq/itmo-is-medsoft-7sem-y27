// swift-tools-version: 6.0
import PackageDescription

let package = Package(
  name: "HIS",
  platforms: [.macOS(.v14)],
  products: [.executable(name: "HIS", targets: ["HIS"])],
  dependencies: [.package(path: "Generated")],
  targets: [
    .executableTarget(name: "HIS", dependencies: [.product(name: "HISAPI", package: "Generated")])
  ]
)
