import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:uuid/uuid.dart';
import '../../../core/network/services/chat_api_service.dart';
import '../../../core/network/websocket_client.dart';
import '../../../core/storage/storage_service.dart';
import '../../../core/config/app_config.dart';
import '../../../shared/utils/logger.dart';
import '../../../shared/models/chat_message.dart';
import '../../../shared/models/product.dart';
import '../../auth/providers/auth_provider.dart';
import 'chat_state.dart';

/// Chat provider - manages chat state and WebSocket connection
final chatProvider = StateNotifierProvider<ChatNotifier, ChatState>((ref) {
  final chatApiService = ref.watch(chatApiServiceProvider);
  final webSocketClient = ref.watch(webSocketClientProvider);
  final user = ref.watch(currentUserProvider);

  return ChatNotifier(
    chatApiService: chatApiService,
    webSocketClient: webSocketClient,
    userCountry: user?.country,
    userLanguage: user?.language,
    userCurrency: user?.currency,
  );
});

/// Chat notifier - handles chat logic
class ChatNotifier extends StateNotifier<ChatState> {
  final ChatApiService _chatApiService;
  final WebSocketClient _webSocketClient;
  final String? _userCountry;
  final String? _userLanguage;
  final String? _userCurrency;

  StreamSubscription? _messageSubscription;
  StreamSubscription? _connectionSubscription;
  final _uuid = const Uuid();

  ChatNotifier({
    required ChatApiService chatApiService,
    required WebSocketClient webSocketClient,
    String? userCountry,
    String? userLanguage,
    String? userCurrency,
  })  : _chatApiService = chatApiService,
        _webSocketClient = webSocketClient,
        _userCountry = userCountry,
        _userLanguage = userLanguage,
        _userCurrency = userCurrency,
        super(ChatState.initial()) {
    _initialize();
  }

  /// Initialize chat
  Future<void> _initialize() async {
    try {
      // Restore session from storage
      await _restoreSession();

      // Listen to WebSocket messages
      _messageSubscription = _webSocketClient.messages.listen(_handleWebSocketMessage);

      // Listen to connection state changes
      _connectionSubscription = _webSocketClient.connectionState.listen(_handleConnectionStateChange);

      // Connect to WebSocket
      await connectWebSocket();
    } catch (e, stackTrace) {
      AppLogger.error('Failed to initialize chat', e, stackTrace);
    }
  }

  /// Restore session from storage
  Future<void> _restoreSession() async {
    try {
      state = state.copyWith(isLoading: true);

      final prefs = StorageService.prefs;
      final sessionId = prefs.getString(AppConfig.sessionIdKey);

      if (sessionId != null) {
        // Restore messages from Hive
        final messagesBox = StorageService.getBox(AppConfig.messagesBox);
        final storedMessages = messagesBox.get(sessionId);

        if (storedMessages is List) {
          final messages = storedMessages
              .map((msg) => ChatMessage.fromJson(Map<String, dynamic>.from(msg as Map)))
              .toList();

          state = state.copyWith(
            sessionId: sessionId,
            messages: messages,
          );

          AppLogger.info('✅ Session restored: $sessionId with ${messages.length} messages');
        }
      }
    } catch (e, stackTrace) {
      AppLogger.error('Failed to restore session', e, stackTrace);
    } finally {
      state = state.copyWith(isLoading: false);
    }
  }

  /// Connect to WebSocket
  Future<void> connectWebSocket() async {
    try {
      await _webSocketClient.connect();
    } catch (e, stackTrace) {
      AppLogger.error('Failed to connect WebSocket', e, stackTrace);
      state = state.copyWith(error: 'Failed to connect to chat server');
    }
  }

  /// Disconnect from WebSocket
  Future<void> disconnectWebSocket() async {
    try {
      await _webSocketClient.disconnect();
    } catch (e, stackTrace) {
      AppLogger.error('Failed to disconnect WebSocket', e, stackTrace);
    }
  }

  /// Handle WebSocket message
  void _handleWebSocketMessage(Map<String, dynamic> data) {
    try {
      AppLogger.debug('Received WebSocket message: $data');

      // Handle different message types
      final type = data['type'] as String?;

      switch (type) {
        case 'message':
          _handleChatMessage(data);
          break;
        case 'typing':
          state = state.copyWith(isTyping: data['value'] == true);
          break;
        case 'error':
          state = state.copyWith(error: data['message'] as String?);
          break;
        default:
          AppLogger.warning('Unknown message type: $type');
      }
    } catch (e, stackTrace) {
      AppLogger.error('Failed to handle WebSocket message', e, stackTrace);
    }
  }

  /// Handle chat message from WebSocket
  void _handleChatMessage(Map<String, dynamic> data) {
    try {
      final message = ChatMessage(
        id: _uuid.v4(),
        role: MessageRole.assistant,
        content: data['message'] as String,
        timestamp: DateTime.now().millisecondsSinceEpoch,
        quickReplies: (data['quick_replies'] as List<dynamic>?)?.cast<String>(),
        products: (data['products'] as List<dynamic>?)
            ?.map((p) => Product.fromJson(Map<String, dynamic>.from(p as Map)))
            .toList(),
        searchType: data['search_type'] as String?,
      );

      _addMessage(message);

      // Update quick replies
      if (message.quickReplies != null) {
        state = state.copyWith(
          quickReplies: message.quickReplies!,
          isTyping: false,
        );
      }
    } catch (e, stackTrace) {
      AppLogger.error('Failed to handle chat message', e, stackTrace);
    }
  }

  /// Handle connection state change
  void _handleConnectionStateChange(WebSocketState newState) {
    state = state.copyWith(wsState: newState);
    AppLogger.info('WebSocket state changed: $newState');
  }

  /// Send message via WebSocket
  Future<void> sendMessageViaWebSocket(String message) async {
    if (!state.isConnected) {
      // Fallback to REST API if WebSocket is not connected
      await sendMessageViaRest(message);
      return;
    }

    try {
      state = state.copyWith(isSending: true, error: null);

      // Create user message
      final userMessage = ChatMessage(
        id: _uuid.v4(),
        role: MessageRole.user,
        content: message,
        timestamp: DateTime.now().millisecondsSinceEpoch,
        isLocal: true,
      );

      _addMessage(userMessage);

      // Send via WebSocket
      _webSocketClient.send({
        'type': 'message',
        'message': message,
        'session_id': state.sessionId,
        'country': _userCountry,
        'language': _userLanguage,
        'currency': _userCurrency,
      });

      // Clear quick replies
      state = state.copyWith(quickReplies: [], isTyping: true);
    } catch (e, stackTrace) {
      AppLogger.error('Failed to send message via WebSocket', e, stackTrace);
      state = state.copyWith(error: 'Failed to send message');
    } finally {
      state = state.copyWith(isSending: false);
    }
  }

  /// Send message via REST API (fallback)
  Future<void> sendMessageViaRest(String message) async {
    try {
      state = state.copyWith(isSending: true, error: null);

      // Create user message
      final userMessage = ChatMessage(
        id: _uuid.v4(),
        role: MessageRole.user,
        content: message,
        timestamp: DateTime.now().millisecondsSinceEpoch,
        isLocal: true,
      );

      _addMessage(userMessage);

      // Send via REST API
      final response = await _chatApiService.sendMessage(
        message: message,
        sessionId: state.sessionId,
        country: _userCountry,
        language: _userLanguage,
        currency: _userCurrency,
      );

      // Update session ID
      if (state.sessionId != response.sessionId) {
        _updateSessionId(response.sessionId);
      }

      // Create assistant message
      final assistantMessage = ChatMessage(
        id: _uuid.v4(),
        role: MessageRole.assistant,
        content: response.message,
        timestamp: DateTime.now().millisecondsSinceEpoch,
        quickReplies: response.quickReplies,
        products: response.products,
        searchType: response.searchType,
      );

      _addMessage(assistantMessage);

      // Update quick replies
      state = state.copyWith(
        quickReplies: response.quickReplies ?? [],
      );
    } catch (e, stackTrace) {
      AppLogger.error('Failed to send message via REST', e, stackTrace);
      state = state.copyWith(error: 'Failed to send message');
    } finally {
      state = state.copyWith(isSending: false);
    }
  }

  /// Send quick reply
  Future<void> sendQuickReply(String reply) async {
    if (state.isConnected) {
      await sendMessageViaWebSocket(reply);
    } else {
      await sendMessageViaRest(reply);
    }
  }

  /// Add message to state and storage
  void _addMessage(ChatMessage message) {
    final updatedMessages = [...state.messages, message];
    state = state.copyWith(messages: updatedMessages);

    // Save to storage
    _saveMessagesToStorage();
  }

  /// Update session ID
  void _updateSessionId(String sessionId) {
    state = state.copyWith(sessionId: sessionId);

    // Save to SharedPreferences
    final prefs = StorageService.prefs;
    prefs.setString(AppConfig.sessionIdKey, sessionId);

    AppLogger.info('Session ID updated: $sessionId');
  }

  /// Save messages to storage
  Future<void> _saveMessagesToStorage() async {
    try {
      if (state.sessionId == null) return;

      final messagesBox = StorageService.getBox(AppConfig.messagesBox);
      final messagesJson = state.messages.map((m) => m.toJson()).toList();

      await messagesBox.put(state.sessionId, messagesJson);
    } catch (e, stackTrace) {
      AppLogger.error('Failed to save messages to storage', e, stackTrace);
    }
  }

  /// Clear chat
  Future<void> clearChat() async {
    try {
      // Clear from storage
      if (state.sessionId != null) {
        final messagesBox = StorageService.getBox(AppConfig.messagesBox);
        await messagesBox.delete(state.sessionId);
      }

      // Clear from SharedPreferences
      final prefs = StorageService.prefs;
      await prefs.remove(AppConfig.sessionIdKey);

      // Reset state
      state = ChatState.initial();

      AppLogger.info('✅ Chat cleared');
    } catch (e, stackTrace) {
      AppLogger.error('Failed to clear chat', e, stackTrace);
    }
  }

  /// Clear error
  void clearError() {
    state = state.copyWith(error: null);
  }

  @override
  void dispose() {
    _messageSubscription?.cancel();
    _connectionSubscription?.cancel();
    super.dispose();
  }
}

/// Helper providers
final chatMessagesProvider = Provider<List<ChatMessage>>((ref) {
  return ref.watch(chatProvider).messages;
});

final chatQuickRepliesProvider = Provider<List<String>>((ref) {
  return ref.watch(chatProvider).quickReplies;
});

final chatIsTypingProvider = Provider<bool>((ref) {
  return ref.watch(chatProvider).isTyping;
});

final chatIsConnectedProvider = Provider<bool>((ref) {
  return ref.watch(chatProvider).isConnected;
});

final chatIsSendingProvider = Provider<bool>((ref) {
  return ref.watch(chatProvider).isSending;
});
