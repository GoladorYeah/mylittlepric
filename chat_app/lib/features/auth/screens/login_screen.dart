import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../theme/app_colors.dart';
import '../../../shared/widgets/widgets.dart';
import '../providers/auth_provider.dart';
import '../widgets/widgets.dart';

/// Login screen with OAuth providers
class LoginScreen extends ConsumerWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);

    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                // Logo and title
                Icon(
                  Icons.shopping_bag,
                  size: 80,
                  color: AppColors.primary,
                ),
                const SizedBox(height: 24),
                Text(
                  'MyLittlePrice',
                  style: Theme.of(context).textTheme.headlineLarge?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: AppColors.primary,
                      ),
                ),
                const SizedBox(height: 12),
                Text(
                  'AI-powered product search assistant',
                  style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                        color: Theme.of(context)
                            .colorScheme
                            .onSurface
                            .withOpacity(0.6),
                      ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 48),
                // Loading indicator
                if (authState.isLoading)
                  const Column(
                    children: [
                      CircularProgressIndicator(),
                      SizedBox(height: 24),
                      Text('Signing in...'),
                    ],
                  )
                else ...[
                  // OAuth buttons
                  const Text(
                    'Sign in with:',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 24),
                  OAuthButton(
                    provider: 'Google',
                    icon: Icons.login,
                    backgroundColor: Colors.white,
                    textColor: Colors.black87,
                    onPressed: () => _handleGoogleLogin(context, ref),
                  ),
                  const SizedBox(height: 16),
                  OAuthButton(
                    provider: 'Facebook',
                    icon: Icons.facebook,
                    backgroundColor: const Color(0xFF1877F2),
                    textColor: Colors.white,
                    onPressed: () => _handleFacebookLogin(context, ref),
                  ),
                  const SizedBox(height: 16),
                  OAuthButton(
                    provider: 'Apple',
                    icon: Icons.apple,
                    backgroundColor: Colors.black,
                    textColor: Colors.white,
                    onPressed: () => _handleAppleLogin(context, ref),
                  ),
                  const SizedBox(height: 32),
                  // Continue as guest
                  OutlinedButton(
                    onPressed: () => _continueAsGuest(context),
                    style: OutlinedButton.styleFrom(
                      minimumSize: const Size(double.infinity, 56),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: const Text('Continue as Guest'),
                  ),
                ],
                const SizedBox(height: 32),
                // Error message
                if (authState.error != null)
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: AppColors.error.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(
                        color: AppColors.error,
                        width: 1,
                      ),
                    ),
                    child: Row(
                      children: [
                        const Icon(
                          Icons.error_outline,
                          color: AppColors.error,
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Text(
                            authState.error!,
                            style: const TextStyle(
                              color: AppColors.error,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                const SizedBox(height: 32),
                // Terms and privacy
                Text(
                  'By continuing, you agree to our',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: Theme.of(context)
                            .colorScheme
                            .onSurface
                            .withOpacity(0.5),
                      ),
                ),
                const SizedBox(height: 4),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    TextButton(
                      onPressed: () {
                        // TODO: Open terms
                      },
                      child: const Text('Terms of Service'),
                    ),
                    Text(
                      ' and ',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: Theme.of(context)
                                .colorScheme
                                .onSurface
                                .withOpacity(0.5),
                          ),
                    ),
                    TextButton(
                      onPressed: () {
                        // TODO: Open privacy policy
                      },
                      child: const Text('Privacy Policy'),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _handleGoogleLogin(BuildContext context, WidgetRef ref) async {
    try {
      // TODO: Implement Google OAuth flow
      // 1. Open Google OAuth URL
      // 2. Get authorization code
      // 3. Exchange for token
      // 4. Call login with token

      // For now, show a message
      if (context.mounted) {
        ErrorSnackBar.show(
          context,
          'Google OAuth integration coming soon!',
        );
      }
    } catch (e) {
      if (context.mounted) {
        ErrorSnackBar.show(context, e.toString());
      }
    }
  }

  Future<void> _handleFacebookLogin(BuildContext context, WidgetRef ref) async {
    try {
      // TODO: Implement Facebook OAuth flow

      if (context.mounted) {
        ErrorSnackBar.show(
          context,
          'Facebook OAuth integration coming soon!',
        );
      }
    } catch (e) {
      if (context.mounted) {
        ErrorSnackBar.show(context, e.toString());
      }
    }
  }

  Future<void> _handleAppleLogin(BuildContext context, WidgetRef ref) async {
    try {
      // TODO: Implement Apple OAuth flow

      if (context.mounted) {
        ErrorSnackBar.show(
          context,
          'Apple OAuth integration coming soon!',
        );
      }
    } catch (e) {
      if (context.mounted) {
        ErrorSnackBar.show(context, e.toString());
      }
    }
  }

  void _continueAsGuest(BuildContext context) {
    // Navigate to chat without authentication
    context.go('/chat');
  }
}
