import Foundation
import HISAPI
import Observation

@MainActor @Observable
final class UserSession {
  let gateway: Gateway
  var profile: AppProfile?

  init(address: String) {
    let secureAddress = Gateway(address: address).address
    gateway = Gateway(address: secureAddress, tokens: SessionTokens(address: secureAddress))
  }

  func enter(_ session: Session) async throws {
    guard AppConfiguration.accepts(role: session.role) else {
      _ = try? await gateway.security.logout(tokenInput: TokenInput(refreshToken: session.refreshToken))
      throw SessionRoleError()
    }
    await gateway.tokens?.set(session)
    do {
      if session.role == "admin" {
        profile = .administrator
      } else if let staffID = session.staffId {
        var worker = try await gateway.staff.getStaff(id: staffID)
        worker.role = StaffRole(rawValue: session.role) ?? .other
        profile = .worker(worker)
      } else { throw URLError(.userAuthenticationRequired) }
    } catch {
      await gateway.tokens?.set(nil)
      throw error
    }
  }

  func check() async {
    do {
      guard let current = try await gateway.tokens?.current(),
        AppConfiguration.accepts(role: current.role) else {
        await logout()
        return
      }
      if case .worker(var worker) = profile {
        worker.role = StaffRole(rawValue: current.role) ?? .other
        profile = .worker(worker)
      }
    } catch { profile = nil }
  }

  func logout() async {
    if let current = try? await gateway.tokens?.current() {
      _ = try? await gateway.security.logout(tokenInput: TokenInput(refreshToken: current.refreshToken))
    }
    await gateway.tokens?.set(nil)
    profile = nil
  }
}

private struct SessionRoleError: LocalizedError {
  var errorDescription: String? {
    AppConfiguration.isAdministrator
      ? "Для этого приложения требуется учётная запись администратора."
      : "Для учётной записи администратора откройте приложение HISAdmin."
  }
}
