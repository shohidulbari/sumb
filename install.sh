#!/usr/bin/env bash

set -e

REPO="shohidulbari/sumb"
VERSION="v1.0.0-alpha" # Pre-release tag
BINARY_NAME="sumb"
INSTALL_DIR="/usr/local/bin"
DRY_RUN=false

# --- Parse args ---
for arg in "$@"; do
  case $arg in
  --dry-run)
    DRY_RUN=true
    shift
    ;;
  --version=*)
    VERSION="${arg#*=}"
    shift
    ;;
  *) ;;
  esac
done

# --- Detect OS (allow override via env var) ---
if [ -z "$OS" ]; then
  OS=$(uname -s)
fi

case "$OS" in
Linux* | linux) OS="linux" ;;
Darwin* | darwin) OS="darwin" ;;
Windows* | windows | MINGW* | MSYS*) OS="windows" ;;
*)
  echo "❌ Unsupported OS: $OS"
  exit 1
  ;;
esac

# --- Detect architecture (allow override via env var) ---
if [ -z "$ARCH" ]; then
  ARCH=$(uname -m)
fi

case "$ARCH" in
x86_64 | amd64) ARCH="amd64" ;;
arm64 | aarch64) ARCH="arm64" ;;
*)
  echo "❌ Unsupported architecture: $ARCH"
  exit 1
  ;;
esac
ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
  ASSET_NAME="${ASSET_NAME}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET_NAME}"

echo "🧠 Detected system:"
echo "   OS:          $OS"
echo "   Architecture: $ARCH"
echo "   Version:      $VERSION"
echo "   Binary:       $ASSET_NAME"
echo "   Download URL: $DOWNLOAD_URL"
echo

if [ "$DRY_RUN" = true ]; then
  echo "✅ Dry run complete — no files downloaded or installed."
  exit 0
fi

echo "⬇️  Downloading ${ASSET_NAME}..."
curl -L "$DOWNLOAD_URL" -o "$ASSET_NAME"

echo "🚀 Installing ${BINARY_NAME} to ${INSTALL_DIR}..."
chmod +x "$ASSET_NAME"
sudo mv "$ASSET_NAME" "$INSTALL_DIR/$BINARY_NAME"

echo "✅ Installation complete!"
echo "Run with: $BINARY_NAME --help"
