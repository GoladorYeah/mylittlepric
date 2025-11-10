import 'dart:async';
import 'dart:convert';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:logger/logger.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:web_socket_channel/status.dart' as status;
import '../config/app_config.dart';
import '../storage/storage_service.dart';

/// WebSocket client with automatic reconnection and heartbeat
class WebSocketClient {
  final Logger _logger = Logger();

  WebSocketChannel? _channel;
  StreamSubscription? _subscription;
  Timer? _heartbeatTimer;
  Timer? _reconnectTimer;

  int _reconnectAttempts = 0;
  bool _isConnecting = false;
  bool _isClosed = false;

  /// Stream controller for incoming messages
  final _messageController = StreamController<Map<String, dynamic>>.broadcast();

  /// Stream controller for connection state
  final _connectionStateController = StreamController<WebSocketState>.broadcast();

  WebSocketClient();

  /// Stream of incoming messages
  Stream<Map<String, dynamic>> get messages => _messageController.stream;

  /// Stream of connection state changes
  Stream<WebSocketState> get connectionState => _connectionStateController.stream;

  /// Current connection state
  WebSocketState _state = WebSocketState.disconnected;
  WebSocketState get state => _state;

  /// Connect to WebSocket server
  Future<void> connect() async {
    if (_isConnecting || _state == WebSocketState.connected) {
      _logger.w('Already connected or connecting');
      return;
    }

    _isConnecting = true;
    _isClosed = false;
    _updateState(WebSocketState.connecting);

    try {
      // Build WebSocket URL
      final wsUrl = await _buildWebSocketUrl();
      _logger.d('Connecting to WebSocket: $wsUrl');

      // Create WebSocket channel
      _channel = WebSocketChannel.connect(Uri.parse(wsUrl));

      // Listen to messages
      _subscription = _channel!.stream.listen(
        _handleMessage,
        onError: _handleError,
        onDone: _handleClose,
        cancelOnError: false,
      );

      // Start heartbeat
      _startHeartbeat();

      // Reset reconnect attempts
      _reconnectAttempts = 0;

      _updateState(WebSocketState.connected);
      _logger.i('WebSocket connected successfully');
    } catch (e) {
      _logger.e('WebSocket connection failed', error: e);
      _isConnecting = false;
      _updateState(WebSocketState.error);
      _scheduleReconnect();
    }

    _isConnecting = false;
  }

  /// Disconnect from WebSocket server
  Future<void> disconnect() async {
    _logger.d('Disconnecting WebSocket');
    _isClosed = true;

    _stopHeartbeat();
    _cancelReconnect();

    await _subscription?.cancel();
    await _channel?.sink.close(status.goingAway);

    _updateState(WebSocketState.disconnected);
  }

  /// Send message to WebSocket server
  void send(Map<String, dynamic> data) {
    if (_state != WebSocketState.connected) {
      _logger.w('Cannot send message: WebSocket not connected');
      throw Exception('WebSocket not connected');
    }

    try {
      final message = jsonEncode(data);
      _channel?.sink.add(message);
      _logger.d('Sent message: $message');
    } catch (e) {
      _logger.e('Failed to send message', error: e);
      rethrow;
    }
  }

  /// Build WebSocket URL with session ID
  Future<String> _buildWebSocketUrl() async {
    final baseUrl = AppConfig.wsBaseUrl;
    final endpoint = AppConfig.wsEndpoint;

    // Get session ID from storage
    final prefs = StorageService.prefs;
    final sessionId = prefs.getString(AppConfig.sessionIdKey);

    // Build URL with session ID if available
    var url = '$baseUrl$endpoint';
    if (sessionId != null && sessionId.isNotEmpty) {
      url += '?session_id=$sessionId';
    }

    return url;
  }

  /// Handle incoming WebSocket message
  void _handleMessage(dynamic data) {
    try {
      final message = jsonDecode(data as String) as Map<String, dynamic>;
      _logger.d('Received message: $message');

      // Handle ping/pong
      if (message['type'] == 'ping') {
        send({'type': 'pong'});
        return;
      }

      // Emit message to stream
      _messageController.add(message);
    } catch (e) {
      _logger.e('Failed to parse message', error: e);
    }
  }

  /// Handle WebSocket error
  void _handleError(dynamic error) {
    _logger.e('WebSocket error', error: error);
    _updateState(WebSocketState.error);
    _scheduleReconnect();
  }

  /// Handle WebSocket close
  void _handleClose() {
    _logger.w('WebSocket closed');
    _stopHeartbeat();

    if (!_isClosed) {
      _updateState(WebSocketState.disconnected);
      _scheduleReconnect();
    }
  }

  /// Start heartbeat timer
  void _startHeartbeat() {
    _stopHeartbeat();

    _heartbeatTimer = Timer.periodic(
      AppConfig.wsPingInterval,
      (_) {
        if (_state == WebSocketState.connected) {
          try {
            send({'type': 'ping'});
          } catch (e) {
            _logger.e('Heartbeat failed', error: e);
          }
        }
      },
    );
  }

  /// Stop heartbeat timer
  void _stopHeartbeat() {
    _heartbeatTimer?.cancel();
    _heartbeatTimer = null;
  }

  /// Schedule reconnection attempt
  void _scheduleReconnect() {
    if (_isClosed || _reconnectAttempts >= AppConfig.wsMaxReconnectAttempts) {
      _logger.w('Max reconnect attempts reached');
      _updateState(WebSocketState.disconnected);
      return;
    }

    _cancelReconnect();

    // Calculate delay with exponential backoff
    final delay = AppConfig.wsReconnectDelay * (_reconnectAttempts + 1);
    _reconnectAttempts++;

    _logger.d('Scheduling reconnect in ${delay.inSeconds}s (attempt $_reconnectAttempts)');

    _reconnectTimer = Timer(delay, () {
      if (!_isClosed) {
        connect();
      }
    });
  }

  /// Cancel reconnection timer
  void _cancelReconnect() {
    _reconnectTimer?.cancel();
    _reconnectTimer = null;
  }

  /// Update connection state
  void _updateState(WebSocketState newState) {
    if (_state != newState) {
      _state = newState;
      _connectionStateController.add(newState);
      _logger.d('WebSocket state changed: $newState');
    }
  }

  /// Dispose resources
  Future<void> dispose() async {
    await disconnect();
    await _messageController.close();
    await _connectionStateController.close();
  }
}

/// WebSocket connection state
enum WebSocketState {
  /// Not connected
  disconnected,

  /// Connecting to server
  connecting,

  /// Connected and ready
  connected,

  /// Error occurred
  error,
}

/// Provider for WebSocketClient
final webSocketClientProvider = Provider<WebSocketClient>((ref) {
  final client = WebSocketClient();

  // Dispose when provider is disposed
  ref.onDispose(() {
    client.dispose();
  });

  return client;
});
