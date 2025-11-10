import { Component, input, output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { Product, ProductCard } from '../../../../shared/types';

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
  product = input.required<Product | ProductCard>();
  index = input<number>();

  // Output event when user wants to view product details
  productDetailsRequested = output<string>();

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
    const product = this.product() as any;
    const rating = product.rating;
    return rating ? `${rating}` : undefined;
  }

  getPrice(): string {
    return this.product().price || '';
  }

  getProductLink(): string {
    const product = this.product() as any;
    return product.product_link || product.link || '#';
  }

  getPageToken(): string | undefined {
    const product = this.product() as any;
    return product.page_token;
  }

  openProductLink(event: Event): void {
    event.stopPropagation();
    const link = this.getProductLink();
    if (link && link !== '#') {
      window.open(link, '_blank', 'noopener,noreferrer');
    }
  }

  onCardClick(): void {
    const pageToken = this.getPageToken();
    if (pageToken) {
      this.productDetailsRequested.emit(pageToken);
    } else {
      // Fallback: open product link if no page token
      const link = this.getProductLink();
      if (link && link !== '#') {
        window.open(link, '_blank', 'noopener,noreferrer');
      }
    }
  }
}
