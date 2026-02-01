#!/bin/bash

# Check if cloudflared is installed
if ! command -v cloudflared &> /dev/null; then
    echo "cloudflared is not installed. Please install it first:"
    echo "brew install cloudflared"
    exit 1
fi

echo "Starting Cloudflare Tunnel..."
echo "Config: scripts/cloudflared.yml"
echo "Tunnel: brobar-dev (5cf7266c-daf2-47a1-a6ac-91abecd85c26)"

# Run the tunnel using the config file
cloudflared tunnel --config scripts/cloudflared.yml run
