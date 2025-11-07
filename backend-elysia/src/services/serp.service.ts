/**
 * SerpService handles product search via Google Shopping API (SerpAPI)
 * Includes relevance scoring and result filtering
 */

import { getJson } from 'serpapi';
import type { Config } from '../config';
import type { KeyRotator } from '../utils/key-rotator';
import type { ProductCard } from '../types';

interface ShoppingItem {
  position: number;
  title: string;
  link: string;
  product_link?: string;
  product_id?: string;
  thumbnail?: string;
  price?: string;
  source?: string;
  rating?: number;
  reviews?: number;
  serpapi_product_api?: string;
  page_token?: string;
}

export class SerpService {
  private keyRotator: KeyRotator;
  private config: Config;

  constructor(keyRotator: KeyRotator, config: Config) {
    this.keyRotator = keyRotator;
    this.config = config;
  }

  /**
   * Search for products using Google Shopping
   */
  async searchProducts(
    query: string,
    searchType: string,
    country: string,
    minPrice?: number,
    maxPrice?: number
  ): Promise<ProductCard[]> {
    // Validate query
    if (!query || query.trim().length === 0) {
      throw new Error('Invalid search query');
    }

    const maxRetries = 2;
    let lastError: Error | null = null;

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      if (attempt > 0) {
        // Exponential backoff
        const backoffMs = 500 * Math.pow(2, attempt - 1);
        console.log(
          `⏳ SERP retry attempt ${attempt + 1}/${maxRetries + 1} after ${backoffMs}ms...`
        );
        await new Promise((resolve) => setTimeout(resolve, backoffMs));
      }

      try {
        const { key, index } = await this.keyRotator.getNextKey();

        if (attempt === 0) {
          console.log('\n🔍 SERP API Request:');
          console.log(`   Original Query: ${query}`);
          console.log(`   Type: ${searchType}`);
          console.log(`   Country: ${country}`);
          if (minPrice || maxPrice) {
            console.log(`   Price Range: ${minPrice || 'any'} - ${maxPrice || 'any'}`);
          }
        }
        console.log(`   Key Index: ${index} (attempt ${attempt + 1})`);

        const params: any = {
          engine: 'google_shopping',
          q: query,
          gl: country,
          hl: this.getLanguageForCountry(country),
          api_key: key,
        };

        if (minPrice) {
          params.min_price = Math.floor(minPrice).toString();
        }
        if (maxPrice) {
          params.max_price = Math.floor(maxPrice).toString();
        }

        const startTime = Date.now();
        const data = await getJson(params);
        const elapsed = (Date.now() - startTime) / 1000;

        if (attempt > 0) {
          console.log(`   ✅ SERP request succeeded on attempt ${attempt + 1}`);
        }
        console.log(`   ⏱️ Response time: ${elapsed.toFixed(2)}s`);

        const shoppingItems: ShoppingItem[] = [];

        if (data.shopping_results && Array.isArray(data.shopping_results)) {
          console.log(`   📦 Raw results: ${data.shopping_results.length} products`);

          for (const item of data.shopping_results) {
            shoppingItems.push({
              position: item.position || 0,
              title: item.title || '',
              link: item.link || '',
              product_link: item.product_link,
              product_id: item.product_id,
              thumbnail: item.thumbnail,
              price: item.price,
              source: item.source,
              rating: item.rating,
              reviews: item.reviews,
              serpapi_product_api: item.serpapi_product_api,
              page_token: item.immersive_product_page_token,
            });
          }
        } else {
          console.log('   ⚠️ No shopping_results in response');
        }

        const cards = this.convertToProductCards(shoppingItems, searchType);

        console.log(`   ✅ Found ${cards.length} products\n`);

        // Record successful usage
        await this.keyRotator.recordUsage(index, true, Date.now() - startTime);

        return cards;
      } catch (error: any) {
        lastError = error;
        console.error(
          `   ❌ SERP API Error (attempt ${attempt + 1}/${maxRetries + 1}):`,
          error.message
        );

        // Check if error is retryable
        const isRetryable =
          error.message?.includes('timeout') ||
          error.message?.includes('503') ||
          error.message?.includes('502') ||
          error.message?.includes('500');

        if (!isRetryable || attempt >= maxRetries) {
          throw new Error(`SERP API error: ${error.message}`);
        }
      }
    }

    throw new Error(
      `SERP API failed after ${maxRetries + 1} retries: ${lastError?.message}`
    );
  }

  /**
   * Convert shopping items to product cards
   */
  private convertToProductCards(
    items: ShoppingItem[],
    searchType: string
  ): ProductCard[] {
    const maxProducts = this.getMaxProducts(searchType);
    const products = items.slice(0, Math.min(maxProducts, 10));

    return products.map((item) => {
      // Extract old price if exists
      let oldPrice: string | undefined;
      if (item.price && item.price.includes('Was:')) {
        const match = item.price.match(/Was:\s*([^\s]+)/);
        if (match) {
          oldPrice = match[1];
        }
      }

      return {
        name: item.title,
        price: item.price || 'N/A',
        old_price: oldPrice,
        link: item.link,
        image: item.thumbnail || '',
        description: item.source || '',
        badge: this.getBadge(item),
        page_token: item.page_token || '',
      };
    });
  }

  /**
   * Get maximum products based on search type
   */
  private getMaxProducts(searchType: string): number {
    switch (searchType) {
      case 'exact':
        return this.config.serpMaxProductsExact;
      case 'parameters':
        return this.config.serpMaxProductsParameters;
      case 'category':
        return this.config.serpMaxProductsCategory;
      default:
        return this.config.serpMaxProductsDefault;
    }
  }

  /**
   * Get badge for product
   */
  private getBadge(item: ShoppingItem): string | undefined {
    if (item.rating && item.rating >= 4.5) {
      return 'Top Rated';
    }
    if (item.reviews && item.reviews >= 1000) {
      return 'Popular';
    }
    return undefined;
  }

  /**
   * Get language code for country
   */
  private getLanguageForCountry(country: string): string {
    const languageMap: Record<string, string> = {
      CH: 'de',
      DE: 'de',
      US: 'en',
      UK: 'en',
      FR: 'fr',
      ES: 'es',
      IT: 'it',
      RU: 'ru',
    };

    return languageMap[country] || 'en';
  }

  /**
   * Get product details by page token
   */
  async getProductDetails(pageToken: string, country: string): Promise<any> {
    try {
      const { key } = await this.keyRotator.getNextKey();

      const params = {
        engine: 'google_shopping_product',
        page_token: pageToken,
        gl: country,
        hl: this.getLanguageForCountry(country),
        api_key: key,
      };

      const data = await getJson(params);

      return {
        type: 'product_details',
        title: data.product_results?.title || '',
        price: data.product_results?.price || '',
        rating: data.product_results?.rating,
        reviews: data.product_results?.reviews,
        description: data.product_results?.description || '',
        images: data.product_results?.images || [],
        specifications: data.product_results?.specifications || [],
        variants: data.product_results?.variants || [],
        offers: data.product_results?.offers || [],
        videos: data.product_results?.videos || [],
        more_options: data.product_results?.more_options || [],
        rating_breakdown: data.product_results?.rating_breakdown || [],
      };
    } catch (error: any) {
      console.error('Error getting product details:', error);
      throw new Error(`Failed to get product details: ${error.message}`);
    }
  }
}
