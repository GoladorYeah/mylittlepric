import { Component, input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';

export interface SimilarProduct {
  title: string;
  price: string;
  thumbnail: string;
  rating?: string;
}

@Component({
  selector: 'app-product-similar-items',
  standalone: true,
  imports: [CommonModule, MatIconModule],
  templateUrl: './product-similar-items.component.html',
  styleUrl: './product-similar-items.component.scss',
})
export class ProductSimilarItemsComponent {
  products = input.required<SimilarProduct[]>();
}
