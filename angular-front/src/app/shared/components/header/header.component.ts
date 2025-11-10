import { Component, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { AuthStore } from '../../../core/stores/auth.store';
import { ThemeService } from '../../../core/services/theme.service';
import { ButtonComponent } from '../button/button.component';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [RouterLink, ButtonComponent],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss',
})
export class HeaderComponent {
  readonly authStore = inject(AuthStore);
  readonly themeService = inject(ThemeService);

  onLogin(): void {
    this.authStore.login('google');
  }

  onLogout(): void {
    this.authStore.logout();
  }

  toggleTheme(): void {
    this.themeService.toggleTheme();
  }
}
