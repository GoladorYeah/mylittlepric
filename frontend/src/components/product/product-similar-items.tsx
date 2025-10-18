import { Star } from "lucide-react";

interface SimilarProduct {
  title: string;
  price: string;
  thumbnail: string;
  rating?: string;
}

interface ProductSimilarItemsProps {
  products: SimilarProduct[];
}

export function ProductSimilarItems({ products }: ProductSimilarItemsProps) {
  if (!products || products.length === 0) return null;

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-bold">Similar Products</h2>
      <div className="grid grid-cols-2 gap-4">
        {products.map((product, index) => (
          <div
            key={index}
            className="space-y-2 p-3 rounded-xl bg-secondary border border-border hover:border-primary transition-colors cursor-pointer"
          >
            <div className="aspect-square bg-muted rounded-lg overflow-hidden">
              <img
                src={product.thumbnail}
                alt={product.title}
                className="w-full h-full object-cover"
              />
            </div>
            <div className="text-sm font-medium line-clamp-2">
              {product.title}
            </div>
            <div className="text-lg font-bold">{product.price}</div>
            {product.rating && (
              <div className="flex items-center gap-1 text-xs">
                <Star className="w-3 h-3 fill-current text-yellow-500" />
                <span>{product.rating}</span>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
