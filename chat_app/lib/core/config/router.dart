import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:chat_app/core/config/constants.dart';

// Import screens (will be created later)
// import 'package:chat_app/features/chat/screens/chat_screen.dart';
// import 'package:chat_app/features/history/screens/history_screen.dart';
// import 'package:chat_app/features/settings/screens/settings_screen.dart';
// import 'package:chat_app/features/auth/screens/login_screen.dart';

class AppRouter {
  static final GoRouter router = GoRouter(
    initialLocation: AppConstants.chatRoute,
    debugLogDiagnostics: true,
    routes: [
      GoRoute(
        path: AppConstants.chatRoute,
        name: 'chat',
        builder: (context, state) {
          // Extract query parameters
          final query = state.uri.queryParameters['q'];
          final sessionId = state.uri.queryParameters['session_id'];

          // TODO: Pass to ChatScreen when implemented
          return const Scaffold(
            body: Center(
              child: Text('Chat Screen - Coming Soon'),
            ),
          );
          // return ChatScreen(
          //   initialQuery: query,
          //   sessionId: sessionId,
          // );
        },
      ),
      GoRoute(
        path: AppConstants.historyRoute,
        name: 'history',
        builder: (context, state) {
          return const Scaffold(
            body: Center(
              child: Text('History Screen - Coming Soon'),
            ),
          );
          // return const HistoryScreen();
        },
      ),
      GoRoute(
        path: AppConstants.settingsRoute,
        name: 'settings',
        builder: (context, state) {
          return const Scaffold(
            body: Center(
              child: Text('Settings Screen - Coming Soon'),
            ),
          );
          // return const SettingsScreen();
        },
      ),
      GoRoute(
        path: AppConstants.loginRoute,
        name: 'login',
        builder: (context, state) {
          return const Scaffold(
            body: Center(
              child: Text('Login Screen - Coming Soon'),
            ),
          );
          // return const LoginScreen();
        },
      ),
    ],
    errorBuilder: (context, state) {
      return Scaffold(
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 64, color: Colors.red),
              const SizedBox(height: 16),
              Text(
                'Page not found',
                style: Theme.of(context).textTheme.headlineMedium,
              ),
              const SizedBox(height: 8),
              Text(
                state.uri.toString(),
                style: Theme.of(context).textTheme.bodySmall,
              ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: () => context.go(AppConstants.chatRoute),
                child: const Text('Go to Chat'),
              ),
            ],
          ),
        ),
      );
    },
  );
}
