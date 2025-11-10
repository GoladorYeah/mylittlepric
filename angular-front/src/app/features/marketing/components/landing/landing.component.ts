import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { HeaderComponent } from '../../../../shared/components/header/header.component';
import { ButtonComponent } from '../../../../shared/components/button/button.component';
import { CardComponent } from '../../../../shared/components/card/card.component';

@Component({
  selector: 'app-landing',
  standalone: true,
  imports: [HeaderComponent, ButtonComponent, CardComponent],
  templateUrl: './landing.component.html',
  styleUrl: './landing.component.scss',
})
export class LandingComponent {
  constructor(private router: Router) {}

  navigateToChat(): void {
    this.router.navigate(['/chat']);
  }
}
