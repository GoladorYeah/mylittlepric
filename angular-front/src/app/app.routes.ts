import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./features/marketing/components/landing/landing.component').then(
        (m) => m.LandingComponent
      ),
  },
  {
    path: 'chat',
    loadComponent: () =>
      import('./features/chat/components/chat-page/chat-page.component').then(
        (m) => m.ChatPageComponent
      ),
  },
  {
    path: 'auth/callback',
    loadComponent: () =>
      import('./features/auth/components/auth-callback/auth-callback.component').then(
        (m) => m.AuthCallbackComponent
      ),
  },
  {
    path: '**',
    redirectTo: '',
  },
];
