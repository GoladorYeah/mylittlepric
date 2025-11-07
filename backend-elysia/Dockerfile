# MyLittlePrice Backend - Elysia.js Edition
FROM oven/bun:1 AS base
WORKDIR /app

# Install dependencies
FROM base AS install
COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile --production

# Copy source code
FROM base AS release
COPY --from=install /app/node_modules ./node_modules
COPY . .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD bun run --bun fetch http://localhost:8080/api/health || exit 1

# Set environment
ENV NODE_ENV=production

# Run the application
CMD ["bun", "run", "src/main.ts"]
