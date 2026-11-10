import Foundation
import HISAPI

// Tokens live only in memory. A shared task prevents concurrent refresh requests.
actor SessionTokens {
  private let address: String
  private var session: Session?
  private var refreshing: Task<Session, Error>?
  private var generation = 0

  init(address: String) { self.address = address }

  func set(_ value: Session?) {
    generation += 1
    refreshing?.cancel()
    refreshing = nil
    session = value
  }

  func current() async throws -> Session {
    guard let session else { throw URLError(.userAuthenticationRequired) }
    if session.expiresAt.timeIntervalSinceNow > 15 { return session }
    if let refreshing { return try await refreshing.value }
    let api = Gateway(address: address).security
    let startedGeneration = generation
    let task = Task { try await api.refreshSession(tokenInput: TokenInput(refreshToken: session.refreshToken)) }
    refreshing = task
    defer { refreshing = nil }
    do {
      let updated = try await task.value
      guard generation == startedGeneration else { throw CancellationError() }
      self.session = updated
      return updated
    } catch {
      if generation == startedGeneration { self.session = nil }
      throw error
    }
  }
}

struct SessionInterceptor: OpenAPIInterceptor {
  let tokens: SessionTokens

  func intercept<T>(urlRequest: URLRequest, urlSession: URLSessionProtocol,
    requestBuilder: RequestBuilder<T>,
    completion: @Sendable @escaping (Result<URLRequest, Error>) -> Void) {
    Task {
      do {
        let session = try await tokens.current()
        var request = urlRequest
        request.setValue("Bearer \(session.accessToken)", forHTTPHeaderField: "Authorization")
        completion(.success(request))
      } catch { completion(.failure(error)) }
    }
  }

  func retry<T>(urlRequest: URLRequest, urlSession: URLSessionProtocol,
    requestBuilder: RequestBuilder<T>, data: Data?, response: URLResponse?, error: Error,
    completion: @Sendable @escaping (OpenAPIInterceptorRetry) -> Void) {
    completion(.dontRetry)
  }
}
