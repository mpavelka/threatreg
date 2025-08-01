#!/bin/bash

# Pre-deploy script for threatreg
# Installs Atlas and runs database migrations

set -e  # Exit on any error

echo "🚀 Starting pre-deploy setup..."

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