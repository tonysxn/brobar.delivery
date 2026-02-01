import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    reactStrictMode: false,
    output: "standalone",
    images: {
        domains: ["localhost"],
    },
    eslint: {
        ignoreDuringBuilds: true,
    },
    typescript: {
        ignoreBuildErrors: true,
    },
};

export default nextConfig;