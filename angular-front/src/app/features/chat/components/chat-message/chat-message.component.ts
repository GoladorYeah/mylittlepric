import { Component, input, output, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { ChatMessage } from '../../../../shared/types';
import { formatRelativeTime } from '../../../../shared/utils';
import { AuthStore } from '../../../../core/stores/auth.store';
import { ProductCardComponent } from '../../../products/components/product-card/product-card.component';

@Component({
  selector: 'app-chat-message',
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    ProductCardComponent,
  ],
  templateUrl: './chat-message.component.html',
  styleUrl: './chat-message.component.scss',
})
export class ChatMessageComponent {
  message = input.required<ChatMessage>();
  onQuickReply = input<(reply: string) => void>();
  productDetailsRequested = output<string>();

  private authStore = inject(AuthStore);

  // Make Math available in template
  protected readonly Math = Math;

  formatTime(timestamp: number): string {
    return formatRelativeTime(timestamp);
  }

  handleQuickReply(reply: string): void {
    const callback = this.onQuickReply();
    if (callback) {
      callback(reply);
    }
  }

  getUserInitials(): string {
    const user = this.authStore.user();
    if (!user) return 'U';

    const fullName = (user as any).full_name;
    if (fullName) {
      const names = fullName.trim().split(/\s+/);
      if (names.length >= 2) {
        return (names[0][0] + names[names.length - 1][0]).toUpperCase();
      }
      return names[0][0].toUpperCase();
    }

    return user.email[0].toUpperCase();
  }

  parseQuickReply(reply: string): { text: string; price: string | null } {
    const priceMatch = reply.match(/\(([≈~]?[A-Z$€£¥]{1,4}[\s]?[\d,.\-–—]+[\+]?(?:[\s]?[kK]|[\s]?[\-–—][\s]?[\d,.\-–—]+[\+]?(?:[kK])?)?)\)$/);

    if (priceMatch) {
      const text = reply.substring(0, priceMatch.index).trim();
      const price = priceMatch[1];
      return { text, price };
    }

    return { text: reply, price: null };
  }

  scrollToProduct(direction: 'left' | 'right', container: HTMLElement): void {
    const scrollAmount = 224; // Width of card (210px) + gap (14px)
    const newScrollLeft = direction === 'left'
      ? container.scrollLeft - scrollAmount
      : container.scrollLeft + scrollAmount;

    container.scrollTo({
      left: newScrollLeft,
      behavior: 'smooth'
    });
  }

  handleProductDetailsRequest(pageToken: string): void {
    this.productDetailsRequested.emit(pageToken);
  }
}
