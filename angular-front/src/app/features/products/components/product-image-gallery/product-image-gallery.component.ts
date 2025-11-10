import { Component, input, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-product-image-gallery',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatIconModule],
  templateUrl: './product-image-gallery.component.html',
  styleUrl: './product-image-gallery.component.scss',
})
export class ProductImageGalleryComponent {
  images = input.required<string[]>();
  title = input.required<string>();

  currentIndex = signal(0);

  nextImage(): void {
    const current = this.currentIndex();
    if (current < this.images().length - 1) {
      this.currentIndex.set(current + 1);
    }
  }

  prevImage(): void {
    const current = this.currentIndex();
    if (current > 0) {
      this.currentIndex.set(current - 1);
    }
  }

  selectImage(index: number): void {
    this.currentIndex.set(index);
  }

  getCurrentImage(): string {
    const images = this.images();
    return images[this.currentIndex()] || '/placeholder.png';
  }
}
