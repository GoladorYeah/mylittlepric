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
    path: 'privacy-policy',
    loadComponent: () =>
      import('./features/policies/components/privacy-policy/privacy-policy.component').then(
        (m) => m.PrivacyPolicyComponent
      ),
  },
  {
    path: 'terms-of-use',
    loadComponent: () =>
      import('./features/policies/components/terms-of-use/terms-of-use.component').then(
        (m) => m.TermsOfUseComponent
      ),
  },
  {
    path: 'cookie-policy',
    loadComponent: () =>
      import('./features/policies/components/cookie-policy/cookie-policy.component').then(
        (m) => m.CookiePolicyComponent
      ),
  },
  {
    path: 'advertising-policy',
    loadComponent: () =>
      import('./features/policies/components/advertising-policy/advertising-policy.component').then(
        (m) => m.AdvertisingPolicyComponent
      ),
  },
  {
    path: 'history',
    loadComponent: () =>
      import('./features/history/components/history-page/history-page.component').then(
        (m) => m.HistoryPageComponent
      ),
  },
  {
    path: 'settings',
    loadComponent: () =>
      import('./features/settings/components/settings-page/settings-page.component').then(
        (m) => m.SettingsPageComponent
      ),
  },
  {
    path: '**',
    redirectTo: '',
  },
];
