#!/bin/bash


# Detect WSL and Host IP
if grep -q "microsoft" /proc/version; then
    echo "WSL detected. Resolving Host IP..."
    HOST_IP=$(grep nameserver /etc/resolv.conf | awk '{print $2}')
    if [ -n "$HOST_IP" ]; then
        echo "Host IP: $HOST_IP"
    else
        echo "Could not determine Host IP, crashing back to localhost"
        HOST_IP="localhost"
    fi
else
    HOST_IP="localhost"
fi

# Define constants
TUNNEL_NAME="brobar-dev"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONFIG="$SCRIPT_DIR/cloudflared.yml"
CRED_DIR="$HOME/.cloudflared"

# Function to check command existence
check_command() {
    if ! command -v "$1" &> /dev/null; then
        echo "Error: $1 is not installed."
        if [ "$1" == "cloudflared" ]; then
            echo "Please install it using: brew install cloudflared" # Or instructions for linux
        fi
        exit 1
    fi
}

check_command cloudflared
check_command awk
check_command sed
check_command grep

echo "Starting Cloudflare Tunnel setup..."

# 1. Check for Login (cert.pem)
if [ ! -f "$CRED_DIR/cert.pem" ]; then
    echo "You are not logged in to Cloudflare."
    echo "Opening login page..."
    cloudflared tunnel login
    if [ ! -f "$CRED_DIR/cert.pem" ]; then
        echo "Login failed or cancelled. Exiting."
        exit 1
    fi
fi

# 2. Check Config and Tunnel ID
if [ ! -f "$CONFIG" ]; then
    echo "Error: Config file not found at $CONFIG"
    exit 1
fi

# Extract current Tunnel ID from config
CURRENT_UUID=$(grep '^tunnel:' "$CONFIG" | awk '{print $2}' | tr -d '"' | tr -d "'")

NEEDS_RECREATION=false

if [ -z "$CURRENT_UUID" ]; then
    echo "No tunnel ID found in config."
    NEEDS_RECREATION=true
else
    # Check if credentials exist for this UUID
    if [ ! -f "$CRED_DIR/$CURRENT_UUID.json" ]; then
        echo "Credentials for tunnel $CURRENT_UUID not found."
        NEEDS_RECREATION=true
    fi
fi

# 3. Create/Recreate Tunnel if needed
if [ "$NEEDS_RECREATION" = true ]; then
    echo "--------------------------------------------------------"
    echo "Tunnel credentials are missing or config is invalid."
    read -p "Do you want to automatically create a NEW tunnel? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborting. Please restore credentials manually."
        exit 1
    fi

    echo "Creating new tunnel '$TUNNEL_NAME'..."
    
    # Clean up old tunnel with same name if exists locally (ignore errors)
    cloudflared tunnel delete -f "$TUNNEL_NAME" &> /dev/null
    
    # Create new tunnel
    # Output format involves a line saying "Created tunnel <NAME> with id <UUID>"
    CREATE_OUTPUT=$(cloudflared tunnel create "$TUNNEL_NAME" 2>&1)
    if [ $? -ne 0 ]; then
        echo "Failed to create tunnel. Output:"
        echo "$CREATE_OUTPUT"
        exit 1
    fi
    
    echo "Tunnel created successfully."
    
    # Extract new UUID from the JSON file that was created in .cloudflared
    # We find the most recent .json file or parse output. 
    # Reliability: cloudflared usually saves as ~/.cloudflared/<UUID>.json
    # parsing UUID from output is safer.
    # Output example: "Created tunnel brobar-dev with id 12345-..."
    NEW_UUID=$(echo "$CREATE_OUTPUT" | grep "Created tunnel" | grep -oE '[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}')
    
    if [ -z "$NEW_UUID" ]; then
        echo "Could not detect new Tunnel UUID from output."
        exit 1
    fi
    
    echo "New Tunnel ID: $NEW_UUID"
    
    # Update config file
    # Replace the tunnel: line
    sed -i "s/^tunnel: .*/tunnel: $NEW_UUID/" "$CONFIG"
    echo "Updated $CONFIG with new UUID."
fi

# ALWAYS update ingress to point to the correct Host IP (for WSL)
if [ "$HOST_IP" != "localhost" ]; then
    echo "Updating ingress to use Host IP: $HOST_IP"
    # Replace localhost with HOST_IP in http services
    sed -i "s|service: http://localhost|service: http://$HOST_IP|g" "$CONFIG"
    # Also handle previous runs where it might be an old IP (regex match ip)
    sed -i "s|service: http://[0-9]\+\.[0-9]\+\.[0-9]\+\.[0-9]\+|service: http://$HOST_IP|g" "$CONFIG"
fi
    
# 4. Setup DNS Routes
    echo "Configuring DNS routes..."
    # Extract hostnames from ingress in yaml
    HOSTNAMES=$(grep -E "^\s*-\s*hostname:" "$CONFIG" | awk '{print $3}')
    
    for HOST in $HOSTNAMES; do
        echo "Routing $HOST -> $NEW_UUID"
        cloudflared tunnel route dns "$TUNNEL_NAME" "$HOST"
    done
    
    echo "Setup complete!"


# 5. Run the tunnel
echo "Starting tunnel..."
# exec replaces the shell process, so Ctrl+C kills the tunnel properly
exec cloudflared tunnel --config "$CONFIG" run
