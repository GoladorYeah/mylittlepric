import {
  Component,
  OnInit,
  OnDestroy,
  inject,
  signal,
  effect,
  ViewChild,
  ElementRef,
} from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { ChatStore } from '../../../../core/stores/chat.store';
import { AuthStore } from '../../../../core/stores/auth.store';
import { WebSocketService } from '../../../../core/services/websocket.service';
import { HeaderComponent } from '../../../../shared/components/header/header.component';
import { ChatMessageComponent } from '../chat-message/chat-message.component';
import { ProductDrawerComponent } from '../../../products/components/product-drawer/product-drawer.component';
import { Product } from '../../../../shared/types';
import { detectCountry, detectLanguage, getCurrencyForCountry } from '../../../../shared/utils/locale';

@Component({
  selector: 'app-chat-page',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatButtonModule,
    MatIconModule,
    MatInputModule,
    MatFormFieldModule,
    MatProgressSpinnerModule,
    HeaderComponent,
    ChatMessageComponent,
    ProductDrawerComponent,
  ],
  templateUrl: './chat-page.component.html',
  styleUrl: './chat-page.component.scss',
})
export class ChatPageComponent implements OnInit, OnDestroy {
  @ViewChild('messagesContainer') messagesContainer?: ElementRef<HTMLDivElement>;

  private readonly route = inject(ActivatedRoute);
  readonly chatStore = inject(ChatStore);
  readonly authStore = inject(AuthStore);
  readonly wsService = inject(WebSocketService);

  messageInput = signal('');
  isConnecting = signal(false);
  isInputFocused = signal(false);

  // Product drawer state
  selectedPageToken = signal<string>('');
  isDrawerOpen = signal(false);

  constructor() {
    // Auto-scroll when messages change
    effect(() => {
      if (this.chatStore.messages()) {
        this.scrollToBottom();
      }
    });
  }

  ngOnInit(): void {
    // Check for query parameter
    this.route.queryParams.subscribe((params) => {
      const query = params['q'];
      if (query && !this.chatStore.hasMessages()) {
        this.messageInput.set(query);
        // Auto-send after a short delay
        setTimeout(() => this.sendMessage(), 500);
      }
    });

    // Connect WebSocket
    this.connectWebSocket();
  }

  ngOnDestroy(): void {
    this.wsService.disconnect();
  }

  private async connectWebSocket(): Promise<void> {
    const sessionId = this.chatStore.sessionId();
    const token = this.authStore.accessToken();
    const country = this.authStore.preferences().country || await detectCountry();
    const language = this.authStore.preferences().language || detectLanguage();
    const currency = this.authStore.preferences().currency || getCurrencyForCountry(country);

    this.isConnecting.set(true);
    this.wsService.connect(sessionId, token || undefined, country, language, currency);

    // Wait for connection
    const checkConnection = setInterval(() => {
      if (this.wsService.isConnected()) {
        this.isConnecting.set(false);
        clearInterval(checkConnection);
      }
    }, 100);

    // Timeout after 5 seconds
    setTimeout(() => {
      if (!this.wsService.isConnected()) {
        this.isConnecting.set(false);
        clearInterval(checkConnection);
      }
    }, 5000);
  }

  async sendMessage(): Promise<void> {
    const content = this.messageInput().trim();
    if (!content) return;

    const token = this.authStore.accessToken();
    const country = this.authStore.preferences().country || await detectCountry();
    const language = this.authStore.preferences().language || detectLanguage();
    const currency = this.authStore.preferences().currency || getCurrencyForCountry(country);

    // Try WebSocket first, fallback to HTTP
    if (this.wsService.isConnected()) {
      this.chatStore.sendMessage(
        content,
        token || undefined,
        country,
        language,
        currency
      );
    } else {
      this.chatStore.sendMessageHttp(
        content,
        token || undefined,
        country,
        language,
        currency
      );
    }

    this.messageInput.set('');
  }

  handleQuickReply(reply: string): void {
    this.messageInput.set(reply);
    this.sendMessage();
  }

  startNewChat(): void {
    this.chatStore.startNewSession();
    this.connectWebSocket();
  }

  private scrollToBottom(): void {
    setTimeout(() => {
      const container = this.messagesContainer?.nativeElement;
      if (container) {
        container.scrollTop = container.scrollHeight;
      }
    }, 100);
  }

  // Get all products from messages
  getAllProducts(): Product[] {
    const products: Product[] = [];
    for (const message of this.chatStore.messages()) {
      if (message.products) {
        products.push(...message.products);
      }
    }
    return products;
  }

  // Track products by their position to avoid re-rendering
  trackByPosition(index: number, product: Product): number {
    return product.position;
  }

  // Handle product details request
  handleProductDetailsRequest(pageToken: string): void {
    this.selectedPageToken.set(pageToken);
    this.isDrawerOpen.set(true);
  }

  // Handle drawer close
  handleDrawerClose(): void {
    this.isDrawerOpen.set(false);
    // Clear page token after animation
    setTimeout(() => this.selectedPageToken.set(''), 300);
  }
}
