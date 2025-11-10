import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-policy-layout',
  imports: [RouterLink, MatIconModule, MatButtonModule],
  templateUrl: './policy-layout.component.html',
  styleUrl: './policy-layout.component.scss',
})
export class PolicyLayoutComponent {
  title = input.required<string>();
  lastUpdated = input.required<string>();
}
