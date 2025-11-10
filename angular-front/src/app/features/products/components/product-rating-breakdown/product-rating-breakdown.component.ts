import { Component, input, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';

export interface RatingItem {
  stars: number;
  amount: number;
}

@Component({
  selector: 'app-product-rating-breakdown',
  standalone: true,
  imports: [CommonModule, MatIconModule],
  templateUrl: './product-rating-breakdown.component.html',
  styleUrl: './product-rating-breakdown.component.scss',
})
export class ProductRatingBreakdownComponent {
  ratings = input.required<RatingItem[]>();

  maxAmount = computed(() => {
    const ratingList = this.ratings();
    if (!ratingList || ratingList.length === 0) return 0;
    return Math.max(...ratingList.map((r) => r.amount));
  });

  getPercentageWidth(amount: number): string {
    const max = this.maxAmount();
    if (max === 0) return '0%';
    return `${(amount / max) * 100}%`;
  }

  formatNumber(num: number): string {
    return num.toLocaleString();
  }
}
