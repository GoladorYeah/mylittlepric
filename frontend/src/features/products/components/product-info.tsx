import { Star } from "lucide-react";

interface ProductInfoProps {
  title: string;
  price: string;
  rating?: number;
  reviews?: number;
  description?: string;
  specifications?: Array<{ title: string; value: string }>;
}

export function ProductInfo({
  title,
  price,
  rating,
  reviews,
  description,
  specifications,
}: ProductInfoProps) {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">{title}</h1>

      {rating && (
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-1">
            <Star className="w-5 h-5 fill-current text-yellow-500" />
            <span className="font-semibold">{rating}</span>
          </div>
          {reviews && (
            <span className="text-sm text-muted-foreground">
              ({reviews.toLocaleString()} reviews)
            </span>
          )}
        </div>
      )}

      <div className="text-3xl font-bold text-primary">{price}</div>

      {description && (
        <div className="space-y-2">
          <h3 className="font-semibold">Description</h3>
          <p className="text-muted-foreground">{description}</p>
        </div>
      )}

      {specifications && specifications.length > 0 && (
        <div className="space-y-2">
          <h3 className="font-semibold">Specifications</h3>
          <div className="grid grid-cols-2 gap-2 text-sm">
            {specifications.map((spec, index) => (
              <div key={index} className="flex flex-col">
                <span className="text-muted-foreground">{spec.title}</span>
                <span className="font-medium">{spec.value}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
