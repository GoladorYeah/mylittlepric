import { Component, signal } from '@angular/core';
import { Router } from '@angular/router';
import { ButtonComponent } from '../../../../shared/components/button/button.component';
import { ThemeToggleComponent } from '../../../../shared/components/theme-toggle/theme-toggle.component';

@Component({
  selector: 'app-landing',
  imports: [ButtonComponent, ThemeToggleComponent],
  templateUrl: './landing.component.html',
  styleUrl: './landing.component.scss',
})
export class LandingComponent {
  searchQuery = signal('');
  isSearchAnimating = signal(false);

  constructor(private router: Router) {}

  handleSearch(): void {
    const query = this.searchQuery().trim();

    if (query) {
      this.isSearchAnimating.set(true);
      setTimeout(() => {
        this.router.navigate(['/chat'], {
          queryParams: { q: encodeURIComponent(query) }
        });
      }, 500);
    } else {
      this.router.navigate(['/chat']);
    }
  }

  navigateToChat(): void {
    this.router.navigate(['/chat']);
  }
}
