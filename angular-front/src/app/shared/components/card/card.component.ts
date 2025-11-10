import { Component, input } from '@angular/core';
import { cn } from '../../utils';

@Component({
  selector: 'app-card',
  standalone: true,
  template: `
    <div [class]="cardClasses()">
      @if (title()) {
        <div class="card-header">
          <h3 class="card-title">{{ title() }}</h3>
          @if (subtitle()) {
            <p class="card-subtitle">{{ subtitle() }}</p>
          }
        </div>
      }
      <div class="card-content">
        <ng-content />
      </div>
      @if (hasFooter()) {
        <div class="card-footer">
          <ng-content select="[footer]" />
        </div>
      }
    </div>
  `,
  styles: [],
})
export class CardComponent {
  title = input<string>('');
  subtitle = input<string>('');
  padding = input<boolean>(true);
  hover = input<boolean>(false);
  className = input<string>('');
  hasFooter = input<boolean>(false);

  cardClasses() {
    const baseClasses = 'card';
    const paddingClass = this.padding() ? 'card-padded' : '';
    const hoverClass = this.hover() ? 'card-hover' : '';

    return cn(baseClasses, paddingClass, hoverClass, this.className());
  }
}
