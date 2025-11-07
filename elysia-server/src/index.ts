/**
 * MyLittlePrice Backend - Elysia.js Edition
 * Main entry point
 */

import { Elysia } from 'elysia';
import { cors } from '@elysiajs/cors';
import { getConfig } from './config';
import { Container } from './container';
import { chatModule } from './modules/chat';
import { authModule } from './modules/auth';
import { statsModule } from './modules/stats';

// Load configuration
const config = getConfig();

console.log('🚀 MyLittlePrice Backend - Elysia.js Edition');
console.log('='.repeat(50));
console.log(`Environment: ${config.env}`);
console.log(`Port: ${config.port}`);
console.log(`Allowed origins: ${config.corsOrigins.join(', ')}`);
console.log('='.repeat(50));

// Initialize container
const container = new Container(config);

// Create Elysia app
const app = new Elysia()
  // CORS middleware
  .use(
    cors({
      origin: (origin) => {
        // Check if origin is in allowed list
        const allowed = config.corsOrigins.includes(origin);

        if (allowed && config.env === 'development') {
          console.log(`✅ CORS allowed origin: ${origin}`);
        }

        if (!allowed) {
          console.log(
            `❌ CORS REJECTED - Origin: '${origin}' not in allowed list: ${config.corsOrigins.join(', ')}`
          );
          console.log(
            `💡 Fix: Add '${origin}' to CORS_ORIGINS environment variable`
          );
        }

        return allowed;
      },
      methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
      allowedHeaders: ['Origin', 'Content-Type', 'Accept', 'Authorization'],
      credentials: true,
      maxAge: 86400,
    })
  )

  // Request logging
  .onRequest(({ request }) => {
    const timestamp = new Date().toISOString().substring(11, 19);
    console.log(`[${timestamp}] ${request.method} ${new URL(request.url).pathname}`);
  })

  // Error handling
  .onError(({ code, error, set }) => {
    console.error('❌ Error:', error);

    if (code === 'NOT_FOUND') {
      set.status = 404;
      return {
        error: true,
        message: 'Not found',
        code: 404,
      };
    }

    if (code === 'VALIDATION') {
      set.status = 400;
      return {
        error: true,
        message: 'Validation error',
        code: 400,
        details: error.message,
      };
    }

    set.status = 500;
    return {
      error: true,
      message: 'Internal server error',
      code: 500,
    };
  })

  // Register modules
  .use(chatModule(container))
  .use(authModule(container))
  .use(statsModule(container))

  // Root endpoint
  .get('/', () => ({
    app: 'MyLittlePrice API',
    version: '2.0.0',
    runtime: 'Elysia.js',
    status: 'running',
  }))

  // Start server
  .listen(config.port);

console.log('='.repeat(50));
console.log(`🚀 Server running on http://localhost:${config.port}`);
console.log(`🔒 Environment: ${config.env}`);
console.log(`🌍 Allowed origins: ${config.corsOrigins.join(', ')}`);
console.log('='.repeat(50));

// Graceful shutdown
process.on('SIGINT', async () => {
  console.log('\n🛑 Shutting down server...');

  try {
    await container.close();
    app.stop();
    console.log('✅ Server stopped gracefully');
    process.exit(0);
  } catch (error) {
    console.error('❌ Error during shutdown:', error);
    process.exit(1);
  }
});

process.on('SIGTERM', async () => {
  console.log('\n🛑 Shutting down server...');

  try {
    await container.close();
    app.stop();
    console.log('✅ Server stopped gracefully');
    process.exit(0);
  } catch (error) {
    console.error('❌ Error during shutdown:', error);
    process.exit(1);
  }
});
