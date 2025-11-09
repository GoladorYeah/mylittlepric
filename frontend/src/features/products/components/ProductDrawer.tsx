"use client";

import { useEffect, useState } from "react";
import { getProductDetails } from "@/lib/api";
import { ProductDetailsResponse } from "@/types";
import { useChatStore } from "@/lib/store";
import { Drawer } from "./ui/drawer";
import { ProductImageGallery } from "./product/product-image-gallery";
import { ProductInfo } from "./product/product-info";
import { ProductOffers } from "./product/product-offers";
import { ProductRatingBreakdown } from "./product/product-rating-breakdown";
import { ProductSimilarItems } from "./product/product-similar-items";
import { ProductDrawerSkeleton } from "./ui/product-skeletons";

interface ProductDrawerProps {
  pageToken: string;
  onClose: () => void;
}

export function ProductDrawer({ pageToken, onClose }: ProductDrawerProps) {
  const [product, setProduct] = useState<ProductDetailsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const { country } = useChatStore();

  useEffect(() => {
    const fetchDetails = async () => {
      try {
        const details = await getProductDetails(pageToken, country);
        setProduct(details);
      } catch (error) {
        console.error("Failed to load product details:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchDetails();
  }, [pageToken, country]);

  return (
    <Drawer isOpen={true} onClose={onClose}>
      {loading ? (
        <ProductDrawerSkeleton />
      ) : product ? (
        <div className="space-y-8 animate-fade-in">
          <ProductImageGallery
            images={product.images || []}
            title={product.title}
          />

          <ProductInfo
            title={product.title}
            price={product.price}
            rating={product.rating}
            reviews={product.reviews}
            description={product.description}
            specifications={product.specifications}
          />

          <ProductOffers offers={product.offers || []} />

          <ProductRatingBreakdown ratings={product.rating_breakdown || []} />

          <ProductSimilarItems products={product.more_options || []} />
        </div>
      ) : (
        <div className="flex items-center justify-center h-96 animate-fade-in">
          <div className="text-center space-y-2">
            <p className="text-muted-foreground text-lg">Failed to load product details</p>
            <button
              onClick={onClose}
              className="text-primary hover:underline text-sm"
            >
              Close and try again
            </button>
          </div>
        </div>
      )}
    </Drawer>
  );
}
