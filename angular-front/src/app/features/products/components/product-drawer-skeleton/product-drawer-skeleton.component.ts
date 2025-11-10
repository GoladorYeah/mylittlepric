import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-product-drawer-skeleton',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './product-drawer-skeleton.component.html',
  styleUrl: './product-drawer-skeleton.component.scss',
})
export class ProductDrawerSkeletonComponent {
  // Array helpers for *ngFor in template
  thumbnails = Array(5).fill(0);
  highlights = Array(4).fill(0);
  specs = Array(6).fill(0);
  offers = Array(3).fill(0);
  ratings = Array(5).fill(0);
  similarProducts = Array(4).fill(0);
}
