import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../shared/models/chat_message.dart';
import '../../../core/network/websocket_client.dart';

part 'chat_state.freezed.dart';

/// Chat state
@freezed
class ChatState with _$ChatState {
  const factory ChatState({
    /// Current session ID
    String? sessionId,

    /// List of chat messages
    @Default([]) List<ChatMessage> messages,

    /// WebSocket connection state
    @Default(WebSocketState.disconnected) WebSocketState wsState,

    /// Whether a message is being sent
    @Default(false) bool isSending,

    /// Whether messages are being loaded
    @Default(false) bool isLoading,

    /// Error message
    String? error,

    /// Quick replies from the assistant
    @Default([]) List<String> quickReplies,

    /// Whether the assistant is typing
    @Default(false) bool isTyping,
  }) = _ChatState;

  const ChatState._();

  /// Initial state
  factory ChatState.initial() => const ChatState();

  /// Check if connected to WebSocket
  bool get isConnected => wsState == WebSocketState.connected;

  /// Check if WebSocket is connecting
  bool get isConnecting => wsState == WebSocketState.connecting;

  /// Check if there are any messages
  bool get hasMessages => messages.isNotEmpty;

  /// Get the last message
  ChatMessage? get lastMessage => messages.isNotEmpty ? messages.last : null;
}
