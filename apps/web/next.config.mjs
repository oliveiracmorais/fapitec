/** @type {import('next').NextConfig} */
const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

const nextConfig = {
  typedRoutes: true,
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        destination: `${API_URL}/api/v1/:path*`,
      },
    ];
  },
  output: process.env.NEXT_OUTPUT === "standalone" ? "standalone" : process.env.NEXT_STATIC_OUTPUT === "true" ? "export" : undefined,
};

export default nextConfig;
