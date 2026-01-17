import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    reactStrictMode: false,
    output: "standalone",
    images: {
        domains: ["localhost"],
    },
};

export default nextConfig;