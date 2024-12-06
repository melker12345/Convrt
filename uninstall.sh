#!/bin/bash

# Exit on error
set -e

# Define colors for output
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Define installation directory
INSTALL_DIR="$HOME/.local/bin"

echo "Uninstalling convrt..."

# Remove executable
if [ -f "$INSTALL_DIR/convrt" ]; then
    rm "$INSTALL_DIR/convrt"
    echo "Removed convrt executable"
fi

# Clean up PATH in shell config files
if [ -f "$HOME/.bashrc" ]; then
    sed -i '/export PATH="$HOME\/.local\/bin:$PATH"/d' "$HOME/.bashrc"
fi

if [ -f "$HOME/.zshrc" ]; then
    sed -i '/export PATH="$HOME\/.local\/bin:$PATH"/d' "$HOME/.zshrc"
fi

echo -e "${GREEN}Uninstallation complete!${NC}"
echo "Please restart your terminal for the changes to take effect."
