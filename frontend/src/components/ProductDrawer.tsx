"use client";

import { useEffect, useState } from "react";
import { X, ExternalLink, Star, ChevronLeft, ChevronRight } from "lucide-react";
import { getProductDetails } from "@/lib/api";
import { ProductDetailsResponse } from "@/types";
import { useChatStore } from "@/lib/store";

interface ProductDrawerProps {
  pageToken: string;
  onClose: () => void;
}

export function ProductDrawer({ pageToken, onClose }: ProductDrawerProps) {
  const [product, setProduct] = useState<ProductDetailsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [isClosing, setIsClosing] = useState(false);
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

  const nextImage = () => {
    if (product?.images && currentImageIndex < product.images.length - 1) {
      setCurrentImageIndex(currentImageIndex + 1);
    }
  };

  const prevImage = () => {
    if (currentImageIndex > 0) {
      setCurrentImageIndex(currentImageIndex - 1);
    }
  };

  const handleClose = () => {
    setIsClosing(true);
    setTimeout(() => {
      onClose();
    }, 300); // Match animation duration
  };

  useEffect(() => {
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = "unset";
    };
  }, []);

  // Handle Escape key
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        handleClose();
      }
    };
    window.addEventListener("keydown", handleEscape);
    return () => window.removeEventListener("keydown", handleEscape);
  }, []);

  return (
    <div className="fixed inset-0 z-50 flex">
      {/* Backdrop */}
      <div
        className={`absolute inset-0 bg-background/80 backdrop-blur-sm transition-opacity duration-300 ${
          isClosing ? "opacity-0" : "opacity-100"
        }`}
        onClick={handleClose}
      />

      {/* Drawer */}
      <div
        className={`ml-auto relative bg-background border-l border-border w-full max-w-2xl h-full overflow-y-auto shadow-2xl transition-transform duration-300 ease-in-out ${
          isClosing ? "translate-x-full" : "translate-x-0"
        }`}
      >
        {/* Close button */}
        <div className="sticky top-0 z-10 bg-background/95 backdrop-blur-sm border-b border-border p-4">
          <button
            onClick={handleClose}
            className="w-10 h-10 rounded-full bg-secondary hover:bg-secondary/80 flex items-center justify-center transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {loading ? (
          <div className="flex items-center justify-center h-96">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
          </div>
        ) : product ? (
          <div className="p-6 space-y-6">
            {/* Images */}
            <div className="space-y-4">
              <div className="relative aspect-square bg-muted rounded-xl overflow-hidden">
                <img
                  src={
                    product.images?.[currentImageIndex] ||
                    "/placeholder.png"
                  }
                  alt={product.title}
                  className="w-full h-full object-cover"
                />
                {product.images && product.images.length > 1 && (
                  <>
                    <button
                      onClick={prevImage}
                      disabled={currentImageIndex === 0}
                      className="absolute left-2 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-background/80 backdrop-blur-sm hover:bg-background disabled:opacity-50 flex items-center justify-center transition-all"
                    >
                      <ChevronLeft className="w-5 h-5" />
                    </button>
                    <button
                      onClick={nextImage}
                      disabled={currentImageIndex === product.images.length - 1}
                      className="absolute right-2 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-background/80 backdrop-blur-sm hover:bg-background disabled:opacity-50 flex items-center justify-center transition-all"
                    >
                      <ChevronRight className="w-5 h-5" />
                    </button>
                  </>
                )}
              </div>

              {product.images && product.images.length > 1 && (
                <div className="flex gap-2 overflow-x-auto pb-2">
                  {product.images.map((img, index) => (
                    <button
                      key={index}
                      onClick={() => setCurrentImageIndex(index)}
                      className={`flex-shrink-0 w-20 h-20 rounded-lg overflow-hidden border-2 transition-colors ${
                        currentImageIndex === index
                          ? "border-primary"
                          : "border-border"
                      }`}
                    >
                      <img
                        src={img}
                        alt={`${product.title} ${index + 1}`}
                        className="w-full h-full object-cover"
                      />
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Product Info */}
            <div className="space-y-4">
              <h1 className="text-2xl font-bold">{product.title}</h1>

              {product.rating && (
                <div className="flex items-center gap-2">
                  <div className="flex items-center gap-1">
                    <Star className="w-5 h-5 fill-current text-yellow-500" />
                    <span className="font-semibold">{product.rating}</span>
                  </div>
                  {product.reviews && (
                    <span className="text-sm text-muted-foreground">
                      ({product.reviews.toLocaleString()} reviews)
                    </span>
                  )}
                </div>
              )}

              <div className="text-3xl font-bold text-primary">
                {product.price}
              </div>

              {product.description && (
                <div className="space-y-2">
                  <h3 className="font-semibold">Description</h3>
                  <p className="text-muted-foreground">{product.description}</p>
                </div>
              )}

              {product.specifications && product.specifications.length > 0 && (
                <div className="space-y-2">
                  <h3 className="font-semibold">Specifications</h3>
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    {product.specifications.map((spec, index) => (
                      <div key={index} className="flex flex-col">
                        <span className="text-muted-foreground">
                          {spec.title}
                        </span>
                        <span className="font-medium">{spec.value}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* Offers */}
            {product.offers && product.offers.length > 0 && (
              <div className="space-y-4">
                <h2 className="text-xl font-bold">Available Offers</h2>
                <div className="space-y-3">
                  {product.offers.map((offer, index) => (
                    <div
                      key={index}
                      className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 rounded-xl bg-secondary border border-border"
                    >
                      <div className="space-y-1">
                        <div className="font-semibold">{offer.merchant}</div>
                        <div className="text-2xl font-bold">{offer.price}</div>
                        {offer.shipping && (
                          <div className="text-sm text-muted-foreground">
                            {offer.shipping}
                          </div>
                        )}
                        {offer.rating && (
                          <div className="flex items-center gap-1 text-sm">
                            <Star className="w-4 h-4 fill-current text-yellow-500" />
                            <span>{offer.rating}</span>
                          </div>
                        )}
                      </div>
                      <a
                        href={offer.link}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="px-6 py-3 bg-primary text-primary-foreground rounded-full font-semibold hover:opacity-90 transition-opacity flex items-center justify-center gap-2 whitespace-nowrap"
                      >
                        View Offer
                        <ExternalLink className="w-4 h-4" />
                      </a>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Rating Breakdown */}
            {product.rating_breakdown && product.rating_breakdown.length > 0 && (
              <div className="space-y-4">
                <h2 className="text-xl font-bold">Rating Breakdown</h2>
                <div className="space-y-2">
                  {product.rating_breakdown.map((rating, index) => (
                    <div key={index} className="flex items-center gap-3">
                      <span className="w-12 text-sm font-medium">
                        {rating.stars} <Star className="w-3 h-3 inline" />
                      </span>
                      <div className="flex-1 h-2 bg-muted rounded-full overflow-hidden">
                        <div
                          className="h-full bg-primary transition-all"
                          style={{
                            width: `${
                              (rating.amount /
                                Math.max(
                                  ...product.rating_breakdown!.map((r) => r.amount)
                                )) *
                              100
                            }%`,
                          }}
                        />
                      </div>
                      <span className="w-16 text-sm text-muted-foreground text-right">
                        {rating.amount.toLocaleString()}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Similar Products */}
            {product.more_options && product.more_options.length > 0 && (
              <div className="space-y-4">
                <h2 className="text-xl font-bold">Similar Products</h2>
                <div className="grid grid-cols-2 gap-4">
                  {product.more_options.map((option, index) => (
                    <div
                      key={index}
                      className="space-y-2 p-3 rounded-xl bg-secondary border border-border hover:border-primary transition-colors cursor-pointer"
                    >
                      <div className="aspect-square bg-muted rounded-lg overflow-hidden">
                        <img
                          src={option.thumbnail}
                          alt={option.title}
                          className="w-full h-full object-cover"
                        />
                      </div>
                      <div className="text-sm font-medium line-clamp-2">
                        {option.title}
                      </div>
                      <div className="text-lg font-bold">{option.price}</div>
                      {option.rating && (
                        <div className="flex items-center gap-1 text-xs">
                          <Star className="w-3 h-3 fill-current text-yellow-500" />
                          <span>{option.rating}</span>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        ) : (
          <div className="flex items-center justify-center h-96">
            <p className="text-muted-foreground">Failed to load product details</p>
          </div>
        )}
      </div>
    </div>
  );
}
