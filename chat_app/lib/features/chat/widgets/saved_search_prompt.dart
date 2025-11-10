import 'package:flutter/material.dart';
import 'package:chat_app/shared/models/saved_search.dart';
import 'package:chat_app/shared/models/chat_message.dart';

/// Dialog prompting user to continue or start new search
/// Displayed when there's an incomplete saved search
class SavedSearchPrompt extends StatelessWidget {
  final SavedSearch? savedSearch;
  final VoidCallback onContinue;
  final VoidCallback onNewSearch;

  const SavedSearchPrompt({
    super.key,
    required this.savedSearch,
    required this.onContinue,
    required this.onNewSearch,
  });

  String _getPreviewText() {
    if (savedSearch == null || savedSearch!.messages.isEmpty) {
      return '';
    }

    // Get last user message
    final userMessages = savedSearch!.messages
        .where((m) => m.role == MessageRole.user)
        .toList();

    if (userMessages.isEmpty) return '';

    final lastMessage = userMessages.last.content;
    const maxLength = 60;

    if (lastMessage.length > maxLength) {
      return '${lastMessage.substring(0, maxLength)}...';
    }

    return lastMessage;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final preview = _getPreviewText();

    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Card(
          elevation: 4,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: Padding(
            padding: const EdgeInsets.all(32.0),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Icon
                Container(
                  width: 64,
                  height: 64,
                  decoration: BoxDecoration(
                    color: theme.colorScheme.primary.withOpacity(0.1),
                    shape: BoxShape.circle,
                  ),
                  child: Icon(
                    Icons.refresh,
                    size: 32,
                    color: theme.colorScheme.primary,
                  ),
                ),
                const SizedBox(height: 24),

                // Title
                Text(
                  'У вас есть незавершенный поиск',
                  style: theme.textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 8),

                // Preview
                if (preview.isNotEmpty) ...[
                  Text(
                    '"$preview"',
                    style: theme.textTheme.bodySmall?.copyWith(
                      fontStyle: FontStyle.italic,
                      color: theme.colorScheme.onSurface.withOpacity(0.6),
                    ),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 24),
                ],

                // Buttons
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: onNewSearch,
                        style: OutlinedButton.styleFrom(
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Icon(Icons.add, size: 20),
                            const SizedBox(width: 8),
                            Text(
                              'Начать новый',
                              style: theme.textTheme.labelLarge,
                            ),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: onContinue,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: theme.colorScheme.primary,
                          foregroundColor: theme.colorScheme.onPrimary,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Icon(Icons.refresh, size: 20),
                            const SizedBox(width: 8),
                            Text(
                              'Продолжить',
                              style: theme.textTheme.labelLarge?.copyWith(
                                color: theme.colorScheme.onPrimary,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                // Footer text
                Text(
                  'Вы можете продолжить незавершенный поиск или начать новый',
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: theme.colorScheme.onSurface.withOpacity(0.6),
                  ),
                  textAlign: TextAlign.center,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
