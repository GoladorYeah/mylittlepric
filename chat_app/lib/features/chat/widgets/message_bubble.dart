import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:chat_app/shared/models/chat_message.dart';
import 'package:chat_app/features/auth/providers/auth_provider.dart';
import 'package:chat_app/features/chat/widgets/quick_replies_widget.dart';
import 'package:chat_app/features/products/widgets/product_card_widget.dart';

/// Message bubble displaying a single chat message
/// Supports user/assistant messages, quick replies, and product cards
class MessageBubble extends ConsumerWidget {
  final ChatMessage message;
  final ValueChanged<String> onQuickReply;

  const MessageBubble({
    super.key,
    required this.message,
    required this.onQuickReply,
  });

  /// Generate initials from user's name or email
  String _getInitials(String? fullName, String? email) {
    if (fullName != null && fullName.isNotEmpty) {
      final names = fullName.trim().split(RegExp(r'\s+'));
      if (names.length >= 2) {
        return '${names.first[0]}${names.last[0]}'.toUpperCase();
      }
      return names.first[0].toUpperCase();
    }

    if (email != null && email.isNotEmpty) {
      return email[0].toUpperCase();
    }

    return 'U';
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final isUser = message.role == MessageRole.user;
    final user = ref.watch(currentUserProvider);

    return Align(
      alignment: isUser ? Alignment.centerRight : Alignment.centerLeft,
      child: Container(
        constraints: BoxConstraints(
          maxWidth: message.products != null && message.products!.isNotEmpty
              ? double.infinity
              : MediaQuery.of(context).size.width * 0.8,
        ),
        child: Column(
          crossAxisAlignment:
              isUser ? CrossAxisAlignment.end : CrossAxisAlignment.start,
          children: [
            // Message content
            if (message.content.isNotEmpty)
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 12,
                ),
                decoration: BoxDecoration(
                  color: isUser
                      ? theme.colorScheme.secondaryContainer
                      : Colors.transparent,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // User avatar (for user messages)
                    if (isUser) ...[
                      CircleAvatar(
                        radius: 18,
                        backgroundColor: theme.colorScheme.primary,
                        child: Text(
                          _getInitials(user?.name, user?.email),
                          style: theme.textTheme.labelMedium?.copyWith(
                            color: theme.colorScheme.onPrimary,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                    ],

                    // Message text
                    Flexible(
                      child: Text(
                        message.content,
                        style: theme.textTheme.bodyMedium?.copyWith(
                          color: isUser
                              ? theme.colorScheme.onSecondaryContainer
                              : theme.colorScheme.onBackground,
                        ),
                      ),
                    ),
                  ],
                ),
              ),

            // Quick replies
            if (message.quickReplies != null &&
                message.quickReplies!.isNotEmpty) ...[
              const SizedBox(height: 8),
              QuickRepliesWidget(
                quickReplies: message.quickReplies!,
                onReplyTap: onQuickReply,
              ),
            ],

            // Products
            if (message.products != null && message.products!.isNotEmpty) ...[
              const SizedBox(height: 12),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Header with count and navigation
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 8),
                    child: Row(
                      children: [
                        Container(
                          width: 4,
                          height: 32,
                          decoration: BoxDecoration(
                            gradient: LinearGradient(
                              colors: [
                                theme.colorScheme.primary,
                                theme.colorScheme.primary.withOpacity(0.5),
                              ],
                              begin: Alignment.topCenter,
                              end: Alignment.bottomCenter,
                            ),
                            borderRadius: BorderRadius.circular(2),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              '${message.products!.length} ${message.products!.length == 1 ? 'Product' : 'Products'} Found',
                              style: theme.textTheme.titleMedium?.copyWith(
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            Text(
                              'Swipe to explore more options',
                              style: theme.textTheme.labelSmall?.copyWith(
                                color:
                                    theme.colorScheme.onBackground.withOpacity(0.6),
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Products carousel
                  SizedBox(
                    height: 340,
                    child: ListView.separated(
                      scrollDirection: Axis.horizontal,
                      padding: const EdgeInsets.symmetric(horizontal: 8),
                      itemCount: message.products!.length,
                      separatorBuilder: (context, index) =>
                          const SizedBox(width: 12),
                      itemBuilder: (context, index) {
                        return ProductCardWidget(
                          product: message.products![index],
                          position: index + 1,
                        );
                      },
                    ),
                  ),

                  // Scroll indicator
                  if (message.products!.length > 1) ...[
                    const SizedBox(height: 12),
                    Center(
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: List.generate(
                          message.products!.length.clamp(0, 10),
                          (index) {
                            final isFirst = index < 3;
                            return Container(
                              margin: const EdgeInsets.symmetric(horizontal: 3),
                              width: isFirst ? 24 : 8,
                              height: 4,
                              decoration: BoxDecoration(
                                color: theme.colorScheme.primary.withOpacity(
                                  isFirst ? 1.0 : 0.3,
                                ),
                                borderRadius: BorderRadius.circular(2),
                              ),
                            );
                          },
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }
}
