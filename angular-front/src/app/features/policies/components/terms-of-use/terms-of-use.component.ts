import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { PolicyLayoutComponent } from '../policy-layout/policy-layout.component';

@Component({
  selector: 'app-terms-of-use',
  imports: [PolicyLayoutComponent, RouterLink],
  templateUrl: './terms-of-use.component.html',
  styleUrl: './terms-of-use.component.scss',
})
export class TermsOfUseComponent {
  title = 'Terms of Use';
  lastUpdated = '2025-01-06';
}
