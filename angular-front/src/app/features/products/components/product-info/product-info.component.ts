import { Component, input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';

export interface ProductSpecification {
  title: string;
  value: string;
}

@Component({
  selector: 'app-product-info',
  standalone: true,
  imports: [CommonModule, MatIconModule],
  templateUrl: './product-info.component.html',
  styleUrl: './product-info.component.scss',
})
export class ProductInfoComponent {
  title = input.required<string>();
  price = input.required<string>();
  rating = input<number>();
  reviews = input<number>();
  description = input<string>();
  specifications = input<ProductSpecification[]>();

  formatNumber(num: number): string {
    return num.toLocaleString();
  }
}
