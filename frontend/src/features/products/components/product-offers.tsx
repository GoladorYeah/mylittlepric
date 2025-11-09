import { ExternalLink, Star, Check, CreditCard, Package, Tag as TagIcon } from "lucide-react";
import Image from "next/image";

interface Offer {
  merchant: string;
  logo?: string;
  price: string;
  extracted_price?: number;
  currency?: string;
  link: string;
  title?: string;
  availability?: string;
  shipping?: string;
  shipping_extracted?: number;
  total?: string;
  extracted_total?: number;
  rating?: number | string;
  reviews?: number;
  payment_methods?: string;
  tag?: string;
  details_and_offers?: string[];
  monthly_payment_duration?: number;
  down_payment?: string;
}

interface ProductOffersProps {
  offers: Offer[];
}

export function ProductOffers({ offers }: ProductOffersProps) {
  if (!offers || offers.length === 0) return null;

  // Find best price offer
  const bestPriceOffer = offers.reduce((best, current) => {
    const bestPrice = best.extracted_total || best.extracted_price || Infinity;
    const currentPrice = current.extracted_total || current.extracted_price || Infinity;
    return currentPrice < bestPrice ? current : best;
  }, offers[0]);

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-bold">
        Available from {offers.length} {offers.length === 1 ? 'store' : 'stores'}
      </h2>
      <div className="space-y-3">
        {offers.map((offer, index) => {
          const isBestPrice = offer === bestPriceOffer && offers.length > 1;
          const hasMonthlyPayment = offer.monthly_payment_duration && offer.monthly_payment_duration > 0;

          return (
            <div
              key={index}
              className={`relative flex flex-col gap-4 p-4 rounded-xl border transition-all ${
                isBestPrice
                  ? 'bg-primary/5 border-primary/40 shadow-md'
                  : 'bg-secondary border-border hover:border-primary/30'
              }`}
            >
              {/* Best Price Tag */}
              {(isBestPrice || offer.tag) && (
                <div className="absolute -top-2 -right-2 px-3 py-1 rounded-full bg-primary text-primary-foreground text-xs font-bold shadow-lg flex items-center gap-1">
                  <TagIcon className="w-3 h-3" />
                  {offer.tag || 'Best price'}
                </div>
              )}

              <div className="flex flex-col sm:flex-row gap-4">
                {/* Store Logo & Info */}
                <div className="flex-1 space-y-3">
                  <div className="flex items-center gap-3">
                    {offer.logo && (
                      <div className="relative w-8 h-8 flex-shrink-0">
                        <Image
                          src={offer.logo}
                          alt={offer.merchant}
                          width={32}
                          height={32}
                          className="object-contain"
                          unoptimized
                          onError={(e) => {
                            const target = e.target as HTMLImageElement;
                            target.style.display = 'none';
                          }}
                        />
                      </div>
                    )}
                    <div className="flex-1 min-w-0">
                      <div className="font-semibold text-lg truncate">{offer.merchant}</div>
                      {offer.title && (
                        <div className="text-sm text-muted-foreground line-clamp-1">
                          {offer.title}
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Pricing */}
                  <div className="space-y-1">
                    {hasMonthlyPayment ? (
                      <div>
                        <div className="text-3xl font-bold text-primary">
                          {offer.price}
                          <span className="text-sm font-normal text-muted-foreground ml-2">
                            /mo for {offer.monthly_payment_duration} months
                          </span>
                        </div>
                        {offer.down_payment && offer.down_payment !== "$0.00" && (
                          <div className="text-sm text-muted-foreground">
                            Down payment: {offer.down_payment}
                          </div>
                        )}
                        {offer.total && (
                          <div className="text-sm text-muted-foreground">
                            Total: {offer.total}
                          </div>
                        )}
                      </div>
                    ) : (
                      <div className="text-3xl font-bold text-primary">{offer.price}</div>
                    )}

                    <div className="flex flex-wrap gap-2 text-sm">
                      {offer.shipping && (
                        <div className={offer.shipping.toLowerCase().includes('free') ? 'text-green-600 dark:text-green-500 font-medium' : 'text-muted-foreground'}>
                          <Package className="w-3.5 h-3.5 inline mr-1" />
                          {offer.shipping}
                        </div>
                      )}
                      {offer.total && !hasMonthlyPayment && (
                        <div className="text-muted-foreground">
                          Total: {offer.total}
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Rating & Reviews */}
                  {(offer.rating || offer.reviews) && (
                    <div className="flex items-center gap-3 text-sm">
                      {offer.rating && (
                        <div className="flex items-center gap-1 text-yellow-600 dark:text-yellow-500">
                          <Star className="w-4 h-4 fill-current" />
                          <span className="font-medium">{offer.rating}</span>
                        </div>
                      )}
                      {offer.reviews && (
                        <div className="text-muted-foreground">
                          ({offer.reviews.toLocaleString()} {offer.reviews === 1 ? 'review' : 'reviews'})
                        </div>
                      )}
                    </div>
                  )}

                  {/* Payment Methods */}
                  {offer.payment_methods && (
                    <div className="text-sm text-muted-foreground flex items-center gap-1">
                      <CreditCard className="w-3.5 h-3.5" />
                      {offer.payment_methods}
                    </div>
                  )}

                  {/* Details and Offers */}
                  {offer.details_and_offers && offer.details_and_offers.length > 0 && (
                    <div className="space-y-1">
                      {offer.details_and_offers.map((detail, idx) => (
                        <div key={idx} className="text-sm text-muted-foreground flex items-start gap-2">
                          <Check className="w-4 h-4 text-green-600 dark:text-green-500 flex-shrink-0 mt-0.5" />
                          <span>{detail}</span>
                        </div>
                      ))}
                    </div>
                  )}

                  {/* Availability */}
                  {offer.availability && (
                    <div className="text-sm text-green-600 dark:text-green-500 font-medium">
                      {offer.availability}
                    </div>
                  )}
                </div>

                {/* View Offer Button */}
                <a
                  href={offer.link}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="px-6 py-3 bg-primary text-primary-foreground rounded-full font-semibold hover:opacity-90 transition-opacity flex items-center justify-center gap-2 whitespace-nowrap self-start sm:self-center h-fit"
                >
                  View Offer
                  <ExternalLink className="w-4 h-4" />
                </a>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
