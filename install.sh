#!/bin/bash

# Sumb Installation Script

set -e

echo "🚀 Installing Sumb Task Management Application..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [[ "$(echo -e "1.21\n$GO_VERSION" | sort -V | head -n1)" != "1.21" ]]; then
    echo "❌ Go version 1.21 or later is required. Current version: $GO_VERSION"
    exit 1
fi

echo "✅ Go version check passed: $GO_VERSION"

# Build the application
echo "🔨 Building sumb..."
make build

# Install to system
if [[ "$EUID" -eq 0 ]]; then
    echo "📦 Installing to /usr/local/bin..."
    cp sumb /usr/local/bin/
    echo "✅ Installation complete! You can now use 'sumb' command."
else
    echo "📦 Installing to ~/.local/bin..."
    mkdir -p ~/.local/bin
    cp sumb ~/.local/bin/
    
    # Check if ~/.local/bin is in PATH
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo "⚠️  ~/.local/bin is not in your PATH."
        echo "Add this line to your shell profile (.bashrc, .zshrc, etc.):"
        echo "export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
    
    echo "✅ Installation complete! You can now use 'sumb' command."
fi

echo ""
echo "🎉 Sumb is now installed!"
echo ""
echo "Usage examples:"
echo "  sumb create -t \"Buy groceries\" -d \"Milk, bread, eggs\""
echo "  sumb list"
echo "  sumb complete 1"
echo "  sumb delete 1"
echo "  sumb --help"
echo ""
echo "For more information, visit: https://github.com/yourusername/sumb" 