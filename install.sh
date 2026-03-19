#!/bin/bash

# Modulation Installer 🎵

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting Modulation Installation (v1.0.1)...${NC}"

# If go.mod is missing, we are likely running via curl/wget one-liner
IS_REMOTE=false
if [ ! -f "go.mod" ]; then
    IS_REMOTE=true
    echo -e "${BLUE}Project files not found locally. Attempting to fetch from GitHub...${NC}"
    if ! command -v git &> /dev/null; then
        echo -e "${RED}Error: git is required for the one-line installer.${NC}"
        echo "Please install git or download the source manually."
        exit 1
    fi
    TMP_DIR=$(mktemp -d)
    echo -e "${BLUE}Cloning into temporary directory: $TMP_DIR${NC}"
    if ! git clone https://github.com/mithunkhatri/modulation.git "$TMP_DIR"; then
        echo -e "${RED}Error: Failed to clone repository.${NC}"
        exit 1
    fi
    cd "$TMP_DIR" || exit 1

    # Checkout latest release tag if it exists
    LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [ -n "$LATEST_TAG" ]; then
        echo -e "${BLUE}Checking out latest release: $LATEST_TAG${NC}"
        git checkout "$LATEST_TAG" -q
    fi
fi

if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found after clone. Please check the repository URL.${NC}"
    exit 1
fi

OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS_TYPE=Linux;;
    Darwin*)    OS_TYPE=macOS;;
    MSYS*|MINGW*|CYGWIN*) OS_TYPE=Windows;;
    *)          OS_TYPE="UNKNOWN:${OS}"
esac

echo -e "Detected OS: ${BLUE}${OS_TYPE}${NC}"

# Define binary name
BINARY_NAME="modulation"
if [ "$OS_TYPE" == "Windows" ]; then
    BINARY_NAME="modulation.exe"
fi

# 1. Check/Install Go
if ! command -v go &> /dev/null; then
    echo -e "${BLUE}Go not found. Attempting to install...${NC}"
    if [ "$OS_TYPE" == "macOS" ] && command -v brew &> /dev/null; then
        brew install go
    elif [ "$OS_TYPE" == "Windows" ]; then
        if command -v choco &> /dev/null; then
            choco install golang -y
        elif command -v scoop &> /dev/null; then
            scoop install go
        else
            echo -e "${RED}Error: Go not found and no supported package manager (Choco/Scoop) found.${NC}"
            exit 1
        fi
    elif [ "$(command -v apt-get)" ]; then
        sudo apt-get update && sudo apt-get install -y golang-go
    elif [ "$(command -v dnf)" ]; then
        sudo dnf install -y golang
    else
        echo -e "${RED}Error: Go is not installed and no supported package manager found.${NC}"
        echo "Please install Go from https://golang.org/dl/"
        exit 1
    fi
fi

# 2. Check/Install FFmpeg
if ! command -v ffmpeg &> /dev/null; then
    echo -e "${BLUE}FFmpeg not found. Attempting to install...${NC}"
    if [ "$OS_TYPE" == "macOS" ] && command -v brew &> /dev/null; then
        brew install ffmpeg
    elif [ "$OS_TYPE" == "Windows" ]; then
        if command -v choco &> /dev/null; then
            choco install ffmpeg -y
        elif command -v scoop &> /dev/null; then
            scoop install ffmpeg
        else
            echo -e "${RED}Error: FFmpeg not found and no supported package manager (Choco/Scoop) found.${NC}"
            exit 1
        fi
    elif [ "$(command -v apt-get)" ]; then
        sudo apt-get update && sudo apt-get install -y ffmpeg
    elif [ "$(command -v dnf)" ]; then
        sudo dnf install -y ffmpeg
    else
        echo -e "${RED}Error: FFmpeg is not installed and no supported package manager found.${NC}"
        exit 1
    fi
fi

# 3. Build the application
echo -e "${BLUE}Building Modulation...${NC}"
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
go build -ldflags "-X main.version=$VERSION" -o "$BINARY_NAME" .

# 4. Install
if [ "$OS_TYPE" == "Windows" ]; then
    echo -e "${GREEN}Build successful!${NC}"
    echo -e "Run the app with: ${BLUE}./$BINARY_NAME${NC}"
else
    echo -e "${BLUE}Installing to /usr/local/bin...${NC}"
    if [ -w /usr/local/bin ]; then
        cp "$BINARY_NAME" /usr/local/bin/modulation
    else
        echo "Requesting sudo permissions to copy to /usr/local/bin..."
        sudo cp "$BINARY_NAME" /usr/local/bin/modulation
    fi
    echo -e "${GREEN}Installation successful!${NC}"
    echo -e "You can now run the app by typing: ${BLUE}modulation${NC}"
fi
echo -e "Made with ☕ by Mithun Khatri"
