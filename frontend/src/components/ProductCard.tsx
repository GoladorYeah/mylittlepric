"use client";

import { useState } from "react";
import { ExternalLink, Star } from "lucide-react";
import { ProductCard as ProductCardType } from "@/types";
import { ProductDrawer } from "./ProductDrawer";

interface ProductCardProps {
  product: ProductCardType;
}

export function ProductCard({ product }: ProductCardProps) {
  const [showModal, setShowModal] = useState(false);

  return (
    <>
      <div
        onClick={() => setShowModal(true)}
        className="group bg-secondary rounded-xl border border-border hover:border-primary transition-all cursor-pointer overflow-hidden"
      >
        <div className="aspect-square bg-muted relative overflow-hidden">
          <img
            src={product.image}
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
          {product.badge && (
            <div className="absolute top-2 right-2 bg-background/90 backdrop-blur-sm px-2 py-1 rounded-full text-xs font-medium flex items-center gap-1">
              <Star className="w-3 h-3 fill-current" />
              {product.badge.replace("⭐ ", "")}
            </div>
          )}
        </div>

        <div className="p-4 space-y-2">
          <h3 className="font-semibold line-clamp-2 group-hover:text-primary transition-colors">
            {product.name}
          </h3>

          {product.description && (
            <p className="text-sm text-muted-foreground line-clamp-1">
              {product.description}
            </p>
          )}

          <div className="flex items-baseline gap-2">
            <span className="text-lg font-bold">{product.price}</span>
            {product.old_price && (
              <span className="text-sm text-muted-foreground line-through">
      {product.old_price}
              </span>
            )}
          </div>

          <a
            href={product.link}
            target="_blank"
            rel="noopener noreferrer"
            onClick={(e) => e.stopPropagation()}
            className="flex items-center gap-2 text-sm text-primary hover:underline"
          >
            View offer
            <ExternalLink className="w-4 h-4" />
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