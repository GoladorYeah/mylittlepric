import { Component, input, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { ProductOffer } from '../../../../shared/types';

@Component({
  selector: 'app-product-offers',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatIconModule],
  templateUrl: './product-offers.component.html',
  styleUrl: './product-offers.component.scss',
})
export class ProductOffersComponent {
  offers = input.required<ProductOffer[]>();

  bestPriceOffer = computed(() => {
    const offerList = this.offers();
    if (!offerList || offerList.length === 0) return null;

    return offerList.reduce((best, current) => {
      const bestPrice = best.extracted_total || best.extracted_price || Infinity;
      const currentPrice =
        current.extracted_total || current.extracted_price || Infinity;
      return currentPrice < bestPrice ? current : best;
    }, offerList[0]);
  });

  isBestPrice(offer: ProductOffer): boolean {
    const best = this.bestPriceOffer();
    return best === offer && this.offers().length > 1;
  }

  hasMonthlyPayment(offer: ProductOffer): boolean {
    return !!(
      offer.monthly_payment_duration && offer.monthly_payment_duration > 0
    );
  }

  formatNumber(num: number): string {
    return num.toLocaleString();
  }

  onImageError(event: Event): void {
    const target = event.target as HTMLImageElement;
    target.style.display = 'none';
  }
}
