import { Skeleton, SkeletonImage, SkeletonText } from "./skeleton"

export default function ProductCardSkeleton({ delay = 0 }: { delay?: number }) {
  return (
    <div
      className="w-[300px] shrink-0 snap-start rounded-xl border border-border bg-card overflow-hidden"
      style={{ animationDelay: `${delay}ms` }}
    >
      {/* Image skeleton */}
      <div className="relative">
        <SkeletonImage aspectRatio="aspect-square" delay={delay} />
        {/* Badge skeletons */}
        <div className="absolute top-2 left-2">
          <Skeleton className="w-8 h-6 rounded-full" delay={delay + 50} />
        </div>
        <div className="absolute top-2 right-2">
          <Skeleton className="w-16 h-6 rounded-full" delay={delay + 100} />
        </div>
      </div>

      {/* Content skeleton */}
      <div className="p-4 space-y-3">
        {/* Title */}
        <div className="space-y-2">
          <Skeleton className="h-5 w-full" delay={delay + 150} />
          <Skeleton className="h-5 w-4/5" delay={delay + 200} />
        </div>

        {/* Description */}
        <SkeletonText lines={2} delay={delay + 250} />

        {/* Price and button */}
        <div className="flex items-center justify-between pt-2">
          <Skeleton className="h-7 w-24" delay={delay + 300} />
          <Skeleton className="h-10 w-28 rounded-lg" delay={delay + 350} />
        </div>
      </div>
    </div>
  )
}

export function OfferSkeleton({ delay = 0 }: { delay?: number }) {
  return (
    <div
      className="flex items-center gap-3 p-4 rounded-lg border border-border bg-card/50"
      style={{ animationDelay: `${delay}ms` }}
    >
      {/* Store logo */}
      <Skeleton variant="circular" className="w-12 h-12" delay={delay} />

      <div className="flex-1 space-y-2">
        {/* Store name */}
        <Skeleton className="h-4 w-32" delay={delay + 50} />
        {/* Condition */}
        <Skeleton className="h-3 w-20" delay={delay + 100} />
      </div>

      <div className="text-right space-y-2">
        {/* Price */}
        <Skeleton className="h-6 w-24" delay={delay + 150} />
        {/* Shipping */}
        <Skeleton className="h-3 w-20" delay={delay + 200} />
      </div>

      {/* Button */}
      <Skeleton className="h-9 w-24 rounded-lg" delay={delay + 250} />
    </div>
  )
}

export function ImageGallerySkeleton({ delay = 0 }: { delay?: number }) {
  return (
    <div className="space-y-4">
      {/* Main image */}
      <div className="relative overflow-hidden rounded-lg">
        <SkeletonImage aspectRatio="aspect-square" delay={delay} />

        {/* Navigation buttons */}
        <div className="absolute inset-0 flex items-center justify-between p-4">
          <Skeleton variant="circular" className="w-10 h-10" delay={delay + 100} />
          <Skeleton variant="circular" className="w-10 h-10" delay={delay + 100} />
        </div>
      </div>

      {/* Thumbnails */}
      <div className="grid grid-cols-5 gap-2">
        {Array.from({ length: 5 }).map((_, i) => (
          <SkeletonImage
            key={i}
            aspectRatio="aspect-square"
            className="rounded-md"
            delay={delay + 150 + i * 50}
          />
        ))}
      </div>
    </div>
  )
}

export function ProductInfoSkeleton({ delay = 0 }: { delay?: number }) {
  return (
    <div className="space-y-6">
      {/* Title and brand */}
      <div className="space-y-3">
        <Skeleton className="h-8 w-full" delay={delay} />
        <Skeleton className="h-8 w-3/4" delay={delay + 50} />
        <Skeleton className="h-5 w-32" delay={delay + 100} />
      </div>

      {/* Rating */}
      <div className="flex items-center gap-4">
        <div className="flex gap-1">
          {Array.from({ length: 5 }).map((_, i) => (
            <Skeleton key={i} className="w-5 h-5" delay={delay + 150 + i * 20} />
          ))}
        </div>
        <Skeleton className="h-4 w-24" delay={delay + 250} />
      </div>

      {/* Price */}
      <div className="space-y-2">
        <Skeleton className="h-10 w-40" delay={delay + 300} />
        <Skeleton className="h-4 w-56" delay={delay + 350} />
      </div>

      {/* Highlights */}
      <div className="space-y-3">
        <Skeleton className="h-6 w-32" delay={delay + 400} />
        <div className="space-y-2">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="flex items-start gap-2">
              <Skeleton className="w-5 h-5 mt-0.5" delay={delay + 450 + i * 50} />
              <Skeleton className="h-5 flex-1" delay={delay + 450 + i * 50} />
            </div>
          ))}
        </div>
      </div>

      {/* Description */}
      <div className="space-y-3">
        <Skeleton className="h-6 w-40" delay={delay + 650} />
        <SkeletonText lines={4} delay={delay + 700} />
      </div>

      {/* Specifications */}
      <div className="space-y-3">
        <Skeleton className="h-6 w-40" delay={delay + 800} />
        <div className="space-y-2">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="flex gap-2">
              <Skeleton className="h-5 w-32" delay={delay + 850 + i * 30} />
              <Skeleton className="h-5 flex-1" delay={delay + 850 + i * 30} />
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export function RatingBreakdownSkeleton({ delay = 0 }: { delay?: number }) {
  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Skeleton className="h-12 w-20" delay={delay} />
        <div className="flex-1 space-y-2">
          <Skeleton className="h-5 w-full" delay={delay + 50} />
          <Skeleton className="h-4 w-32" delay={delay + 100} />
        </div>
      </div>

      {/* Rating bars */}
      <div className="space-y-2">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="flex items-center gap-3">
            <Skeleton className="h-4 w-12" delay={delay + 150 + i * 40} />
            <Skeleton className="h-2 flex-1 rounded-full" delay={delay + 150 + i * 40} />
            <Skeleton className="h-4 w-16" delay={delay + 150 + i * 40} />
          </div>
        ))}
      </div>
    </div>
  )
}

export function ProductDrawerSkeleton() {
  return (
    <div className="space-y-8 animate-fade-in">
      {/* Image Gallery */}
      <ImageGallerySkeleton delay={0} />

      {/* Product Info */}
      <ProductInfoSkeleton delay={200} />

      {/* Offers Section */}
      <div className="space-y-4">
        <Skeleton className="h-6 w-40" delay={400} />
        <div className="space-y-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <OfferSkeleton key={i} delay={450 + i * 100} />
          ))}
        </div>
      </div>

      {/* Rating Breakdown */}
      <RatingBreakdownSkeleton delay={700} />

      {/* Similar Products */}
      <div className="space-y-4">
        <Skeleton className="h-6 w-48" delay={800} />
        <div className="grid grid-cols-2 gap-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="space-y-2">
              <SkeletonImage aspectRatio="aspect-square" className="rounded-lg" delay={850 + i * 50} />
              <Skeleton className="h-4 w-full" delay={900 + i * 50} />
              <Skeleton className="h-5 w-20" delay={950 + i * 50} />
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export function QuickReplySkeleton({ count = 3, delay = 0 }: { count?: number; delay?: number }) {
  return (
    <div className="flex flex-wrap gap-2 stagger-container">
      {Array.from({ length: count }).map((_, i) => (
        <Skeleton
          key={i}
          className="h-9 rounded-full"
          style={{
            width: `${80 + Math.random() * 100}px`,
            animationDelay: `${delay + i * 80}ms`
          }}
        />
      ))}
    </div>
  )
}
