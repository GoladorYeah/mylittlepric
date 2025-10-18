import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Enable standalone output for Docker
  output: 'standalone',

  // WebSocket configuration for API communication
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: process.env.NEXT_PUBLIC_API_URL
          ? `${process.env.NEXT_PUBLIC_API_URL}/api/:path*`
          : 'http://localhost:8080/api/:path*',
      },
      {
        source: '/ws',
        destination: process.env.NEXT_PUBLIC_API_URL
          ? `${process.env.NEXT_PUBLIC_API_URL}/ws`
          : 'http://localhost:8080/ws',
      },
    ];
  },
};

export default nextConfig;
