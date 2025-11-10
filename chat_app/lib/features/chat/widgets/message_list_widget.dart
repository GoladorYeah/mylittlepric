import 'package:flutter/material.dart';
import 'package:chat_app/shared/models/chat_message.dart';
import 'package:chat_app/features/chat/widgets/message_bubble.dart';
import 'package:chat_app/features/chat/widgets/chat_empty_state.dart';
import 'package:chat_app/features/chat/widgets/typing_indicator.dart';

/// Widget displaying the list of chat messages with auto-scroll
class MessageListWidget extends StatefulWidget {
  final List<ChatMessage> messages;
  final bool isTyping;
  final ValueChanged<String> onQuickReply;

  const MessageListWidget({
    super.key,
    required this.messages,
    required this.isTyping,
    required this.onQuickReply,
  });

  @override
  State<MessageListWidget> createState() => _MessageListWidgetState();
}

class _MessageListWidgetState extends State<MessageListWidget> {
  final _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _scrollToBottom();
    });
  }

  @override
  void didUpdateWidget(MessageListWidget oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.messages.length != oldWidget.messages.length ||
        widget.isTyping != oldWidget.isTyping) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _scrollToBottom();
      });
    }
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _scrollToBottom() {
    if (_scrollController.hasClients) {
      _scrollController.animateTo(
        _scrollController.position.maxScrollExtent,
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeOut,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    if (widget.messages.isEmpty && !widget.isTyping) {
      return const ChatEmptyState();
    }

    return ListView.separated(
      controller: _scrollController,
      padding: const EdgeInsets.all(16),
      itemCount: widget.messages.length + (widget.isTyping ? 1 : 0),
      separatorBuilder: (context, index) => const SizedBox(height: 16),
      itemBuilder: (context, index) {
        // Show typing indicator at the end
        if (index == widget.messages.length) {
          return const TypingIndicator();
        }

        final message = widget.messages[index];
        return MessageBubble(
          key: ValueKey(message.id),
          message: message,
          onQuickReply: widget.onQuickReply,
        );
      },
    );
  }
}
