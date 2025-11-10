import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../shared/models/user.dart';

part 'auth_state.freezed.dart';

/// Authentication state
@freezed
class AuthState with _$AuthState {
  const factory AuthState({
    /// Current authenticated user
    User? user,

    /// Access token
    String? accessToken,

    /// Refresh token
    String? refreshToken,

    /// Loading state
    @Default(false) bool isLoading,

    /// Error message
    String? error,

    /// Whether the user is authenticated
    @Default(false) bool isAuthenticated,
  }) = _AuthState;

  const AuthState._();

  /// Initial state
  factory AuthState.initial() => const AuthState();

  /// Loading state
  factory AuthState.loading() => const AuthState(isLoading: true);

  /// Authenticated state
  factory AuthState.authenticated({
    required User user,
    required String accessToken,
    required String refreshToken,
  }) =>
      AuthState(
        user: user,
        accessToken: accessToken,
        refreshToken: refreshToken,
        isAuthenticated: true,
      );

  /// Error state
  factory AuthState.error(String message) => AuthState(
        error: message,
        isAuthenticated: false,
      );
}
