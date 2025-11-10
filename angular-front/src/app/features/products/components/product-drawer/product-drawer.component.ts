import {
  Component,
  input,
  output,
  signal,
  inject,
  OnInit,
  effect,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { ApiService } from '../../../../core/services/api.service';
import { AuthStore } from '../../../../core/stores/auth.store';
import { ProductDetailsResponse } from '../../../../shared/types';
import { ProductImageGalleryComponent } from '../product-image-gallery/product-image-gallery.component';
import { ProductInfoComponent } from '../product-info/product-info.component';
import { ProductOffersComponent } from '../product-offers/product-offers.component';
import { ProductRatingBreakdownComponent } from '../product-rating-breakdown/product-rating-breakdown.component';
import { ProductSimilarItemsComponent } from '../product-similar-items/product-similar-items.component';
import { ProductDrawerSkeletonComponent } from '../product-drawer-skeleton/product-drawer-skeleton.component';

@Component({
  selector: 'app-product-drawer',
  standalone: true,
  imports: [
    CommonModule,
    MatSidenavModule,
    MatButtonModule,
    MatIconModule,
    ProductImageGalleryComponent,
    ProductInfoComponent,
    ProductOffersComponent,
    ProductRatingBreakdownComponent,
    ProductSimilarItemsComponent,
    ProductDrawerSkeletonComponent,
  ],
  templateUrl: './product-drawer.component.html',
  styleUrl: './product-drawer.component.scss',
})
export class ProductDrawerComponent implements OnInit {
  pageToken = input.required<string>();
  isOpen = input.required<boolean>();
  closed = output<void>();

  private apiService = inject(ApiService);
  private authStore = inject(AuthStore);

  product = signal<ProductDetailsResponse | null>(null);
  loading = signal(true);
  error = signal(false);

  constructor() {
    // Watch for pageToken changes and refetch
    effect(() => {
      const token = this.pageToken();
      if (token && this.isOpen()) {
        this.fetchProductDetails();
      }
    });
  }

  ngOnInit(): void {
    if (this.pageToken() && this.isOpen()) {
      this.fetchProductDetails();
    }
  }

  private fetchProductDetails(): void {
    this.loading.set(true);
    this.error.set(false);

    const accessToken = this.authStore.accessToken();

    this.apiService.getProductDetails(this.pageToken(), accessToken ?? undefined).subscribe({
      next: (details) => {
        this.product.set(details);
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to load product details:', err);
        this.error.set(true);
        this.loading.set(false);
      },
    });
  }

  handleClose(): void {
    this.closed.emit();
  }

  retry(): void {
    this.fetchProductDetails();
  }
}
