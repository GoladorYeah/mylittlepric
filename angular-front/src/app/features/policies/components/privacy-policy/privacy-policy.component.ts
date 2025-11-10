import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { PolicyLayoutComponent } from '../policy-layout/policy-layout.component';

@Component({
  selector: 'app-privacy-policy',
  imports: [PolicyLayoutComponent, RouterLink],
  templateUrl: './privacy-policy.component.html',
  styleUrl: './privacy-policy.component.scss',
})
export class PrivacyPolicyComponent {
  title = 'Privacy Policy';
  lastUpdated = '2025-01-06';
}
