import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:chat_app/features/chat/providers/chat_provider.dart';
import 'package:chat_app/features/chat/widgets/message_list_widget.dart';
import 'package:chat_app/features/chat/widgets/chat_input_widget.dart';
import 'package:chat_app/features/chat/widgets/saved_search_prompt.dart';

/// Main chat screen with messages, input, and saved search prompt
class ChatScreen extends ConsumerStatefulWidget {
  final String? initialQuery;
  final String? sessionId;

  const ChatScreen({
    super.key,
    this.initialQuery,
    this.sessionId,
  });

  @override
  ConsumerState<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends ConsumerState<ChatScreen> {
  bool _initialQuerySent = false;

  @override
  void initState() {
    super.initState();
    _initializeChat();
  }

  Future<void> _initializeChat() async {
    final chatNotifier = ref.read(chatProvider.notifier);

    // Initialize WebSocket connection
    await chatNotifier.connect();

    // Send initial query if provided
    if (widget.initialQuery != null &&
        widget.initialQuery!.isNotEmpty &&
        !_initialQuerySent) {
      _initialQuerySent = true;
      WidgetsBinding.instance.addPostFrameCallback((_) {
        chatNotifier.sendMessage(widget.initialQuery!);
      });
    }
  }

  @override
  void dispose() {
    // Disconnect WebSocket when screen is disposed
    ref.read(chatProvider.notifier).disconnect();
    super.dispose();
  }

  void _handleSendMessage(String message) {
    ref.read(chatProvider.notifier).sendMessage(message);
  }

  void _handleQuickReply(String reply) {
    ref.read(chatProvider.notifier).sendMessage(reply);
  }

  void _handleContinueSearch() {
    // TODO: Implement restore saved search
    // For now, just hide the prompt
    setState(() {
      // Hide saved search prompt
    });
  }

  void _handleNewSearch() {
    ref.read(chatProvider.notifier).clearMessages();
    setState(() {
      // Hide saved search prompt
    });
  }

  @override
  Widget build(BuildContext context) {
    final chatState = ref.watch(chatProvider);
    final messages = ref.watch(chatMessagesProvider);
    final isTyping = ref.watch(chatIsTypingProvider);
    final isConnected = ref.watch(chatIsConnectedProvider);
    final isSending = ref.watch(chatIsSendingProvider);

    // Determine connection status text
    String connectionStatus = 'Disconnected';
    if (chatState.connectionState == ConnectionState.connecting) {
      connectionStatus = 'Connecting';
    } else if (chatState.connectionState == ConnectionState.connected) {
      connectionStatus = 'Connected';
    }

    // Check if we should show saved search prompt
    // TODO: Get this from chat state or provider
    final showSavedSearchPrompt = false;
    final savedSearch = null;

    return Scaffold(
      appBar: AppBar(
        title: const Text('MyLittlePrice'),
        actions: [
          // Connection status indicator
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: Center(
              child: Row(
                children: [
                  Icon(
                    isConnected ? Icons.wifi : Icons.wifi_off,
                    size: 16,
                    color: isConnected ? Colors.green : Colors.red,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    connectionStatus,
                    style: Theme.of(context).textTheme.labelSmall?.copyWith(
                          color: isConnected ? Colors.green : Colors.red,
                          fontWeight: FontWeight.w600,
                        ),
                  ),
                ],
              ),
            ),
          ),

          // New search button
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: isConnected ? _handleNewSearch : null,
            tooltip: 'New Search',
          ),

          // Settings button
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () {
              // TODO: Navigate to settings
            },
            tooltip: 'Settings',
          ),
        ],
      ),
      body: Column(
        children: [
          // Main content area
          Expanded(
            child: showSavedSearchPrompt
                ? SavedSearchPrompt(
                    savedSearch: savedSearch,
                    onContinue: _handleContinueSearch,
                    onNewSearch: _handleNewSearch,
                  )
                : MessageListWidget(
                    messages: messages,
                    isTyping: isTyping,
                    onQuickReply: _handleQuickReply,
                  ),
          ),

          // Input field (hidden when showing saved search prompt)
          if (!showSavedSearchPrompt)
            ChatInputWidget(
              onSend: _handleSendMessage,
              isLoading: isSending,
              isConnected: isConnected,
              connectionStatus: connectionStatus,
            ),
        ],
      ),
    );
  }
}
