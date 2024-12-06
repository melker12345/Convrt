#!/bin/bash

# Exit on error
set -e

# Define colors for output
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Define installation directory
INSTALL_DIR="$HOME/.local/bin"

# Create installation directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

echo "Installing convrt..."

# Build for current platform
go build -o convrt

# Copy executable to installation directory
cp convrt "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/convrt"

# Add to PATH if not already present
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zshrc" 2>/dev/null || true
    echo -e "${GREEN}Added convrt to PATH. Please restart your terminal or run:${NC}"
    echo -e "    source ~/.bashrc"
    echo -e "    # or for zsh:"
    echo -e "    source ~/.zshrc"
else
    echo -e "${GREEN}convrt is already in PATH${NC}"
fi

echo -e "${GREEN}Installation complete! You can now use 'convrt' from anywhere.${NC}"
echo "Try running 'convrt --help' to get started."
