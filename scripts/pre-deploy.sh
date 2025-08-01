#!/bin/bash

# Pre-deploy script for threatreg
# Installs Atlas and runs database migrations

set -euo pipefail  # Exit on any error, undefined variables, and pipe failures

echo "🚀 Starting pre-deploy setup..."

# Install curl if not already installed
if ! command -v curl &> /dev/null; then
    echo "📦 Installing curl..."
    if command -v apt-get &> /dev/null; then
        apt-get update && apt-get install -y curl
    elif command -v yum &> /dev/null; then
        yum install -y curl
    elif command -v apk &> /dev/null; then
        apk add --no-cache curl
    elif command -v brew &> /dev/null; then
        brew install curl
    else
        echo "❌ Error: Cannot install curl. Please install curl manually."
        exit 1
    fi
    echo "✅ curl installed successfully"
else
    echo "✅ curl already installed"
fi

# Install Atlas if not already installed
if ! command -v atlas &> /dev/null; then
    echo "📦 Installing Atlas..."
    curl -sSf https://atlasgo.sh | sh
    echo "✅ Atlas installed successfully"
else
    echo "✅ Atlas already installed"
fi

# Verify Atlas installation
atlas version

# Run database migrations
echo "🔄 Running database migrations..."
atlas migrate apply --env postgres

echo "✅ Pre-deploy setup completed successfully!"