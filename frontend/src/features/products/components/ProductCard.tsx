"use client";

import { useState } from "react";
import { ExternalLink, Star } from "lucide-react";
import { ProductCard as ProductCardType } from "@/types";
import { ProductDrawer } from "./ProductDrawer";

interface ProductCardProps {
  product: ProductCardType;
  index?: number;
}

export function ProductCard({ product, index }: ProductCardProps) {
  const [showModal, setShowModal] = useState(false);

  return (
    <>
      <div
        onClick={() => setShowModal(true)}
        className="group bg-card rounded-lg border border-border hover:border-primary/50 transition-all duration-500 cursor-pointer overflow-hidden h-full flex flex-col hover:shadow-lg hover:shadow-primary/10 hover:-translate-y-0.5"
      >
        {/* Image Container - Reduced aspect ratio */}
        <div className="aspect-[4/3] bg-muted relative overflow-hidden">
          {/* Gradient overlay on hover */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 z-10" />

          <img
            src={product.image}
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-[1.02] transition-transform duration-700 ease-out"
          />

          {/* Position Badge - Smaller */}
          {index && (
            <div className="absolute top-1.5 left-1.5 bg-gradient-to-br from-primary to-primary/80 text-primary-foreground w-7 h-7 rounded-full flex items-center justify-center text-[10px] font-bold shadow-md group-hover:scale-110 transition-all duration-300 z-20">
              {index}
            </div>
          )}

          {/* Rating Badge - Smaller */}
          {product.badge && (
            <div className="absolute top-1.5 right-1.5 bg-background/95 backdrop-blur-sm px-2 py-0.5 rounded-full text-[10px] font-semibold flex items-center gap-1 shadow-md border border-border/50 group-hover:border-primary/30 transition-all duration-300 z-20">
              <Star className="w-2.5 h-2.5 fill-yellow-500 text-yellow-500" />
              <span className="text-foreground">{product.badge.replace("⭐ ", "")}</span>
            </div>
          )}
        </div>

        {/* Content Container - More compact */}
        <div className="p-3 space-y-2 flex flex-col flex-grow bg-gradient-to-b from-card to-card/50">
          {/* Title - Smaller */}
          <h3 className="font-semibold text-xs leading-snug line-clamp-2 group-hover:text-primary transition-colors duration-300 min-h-[2rem]">
            {product.name}
          </h3>

          {/* Merchant - Smaller */}
          {product.description && (
            <p className="text-[10px] text-muted-foreground line-clamp-1 flex items-center gap-1 group-hover:text-foreground/70 transition-colors duration-300">
              <span className="w-1 h-1 rounded-full bg-primary/40 group-hover:bg-primary/60 transition-colors" />
              {product.description}
            </p>
          )}

          {/* Spacer to push price and button to bottom */}
          <div className="flex-grow" />

          {/* Price - Smaller */}
          <div className="flex items-baseline gap-1.5 pt-1.5 border-t border-border/50 group-hover:border-primary/30 transition-colors duration-300">
            <span className="text-base font-bold bg-gradient-to-r from-primary to-primary/80 bg-clip-text text-transparent group-hover:from-primary group-hover:to-primary transition-all duration-300">
              {product.price}
            </span>
            {product.old_price && (
              <span className="text-[10px] text-muted-foreground line-through group-hover:text-destructive/70 transition-colors duration-300">
                {product.old_price}
              </span>
            )}
          </div>

          {/* View Offer Link - More compact */}
          <a
            href={product.link}
            target="_blank"
            rel="noopener noreferrer"
            onClick={(e) => e.stopPropagation()}
            className="inline-flex items-center justify-center gap-1.5 text-xs font-medium text-primary-foreground bg-gradient-to-r from-primary to-primary/90 hover:from-primary hover:to-primary px-3 py-1.5 rounded-md transition-all duration-300 hover:shadow-lg hover:shadow-primary/20 hover:scale-[1.02] group/btn relative overflow-hidden"
          >
            {/* Shimmer effect */}
            <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent opacity-0 group-hover/btn:opacity-100 group-hover/btn:animate-shimmer" />

            <span className="relative">View Offer</span>
            <ExternalLink className="w-3 h-3 group-hover/btn:translate-x-0.5 transition-transform duration-300 relative" />
          </a>
        </div>
      </div>

      {showModal && (
        <ProductDrawer
          pageToken={product.page_token}
          onClose={() => setShowModal(false)}
        />
      )}
    </>
  );
}