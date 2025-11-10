import 'package:chat_app/core/config/constants.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:chat_app/features/chat/widgets/widgets.dart';
import 'package:chat_app/features/history/screens/screens.dart';
import 'package:chat_app/features/settings/screens/screens.dart';
import 'package:chat_app/features/auth/screens/screens.dart';

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

          return ChatScreen(
            initialQuery: query,
            sessionId: sessionId,
          );
        },
      ),
      GoRoute(
        path: AppConstants.historyRoute,
        name: 'history',
        builder: (context, state) {
          return const HistoryScreen();
        },
      ),
      GoRoute(
        path: AppConstants.settingsRoute,
        name: 'settings',
        builder: (context, state) {
          return const SettingsScreen();
        },
      ),
      GoRoute(
        path: AppConstants.loginRoute,
        name: 'login',
        builder: (context, state) {
          return const LoginScreen();
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
