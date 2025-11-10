import { Component, input, output, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { cn } from '../../utils';

@Component({
  selector: 'app-input',
  standalone: true,
  imports: [FormsModule],
  template: `
    <div [class]="containerClasses()">
      @if (label()) {
        <label [for]="id()" class="input-label">
          {{ label() }}
        </label>
      }
      <input
        [id]="id()"
        [type]="type()"
        [placeholder]="placeholder()"
        [disabled]="disabled()"
        [value]="value()"
        [class]="inputClasses()"
        (input)="handleInput($event)"
        (blur)="handleBlur()"
        (focus)="handleFocus()"
      />
      @if (error()) {
        <span class="input-error">{{ error() }}</span>
      }
    </div>
  `,
  styles: [],
})
export class InputComponent {
  id = input<string>(`input-${Math.random().toString(36).substring(7)}`);
  type = input<string>('text');
  label = input<string>('');
  placeholder = input<string>('');
  value = input<string>('');
  disabled = input<boolean>(false);
  error = input<string>('');
  fullWidth = input<boolean>(false);

  valueChange = output<string>();
  blur = output<void>();
  focus = output<void>();

  private isFocused = signal(false);

  containerClasses() {
    const widthClass = this.fullWidth() ? 'w-full' : '';
    return cn('input-container', widthClass);
  }

  inputClasses() {
    const baseClasses = 'input';
    const errorClass = this.error() ? 'input-error-state' : '';
    const focusedClass = this.isFocused() ? 'input-focused' : '';

    return cn(baseClasses, errorClass, focusedClass);
  }

  handleInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.valueChange.emit(target.value);
  }

  handleBlur(): void {
    this.isFocused.set(false);
    this.blur.emit();
  }

  handleFocus(): void {
    this.isFocused.set(true);
    this.focus.emit();
  }
}
