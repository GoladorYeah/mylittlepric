import { Component, input } from '@angular/core';
import { ChatMessage } from '../../../../shared/types';
import { formatRelativeTime } from '../../../../shared/utils';
import { ButtonComponent } from '../../../../shared/components/button/button.component';

@Component({
  selector: 'app-chat-message',
  standalone: true,
  imports: [ButtonComponent],
  templateUrl: './chat-message.component.html',
  styleUrl: './chat-message.component.scss',
})
export class ChatMessageComponent {
  message = input.required<ChatMessage>();
  onQuickReply = input<(reply: string) => void>();

  formatTime(timestamp: number): string {
    return formatRelativeTime(timestamp);
  }

  handleQuickReply(reply: string): void {
    const callback = this.onQuickReply();
    if (callback) {
      callback(reply);
    }
  }
}
