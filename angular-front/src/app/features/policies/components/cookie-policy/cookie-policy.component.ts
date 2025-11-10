import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { PolicyLayoutComponent } from '../policy-layout/policy-layout.component';

@Component({
  selector: 'app-cookie-policy',
  imports: [PolicyLayoutComponent, RouterLink],
  templateUrl: './cookie-policy.component.html',
  styleUrl: './cookie-policy.component.scss',
})
export class CookiePolicyComponent {
  title = 'Cookie Policy';
  lastUpdated = '2025-01-06';
}
