import { Component, input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { Product } from '../../../../shared/types';

@Component({
  selector: 'app-product-card',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
  ],
  templateUrl: './product-card.component.html',
  styleUrl: './product-card.component.scss',
})
export class ProductCardComponent {
  product = input.required<Product>();
  index = input<number>();

  getImage(): string {
    const product = this.product() as any;
    return product.thumbnail || product.image || '';
  }

  getTitle(): string {
    const product = this.product() as any;
    return product.title || product.name || '';
  }

  getSource(): string {
    const product = this.product() as any;
    return product.source || product.merchant || '';
  }

  getRating(): string | undefined {
    const rating = this.product().rating;
    return rating ? `${rating}` : undefined;
  }

  getPrice(): string {
    return this.product().price || '';
  }

  getProductLink(): string {
    return this.product().product_link || this.product().link || '#';
  }

  openProductLink(event: Event): void {
    event.stopPropagation();
    const link = this.getProductLink();
    if (link && link !== '#') {
      window.open(link, '_blank', 'noopener,noreferrer');
    }
  }

  onCardClick(): void {
    // TODO: Open product drawer/modal with detailed information
    console.log('Product clicked:', this.product());
  }
}
