import { Star } from "lucide-react";

interface RatingItem {
  stars: number;
  amount: number;
}

interface ProductRatingBreakdownProps {
  ratings: RatingItem[];
}

export function ProductRatingBreakdown({ ratings }: ProductRatingBreakdownProps) {
  if (!ratings || ratings.length === 0) return null;

  const maxAmount = Math.max(...ratings.map((r) => r.amount));

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-bold">Rating Breakdown</h2>
      <div className="space-y-2">
        {ratings.map((rating, index) => (
          <div key={index} className="flex items-center gap-3">
            <span className="w-12 text-sm font-medium">
              {rating.stars} <Star className="w-3 h-3 inline" />
            </span>
            <div className="flex-1 h-2 bg-muted rounded-full overflow-hidden">
              <div
                className="h-full bg-primary transition-all"
                style={{
                  width: `${(rating.amount / maxAmount) * 100}%`,
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
  );
}
