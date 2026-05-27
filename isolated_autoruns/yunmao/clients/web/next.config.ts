import type { NextConfig } from "next";

const cfg: NextConfig = {
  reactStrictMode: true,
  experimental: {
    typedRoutes: true,
  },
  env: {
    NEXT_PUBLIC_API_BASE: process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:18000",
    NEXT_PUBLIC_WS_BASE: process.env.NEXT_PUBLIC_WS_BASE ?? "ws://localhost:18007",
  },
};

export default cfg;
