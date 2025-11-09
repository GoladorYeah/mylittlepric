import { cn } from "@/lib/utils"

interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'circular' | 'text'
  animation?: 'pulse' | 'wave' | 'shimmer'
  delay?: number
}

export function Skeleton({
  className,
  variant = 'default',
  animation = 'shimmer',
  delay = 0,
  ...props
}: SkeletonProps) {
  const animationClass = {
    pulse: 'animate-pulse',
    wave: 'animate-wave',
    shimmer: 'animate-skeleton-shimmer'
  }[animation]

  const variantClass = {
    default: 'rounded-md',
    circular: 'rounded-full',
    text: 'rounded h-4'
  }[variant]

  return (
    <div
      className={cn(
        "bg-gradient-to-r from-muted via-muted/50 to-muted bg-[length:200%_100%]",
        variantClass,
        animationClass,
        className
      )}
      style={{ animationDelay: `${delay}ms` }}
      {...props}
    />
  )
}

// Specialized skeleton components
export function SkeletonText({
  lines = 1,
  className,
  delay = 0
}: {
  lines?: number
  className?: string
  delay?: number
}) {
  return (
    <div className={cn("space-y-2", className)}>
      {Array.from({ length: lines }).map((_, i) => (
        <Skeleton
          key={i}
          variant="text"
          className={i === lines - 1 ? "w-4/5" : "w-full"}
          delay={delay + i * 50}
        />
      ))}
    </div>
  )
}

export function SkeletonImage({
  className,
  aspectRatio = "aspect-square",
  delay = 0
}: {
  className?: string
  aspectRatio?: string
  delay?: number
}) {
  return (
    <Skeleton
      className={cn(aspectRatio, "w-full", className)}
      delay={delay}
    />
  )
}

export function SkeletonButton({
  className,
  delay = 0
}: {
  className?: string
  delay?: number
}) {
  return (
    <Skeleton
      className={cn("h-10 w-24 rounded-lg", className)}
      delay={delay}
    />
  )
}

export function SkeletonAvatar({
  size = 40,
  delay = 0
}: {
  size?: number
  delay?: number
}) {
  return (
    <Skeleton
      variant="circular"
      className="shrink-0"
      style={{ width: size, height: size, animationDelay: `${delay}ms` }}
    />
  )
}

export function SkeletonCard({
  className,
  delay = 0
}: {
  className?: string
  delay?: number
}) {
  return (
    <div className={cn("rounded-xl border border-border bg-card p-4 space-y-3", className)}>
      <SkeletonImage aspectRatio="aspect-video" delay={delay} />
      <div className="space-y-2">
        <Skeleton className="h-5 w-3/4" delay={delay + 100} />
        <SkeletonText lines={2} delay={delay + 150} />
      </div>
      <div className="flex justify-between items-center pt-2">
        <Skeleton className="h-6 w-20" delay={delay + 200} />
        <SkeletonButton delay={delay + 250} />
      </div>
    </div>
  )
}
