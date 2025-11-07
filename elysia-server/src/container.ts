/**
 * Dependency Injection Container
 * Initializes and provides access to all services
 */

import type { Config } from './config';
import { initRedis, initPostgres } from './utils/database';
import { KeyRotator } from './utils/key-rotator';
import { JWTService } from './utils/jwt';
import { EmbeddingService } from './services/embedding.service';
import { CacheService } from './services/cache.service';
import { GroundingStrategy } from './services/grounding-strategy.service';
import { GeminiService } from './services/gemini.service';
import { SerpService } from './services/serp.service';
import { SessionService } from './services/session.service';
import { GoogleOAuthService, AuthService } from './services/auth.service';
import { SearchHistoryService } from './services/search-history.service';
import type { Redis } from 'ioredis';

export class Container {
  public config: Config;
  public redis: Redis;
  public db: any;

  // Utilities
  public geminiRotator: KeyRotator;
  public serpRotator: KeyRotator;
  public jwtService: JWTService;

  // Services
  public embeddingService: EmbeddingService;
  public cacheService: CacheService;
  public groundingStrategy: GroundingStrategy;
  public geminiService: GeminiService;
  public serpService: SerpService;
  public sessionService: SessionService;
  public googleOAuthService: GoogleOAuthService;
  public authService: AuthService;
  public searchHistoryService: SearchHistoryService;

  constructor(config: Config) {
    this.config = config;

    // Initialize databases
    console.log('🚀 Initializing databases...');
    this.redis = initRedis(config);
    this.db = initPostgres(config);

    // Initialize key rotators
    console.log('🔑 Initializing key rotators...');
    this.geminiRotator = new KeyRotator('gemini', config.geminiApiKeys, this.redis);
    this.serpRotator = new KeyRotator('serp', config.serpApiKeys, this.redis);

    console.log(`✅ Gemini Key Rotator: ${this.geminiRotator.getTotalKeys()} keys`);
    console.log(`✅ SERP Key Rotator: ${this.serpRotator.getTotalKeys()} keys`);

    // Initialize JWT service
    console.log('🔐 Initializing JWT service...');
    this.jwtService = new JWTService(
      config.jwtAccessSecret,
      config.jwtRefreshSecret,
      config.jwtAccessTTL,
      config.jwtRefreshTTL
    );

    // Initialize services
    console.log('🛠️ Initializing services...');

    // Get first Gemini key for embedding service
    const { key: geminiKey } = this.geminiRotator.getKeyByIndex(0);

    this.embeddingService = new EmbeddingService(geminiKey, this.redis, config);
    console.log('✅ Embedding Service initialized');

    this.cacheService = new CacheService(this.redis, config, this.embeddingService);
    console.log('✅ Cache Service initialized');

    this.groundingStrategy = new GroundingStrategy(this.embeddingService, config);
    console.log('✅ Grounding Strategy initialized');

    this.geminiService = new GeminiService(
      this.geminiRotator,
      config,
      this.embeddingService,
      this.groundingStrategy
    );
    console.log(`✅ Gemini Service initialized (model: ${config.geminiModel})`);
    console.log(`🎯 Smart Grounding: '${config.geminiGroundingMode}' mode`);
    if (config.geminiUseGrounding) {
      console.log('🔍 Grounding: ENABLED (selective usage)');
    } else {
      console.log('💬 Grounding: DISABLED globally');
    }

    this.serpService = new SerpService(this.serpRotator, config);
    console.log('✅ SERP Service initialized');

    this.sessionService = new SessionService(
      this.redis,
      config.sessionTTL,
      config.maxMessagesPerSession
    );
    console.log('✅ Session Service initialized');

    this.googleOAuthService = new GoogleOAuthService(config);
    console.log('✅ Google OAuth Service initialized');

    this.authService = new AuthService(
      this.redis,
      this.jwtService,
      this.googleOAuthService
    );
    console.log('✅ Auth Service initialized');

    this.searchHistoryService = new SearchHistoryService(this.db);
    console.log('✅ Search History Service initialized');

    console.log('✅ All services initialized successfully');
  }

  /**
   * Health check
   */
  async healthCheck(): Promise<any> {
    try {
      await this.redis.ping();

      return {
        status: 'healthy',
        redis: 'ok',
        gemini_keys: {
          count: this.geminiRotator.getTotalKeys(),
          status: 'ok',
        },
        serp_keys: {
          count: this.serpRotator.getTotalKeys(),
          status: 'ok',
        },
        grounding: {
          mode: this.config.geminiGroundingMode,
          enabled: this.config.geminiUseGrounding,
        },
        embedding: {
          status: 'ok',
        },
      };
    } catch (error) {
      return {
        status: 'unhealthy',
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  /**
   * Close all connections
   */
  async close(): Promise<void> {
    console.log('🛑 Shutting down container...');

    if (this.redis) {
      await this.redis.quit();
    }

    if (this.db) {
      await this.db.end();
    }

    console.log('✅ Container closed gracefully');
  }
}
