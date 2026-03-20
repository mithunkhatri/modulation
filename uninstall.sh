#!/bin/bash

# Modulation Uninstaller 🎵

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting Modulation Uninstallation...${NC}"

# Detect OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS_TYPE=Linux;;
    Darwin*)    OS_TYPE=macOS;;
    MSYS*|MINGW*|CYGWIN*) OS_TYPE=Windows;;
    *)          OS_TYPE="UNKNOWN:${OS}"
esac

# 1. Remove binary
if [ "$OS_TYPE" == "Windows" ]; then
    if [ -f "modulation.exe" ]; then
        rm "modulation.exe"
        echo -e "${GREEN}Removed modulation.exe from current directory.${NC}"
    else
        echo -e "${BLUE}modulation.exe not found in current directory.${NC}"
    fi
else
    BINARY_PATH="/usr/local/bin/modulation"
    if [ -f "$BINARY_PATH" ]; then
        if [ -w "/usr/local/bin" ]; then
            rm "$BINARY_PATH"
        else
            echo "Requesting sudo permissions to remove $BINARY_PATH..."
            sudo rm "$BINARY_PATH"
        fi
        echo -e "${GREEN}Removed $BINARY_PATH successfully.${NC}"
    else
        echo -e "${BLUE}modulation not found in /usr/local/bin.${NC}"
    fi
fi

# 2. Handle configuration files
CONFIG_DIR=""
if [ "$OS_TYPE" == "macOS" ]; then
    CONFIG_DIR="$HOME/Library/Application Support/modulation"
elif [ "$OS_TYPE" == "Linux" ]; then
    CONFIG_DIR="$HOME/.config/modulation"
fi

if [ -n "$CONFIG_DIR" ] && [ -d "$CONFIG_DIR" ]; then
    echo -e "${BLUE}Configuration found at: $CONFIG_DIR${NC}"
    read -p "Do you want to remove all configuration and favorites? [y/N] " confirm
    if [[ "$confirm" =~ ^[Yy]$ ]]; then
        rm -rf "$CONFIG_DIR"
        echo -e "${GREEN}Removed configuration directory.${NC}"
    else
        echo -e "${BLUE}Skipped removing configuration directory.${NC}"
    fi
fi

echo -e "${GREEN}Uninstallation complete!${NC}"
