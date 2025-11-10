import { Component, input, output } from '@angular/core';
import { cn } from '../../utils';

type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive';
type ButtonSize = 'sm' | 'md' | 'lg';

@Component({
  selector: 'app-button',
  standalone: true,
  template: `
    <button
      [class]="buttonClasses()"
      [type]="type()"
      [disabled]="disabled()"
      (click)="handleClick($event)"
    >
      <ng-content />
    </button>
  `,
  styles: [],
})
export class ButtonComponent {
  variant = input<ButtonVariant>('primary');
  size = input<ButtonSize>('md');
  type = input<'button' | 'submit' | 'reset'>('button');
  disabled = input<boolean>(false);
  fullWidth = input<boolean>(false);
  click = output<MouseEvent>();

  buttonClasses() {
    const baseClasses = 'btn';
    const variantClass = `btn-${this.variant()}`;
    const sizeClass = `btn-${this.size()}`;
    const widthClass = this.fullWidth() ? 'w-full' : '';

    return cn(baseClasses, variantClass, sizeClass, widthClass);
  }

  handleClick(event: MouseEvent): void {
    if (!this.disabled()) {
      this.click.emit(event);
    }
  }
}
