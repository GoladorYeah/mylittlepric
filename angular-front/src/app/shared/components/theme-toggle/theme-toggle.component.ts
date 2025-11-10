import { Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ThemeService } from '../../../core/services/theme.service';

@Component({
  selector: 'app-theme-toggle',
  standalone: true,
  imports: [MatButtonModule, MatIconModule, MatTooltipModule],
  templateUrl: './theme-toggle.component.html',
  styleUrl: './theme-toggle.component.scss',
})
export class ThemeToggleComponent {
  readonly themeService = inject(ThemeService);

  getIcon(): string {
    const theme = this.themeService.actualTheme();
    return theme === 'dark' ? 'light_mode' : 'dark_mode';
  }

  getLabel(): string {
    const theme = this.themeService.actualTheme();
    return theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode';
  }

  toggleTheme(): void {
    this.themeService.toggleTheme();
  }
}
