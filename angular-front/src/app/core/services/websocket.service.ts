import { Injectable, signal, effect } from '@angular/core';
import { Subject, Observable, timer } from 'rxjs';
import { environment } from '../../../environments/environment';
import { WebSocketMessage, ChatMessage } from '../../shared/types';
import { generateId } from '../../shared/utils';

export enum WebSocketConnectionState {
  CONNECTING = 'CONNECTING',
  OPEN = 'OPEN',
  CLOSING = 'CLOSING',
  CLOSED = 'CLOSED',
}

@Injectable({
  providedIn: 'root',
})
export class WebSocketService {
  private ws: WebSocket | null = null;
  private messageSubject = new Subject<WebSocketMessage>();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private heartbeatInterval: any = null;
  private sessionId: string | null = null;
  private accessToken: string | null = null;
  private country: string | null = null;
  private language: string | null = null;
  private currency: string | null = null;

  // Signals for reactive state
  readonly connectionState = signal<WebSocketConnectionState>(
    WebSocketConnectionState.CLOSED
  );
  readonly isConnected = signal(false);
  readonly lastError = signal<string | null>(null);

  constructor() {
    // Set up effect to update isConnected based on connectionState
    effect(() => {
      const state = this.connectionState();
      this.isConnected.set(state === WebSocketConnectionState.OPEN);
    });
  }

  get messages$(): Observable<WebSocketMessage> {
    return this.messageSubject.asObservable();
  }

  private getWebSocketUrl(): string {
    if (environment.wsUrl) {
      return environment.wsUrl;
    }

    // Construct WebSocket URL from current location
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    return `${protocol}//${host}/ws`;
  }

  connect(
    sessionId: string,
    accessToken?: string,
    country?: string,
    language?: string,
    currency?: string
  ): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      console.log('[WebSocket] Already connected');
      return;
    }

    this.sessionId = sessionId;
    this.accessToken = accessToken || null;
    this.country = country || null;
    this.language = language || null;
    this.currency = currency || null;

    this.connectionState.set(WebSocketConnectionState.CONNECTING);
    this.lastError.set(null);

    try {
      const url = this.getWebSocketUrl();
      console.log('[WebSocket] Connecting to:', url);

      this.ws = new WebSocket(url);

      this.ws.onopen = () => {
        console.log('[WebSocket] Connected');
        this.connectionState.set(WebSocketConnectionState.OPEN);
        this.reconnectAttempts = 0;
        this.startHeartbeat();
      };

      this.ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          console.log('[WebSocket] Received:', message);
          this.messageSubject.next(message);
        } catch (error) {
          console.error('[WebSocket] Failed to parse message:', error);
        }
      };

      this.ws.onerror = (error) => {
        console.error('[WebSocket] Error:', error);
        this.lastError.set('WebSocket connection error');
      };

      this.ws.onclose = (event) => {
        console.log('[WebSocket] Closed:', event.code, event.reason);
        this.connectionState.set(WebSocketConnectionState.CLOSED);
        this.stopHeartbeat();

        if (!event.wasClean && this.reconnectAttempts < this.maxReconnectAttempts) {
          this.scheduleReconnect();
        }
      };
    } catch (error) {
      console.error('[WebSocket] Failed to create connection:', error);
      this.lastError.set('Failed to create WebSocket connection');
      this.connectionState.set(WebSocketConnectionState.CLOSED);
    }
  }

  private scheduleReconnect(): void {
    this.reconnectAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);

    console.log(
      `[WebSocket] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`
    );

    timer(delay).subscribe(() => {
      if (this.sessionId) {
        this.connect(
          this.sessionId,
          this.accessToken || undefined,
          this.country || undefined,
          this.language || undefined,
          this.currency || undefined
        );
      }
    });
  }

  disconnect(): void {
    if (this.ws) {
      console.log('[WebSocket] Disconnecting...');
      this.connectionState.set(WebSocketConnectionState.CLOSING);
      this.stopHeartbeat();
      this.ws.close(1000, 'Client disconnecting');
      this.ws = null;
      this.sessionId = null;
      this.accessToken = null;
    }
  }

  sendMessage(
    message: string,
    sessionId: string,
    accessToken?: string,
    country?: string,
    language?: string,
    currency?: string,
    newSearch?: boolean
  ): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('[WebSocket] Cannot send message: not connected');
      this.lastError.set('WebSocket is not connected');
      return;
    }

    const payload: WebSocketMessage = {
      type: 'chat',
      session_id: sessionId,
      message,
      access_token: accessToken,
      country,
      language,
      currency,
      new_search: newSearch,
    };

    try {
      console.log('[WebSocket] Sending:', payload);
      this.ws.send(JSON.stringify(payload));
    } catch (error) {
      console.error('[WebSocket] Failed to send message:', error);
      this.lastError.set('Failed to send message');
    }
  }

  private startHeartbeat(): void {
    this.stopHeartbeat();
    this.heartbeatInterval = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        try {
          this.ws.send(JSON.stringify({ type: 'ping' }));
        } catch (error) {
          console.error('[WebSocket] Failed to send heartbeat:', error);
        }
      }
    }, 30000); // Every 30 seconds
  }

  private stopHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }

  // Update session context
  updateContext(
    sessionId: string,
    accessToken?: string,
    country?: string,
    language?: string,
    currency?: string
  ): void {
    this.sessionId = sessionId;
    this.accessToken = accessToken || null;
    this.country = country || null;
    this.language = language || null;
    this.currency = currency || null;
  }
}
