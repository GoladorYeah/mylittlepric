import { Component, OnInit, inject } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { AuthStore } from '../../../../core/stores/auth.store';

@Component({
  selector: 'app-auth-callback',
  standalone: true,
  template: `
    <div class="auth-callback">
      <div class="spinner"></div>
      <p>{{ message }}</p>
    </div>
  `,
  styles: [
    `
      .auth-callback {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 100vh;
        gap: var(--spacing-lg);
      }

      .spinner {
        width: 3rem;
        height: 3rem;
        border: 4px solid var(--color-border);
        border-top-color: var(--color-primary);
        border-radius: 50%;
        animation: spin 1s linear infinite;
      }

      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }

      p {
        color: var(--color-muted-foreground);
        font-size: var(--font-size-lg);
      }
    `,
  ],
})
export class AuthCallbackComponent implements OnInit {
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly authStore = inject(AuthStore);

  message = 'Processing authentication...';

  async ngOnInit(): Promise<void> {
    try {
      // Wait a bit for cookies to be set
      await new Promise((resolve) => setTimeout(resolve, 500));

      // Check authentication status
      await this.authStore.checkAuth();

      if (this.authStore.isAuthenticated()) {
        this.message = 'Authentication successful! Redirecting...';

        // Load user preferences
        await this.authStore.loadPreferences();

        // Redirect to chat page after a short delay
        setTimeout(() => {
          this.router.navigate(['/chat']);
        }, 1000);
      } else {
        this.message = 'Authentication failed. Redirecting...';
        setTimeout(() => {
          this.router.navigate(['/']);
        }, 2000);
      }
    } catch (error) {
      console.error('Auth callback error:', error);
      this.message = 'An error occurred. Redirecting...';
      setTimeout(() => {
        this.router.navigate(['/']);
      }, 2000);
    }
  }
}
