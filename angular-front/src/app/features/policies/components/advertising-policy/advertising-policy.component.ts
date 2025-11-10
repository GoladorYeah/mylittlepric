import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { PolicyLayoutComponent } from '../policy-layout/policy-layout.component';

@Component({
  selector: 'app-advertising-policy',
  imports: [PolicyLayoutComponent, RouterLink],
  templateUrl: './advertising-policy.component.html',
  styleUrl: './advertising-policy.component.scss',
})
export class AdvertisingPolicyComponent {
  title = 'Advertiser & Seller Advertising Policy';
  lastUpdated = '2025-01-06';
}
