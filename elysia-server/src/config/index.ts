/**
 * Configuration module for MyLittlePrice Backend
 * Handles all environment variables and application settings
 */

export interface Config {
  // Server
  port: string;
  env: string;

  // Database
  databaseUrl: string;

  // Redis
  redisUrl: string;
  redisPassword: string;
  redisDb: number;

  // JWT Authentication
  jwtAccessSecret: string;
  jwtRefreshSecret: string;
  jwtAccessTTL: number; // seconds
  jwtRefreshTTL: number; // seconds

  // Google OAuth
  googleClientId: string;
  googleClientSecret: string;
  googleRedirectUrl: string;

  // Session
  sessionTTL: number;
  maxMessagesPerSession: number;
  maxSearchesPerSession: number;

  // API Keys
  geminiApiKeys: string[];
  serpApiKeys: string[];

  // Gemini Configuration
  geminiModel: string;
  geminiFallbackModel: string;
  geminiTemperature: number;
  geminiMaxOutputTokens: number;
  geminiUseGrounding: boolean;

  // Grounding Strategy Settings
  geminiGroundingMode: 'conservative' | 'balanced' | 'aggressive';
  geminiGroundingMinWords: number;

  // Grounding Strategy Thresholds
  geminiBrandQueryConfidence: number;
  geminiGroundingWeightFreshInfo: number;
  geminiGroundingWeightSpecific: number;
  geminiGroundingWeightDrift: number;
  geminiGroundingWeightElectron: number;
  geminiGroundingDecisionThresh: number;
  geminiBrandQueryMaxWords: number;
  geminiBrandSimilarityThresh: number;
  geminiDialogueHistoryWindow: number;
  geminiDialogueDriftThresh: number;
  geminiDriftScoreBonus: number;
  geminiElectronicsThreshHigh: number;
  geminiElectronicsScoreHigh: number;
  geminiCategorySimilarityThresh: number;
  geminiCategoryScore: number;

  // Gemini Translation Settings
  geminiTranslationTemperature: number;
  geminiTranslationMaxTokens: number;

  // Embedding Settings
  geminiEmbeddingModel: string;
  embeddingCategoryDetectionThresh: number;
  cacheQueryEmbeddingTTL: number;

  // SERP Relevance Thresholds
  serpThresholdExact: number;
  serpThresholdParameters: number;
  serpThresholdCategory: number;
  serpLogTopResultsCount: number;
  serpFallbackMinResults: number;

  // SERP Scoring Weights
  serpScorePhraseMatch: number;
  serpScoreAllWords: number;
  serpScorePartialWords: number;
  serpScoreWordOrderWeight: number;
  serpScoreBrandMatch: number;
  serpScoreModelMatch: number;
  serpMinWordLength: number;
  serpModelNumberMinLength: number;

  // SERP Max Products
  serpMaxProductsExact: number;
  serpMaxProductsParameters: number;
  serpMaxProductsCategory: number;
  serpMaxProductsDefault: number;

  // Default Values
  defaultCountry: string;
  defaultLanguage: string;
  defaultCurrency: string;

  // Cache TTL (seconds)
  cacheGeminiTTL: number;
  cacheSerpTTL: number;
  cacheImmersiveTTL: number;

  // Rate Limiting
  rateLimitRequests: number;
  rateLimitWindow: number;

  // CORS
  corsOrigins: string[];

  // Logging
  logLevel: string;
}

function getEnv(key: string, defaultValue: string = ''): string {
  return process.env[key] || defaultValue;
}

function getEnvAsInt(key: string, defaultValue: number): number {
  const value = process.env[key];
  if (value) {
    const parsed = parseInt(value, 10);
    if (!isNaN(parsed)) return parsed;
  }
  return defaultValue;
}

function getEnvAsFloat(key: string, defaultValue: number): number {
  const value = process.env[key];
  if (value) {
    const parsed = parseFloat(value);
    if (!isNaN(parsed)) return parsed;
  }
  return defaultValue;
}

function getEnvAsSlice(key: string, defaultValue: string[] = []): string[] {
  const value = process.env[key];
  if (!value) return defaultValue;

  const values = value
    .split(',')
    .map(v => v.trim())
    .filter(v => v.length > 0);

  return values.length > 0 ? values : defaultValue;
}

function getEnvAsBool(key: string, defaultValue: boolean): boolean {
  const value = process.env[key]?.toLowerCase();
  if (!value) return defaultValue;

  return value === 'true' || value === '1' || value === 'yes';
}

export function loadConfig(): Config {
  const config: Config = {
    // Server
    port: getEnv('PORT', '8080'),
    env: getEnv('ENV', 'development'),

    // Database
    databaseUrl: getEnv(
      'DATABASE_URL',
      'postgres://postgres:postgres@localhost:5432/mylittleprice?sslmode=disable'
    ),

    // Redis
    redisUrl: getEnv('REDIS_URL', 'localhost:6379'),
    redisPassword: getEnv('REDIS_PASSWORD', ''),
    redisDb: getEnvAsInt('REDIS_DB', 0),

    // JWT
    jwtAccessSecret: getEnv('JWT_ACCESS_SECRET', 'change-me-in-production-access-secret-key'),
    jwtRefreshSecret: getEnv('JWT_REFRESH_SECRET', 'change-me-in-production-refresh-secret-key'),
    jwtAccessTTL: getEnvAsInt('JWT_ACCESS_TTL', 900), // 15 minutes
    jwtRefreshTTL: getEnvAsInt('JWT_REFRESH_TTL', 604800), // 7 days

    // Google OAuth
    googleClientId: getEnv('GOOGLE_CLIENT_ID', ''),
    googleClientSecret: getEnv('GOOGLE_CLIENT_SECRET', ''),
    googleRedirectUrl: getEnv('GOOGLE_REDIRECT_URL', ''),

    // Session
    sessionTTL: getEnvAsInt('SESSION_TTL', 86400),
    maxMessagesPerSession: getEnvAsInt('MAX_MESSAGES_PER_SESSION', 8),
    maxSearchesPerSession: getEnvAsInt('MAX_SEARCHES_PER_SESSION', 3),

    // API Keys
    geminiApiKeys: getEnvAsSlice('GEMINI_API_KEYS', []),
    serpApiKeys: getEnvAsSlice('SERP_API_KEYS', []),

    // Gemini Configuration
    geminiModel: getEnv('GEMINI_MODEL', 'gemini-2.0-flash-exp'),
    geminiFallbackModel: getEnv('GEMINI_FALLBACK_MODEL', 'gemini-1.5-flash'),
    geminiTemperature: getEnvAsFloat('GEMINI_TEMPERATURE', 0.7),
    geminiMaxOutputTokens: getEnvAsInt('GEMINI_MAX_OUTPUT_TOKENS', 8192),
    geminiUseGrounding: getEnvAsBool('GEMINI_USE_GROUNDING', true),

    // Grounding Strategy
    geminiGroundingMode: getEnv('GEMINI_GROUNDING_MODE', 'balanced') as any,
    geminiGroundingMinWords: getEnvAsInt('GEMINI_GROUNDING_MIN_WORDS', 2),

    // Grounding Strategy Thresholds
    geminiBrandQueryConfidence: getEnvAsFloat('GEMINI_BRAND_QUERY_CONFIDENCE', 0.95),
    geminiGroundingWeightFreshInfo: getEnvAsFloat('GEMINI_GROUNDING_WEIGHT_FRESH_INFO', 0.3),
    geminiGroundingWeightSpecific: getEnvAsFloat('GEMINI_GROUNDING_WEIGHT_SPECIFIC_PRODUCT', 0.35),
    geminiGroundingWeightDrift: getEnvAsFloat('GEMINI_GROUNDING_WEIGHT_DIALOGUE_DRIFT', 0.2),
    geminiGroundingWeightElectron: getEnvAsFloat('GEMINI_GROUNDING_WEIGHT_ELECTRONICS', 0.15),
    geminiGroundingDecisionThresh: getEnvAsFloat('GEMINI_GROUNDING_DECISION_THRESHOLD', 0.5),
    geminiBrandQueryMaxWords: getEnvAsInt('GEMINI_BRAND_QUERY_MAX_WORDS', 3),
    geminiBrandSimilarityThresh: getEnvAsFloat('GEMINI_BRAND_SIMILARITY_THRESHOLD', 0.65),
    geminiDialogueHistoryWindow: getEnvAsInt('GEMINI_DIALOGUE_HISTORY_WINDOW', 4),
    geminiDialogueDriftThresh: getEnvAsFloat('GEMINI_DIALOGUE_DRIFT_THRESHOLD', 0.4),
    geminiDriftScoreBonus: getEnvAsFloat('GEMINI_DRIFT_SCORE_BONUS', 0.8),
    geminiElectronicsThreshHigh: getEnvAsFloat('GEMINI_ELECTRONICS_THRESHOLD_HIGH', 0.7),
    geminiElectronicsScoreHigh: getEnvAsFloat('GEMINI_ELECTRONICS_SCORE_HIGH', 0.9),
    geminiCategorySimilarityThresh: getEnvAsFloat('GEMINI_CATEGORY_SIMILARITY_THRESHOLD', 0.6),
    geminiCategoryScore: getEnvAsFloat('GEMINI_CATEGORY_SCORE', 0.5),

    // Gemini Translation
    geminiTranslationTemperature: getEnvAsFloat('GEMINI_TRANSLATION_TEMPERATURE', 0.3),
    geminiTranslationMaxTokens: getEnvAsInt('GEMINI_TRANSLATION_MAX_TOKENS', 100),

    // Embedding
    geminiEmbeddingModel: getEnv('GEMINI_EMBEDDING_MODEL', 'text-embedding-004'),
    embeddingCategoryDetectionThresh: getEnvAsFloat('EMBEDDING_CATEGORY_DETECTION_THRESHOLD', 0.6),
    cacheQueryEmbeddingTTL: getEnvAsInt('CACHE_QUERY_EMBEDDING_TTL', 86400),

    // SERP Thresholds
    serpThresholdExact: getEnvAsFloat('SERP_THRESHOLD_EXACT', 0.4),
    serpThresholdParameters: getEnvAsFloat('SERP_THRESHOLD_PARAMETERS', 0.2),
    serpThresholdCategory: getEnvAsFloat('SERP_THRESHOLD_CATEGORY', 0.1),
    serpLogTopResultsCount: getEnvAsInt('SERP_LOG_TOP_RESULTS_COUNT', 5),
    serpFallbackMinResults: getEnvAsInt('SERP_FALLBACK_MIN_RESULTS', 3),

    // SERP Scoring
    serpScorePhraseMatch: getEnvAsFloat('SERP_SCORE_PHRASE_MATCH', 1.0),
    serpScoreAllWords: getEnvAsFloat('SERP_SCORE_ALL_WORDS', 0.6),
    serpScorePartialWords: getEnvAsFloat('SERP_SCORE_PARTIAL_WORDS', 0.5),
    serpScoreWordOrderWeight: getEnvAsFloat('SERP_SCORE_WORD_ORDER_WEIGHT', 0.2),
    serpScoreBrandMatch: getEnvAsFloat('SERP_SCORE_BRAND_MATCH', 0.3),
    serpScoreModelMatch: getEnvAsFloat('SERP_SCORE_MODEL_MATCH', 0.3),
    serpMinWordLength: getEnvAsInt('SERP_MIN_WORD_LENGTH', 2),
    serpModelNumberMinLength: getEnvAsInt('SERP_MODEL_NUMBER_MIN_LENGTH', 2),

    // SERP Max Products
    serpMaxProductsExact: getEnvAsInt('SERP_MAX_PRODUCTS_EXACT', 3),
    serpMaxProductsParameters: getEnvAsInt('SERP_MAX_PRODUCTS_PARAMETERS', 6),
    serpMaxProductsCategory: getEnvAsInt('SERP_MAX_PRODUCTS_CATEGORY', 8),
    serpMaxProductsDefault: getEnvAsInt('SERP_MAX_PRODUCTS_DEFAULT', 6),

    // Defaults
    defaultCountry: getEnv('DEFAULT_COUNTRY', 'CH'),
    defaultLanguage: getEnv('DEFAULT_LANGUAGE', 'en'),
    defaultCurrency: getEnv('DEFAULT_CURRENCY', 'CHF'),

    // Cache
    cacheGeminiTTL: getEnvAsInt('CACHE_GEMINI_TTL', 3600),
    cacheSerpTTL: getEnvAsInt('CACHE_SERP_TTL', 86400),
    cacheImmersiveTTL: getEnvAsInt('CACHE_IMMERSIVE_TTL', 43200),

    // Rate Limiting
    rateLimitRequests: getEnvAsInt('RATE_LIMIT_REQUESTS', 100),
    rateLimitWindow: getEnvAsInt('RATE_LIMIT_WINDOW', 60),

    // CORS
    corsOrigins: getEnvAsSlice('CORS_ORIGINS', ['http://localhost:3000']),

    // Logging
    logLevel: getEnv('LOG_LEVEL', 'info'),
  };

  validateConfig(config);
  return config;
}

function validateConfig(config: Config): void {
  const errors: string[] = [];

  if (config.geminiApiKeys.length === 0) {
    errors.push('At least one GEMINI_API_KEY is required');
  }

  if (config.serpApiKeys.length === 0) {
    errors.push('At least one SERP_API_KEY is required');
  }

  if (!config.googleClientId) {
    errors.push('GOOGLE_CLIENT_ID is required');
  }

  if (!config.googleClientSecret) {
    errors.push('GOOGLE_CLIENT_SECRET is required');
  }

  const validModes = ['conservative', 'balanced', 'aggressive'];
  if (!validModes.includes(config.geminiGroundingMode)) {
    errors.push(`GEMINI_GROUNDING_MODE must be one of: ${validModes.join(', ')}`);
  }

  if (config.geminiGroundingMinWords < 1 || config.geminiGroundingMinWords > 10) {
    errors.push('GEMINI_GROUNDING_MIN_WORDS must be between 1 and 10');
  }

  if (config.maxSearchesPerSession < 1 || config.maxSearchesPerSession > 10) {
    errors.push('MAX_SEARCHES_PER_SESSION must be between 1 and 10');
  }

  if (errors.length > 0) {
    throw new Error(`Configuration validation failed:\n${errors.join('\n')}`);
  }
}

// Singleton instance
let configInstance: Config | null = null;

export function getConfig(): Config {
  if (!configInstance) {
    configInstance = loadConfig();
  }
  return configInstance;
}
