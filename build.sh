#!/bin/bash

APP_NAME="dp-cli"

TARGETS=(
  "windows/amd64"
  "windows/386"
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
)

OUTPUT_DIR="build"
mkdir -p "$OUTPUT_DIR"

for TARGET in "${TARGETS[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "$TARGET"
  EXT=""
  if [ "$GOOS" = "windows" ]; then
    EXT=".exe"
  fi
  OUTPUT_NAME="${APP_NAME}_${GOOS}_${GOARCH}${EXT}"
  echo "Building for $GOOS/$GOARCH..."
  env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$OUTPUT_DIR/$OUTPUT_NAME" main.go
done

echo "Done! Binaries are in ./$OUTPUT_DIR"
