import { ExternalLink, Star } from "lucide-react";

interface Offer {
  merchant: string;
  price: string;
  link: string;
  shipping?: string;
  rating?: number | string;
  currency?: string;
  availability?: string;
}

interface ProductOffersProps {
  offers: Offer[];
}

export function ProductOffers({ offers }: ProductOffersProps) {
  if (!offers || offers.length === 0) return null;

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-bold">Available Offers</h2>
      <div className="space-y-3">
        {offers.map((offer, index) => (
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
  );
}
