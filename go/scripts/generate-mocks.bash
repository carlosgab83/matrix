#!/bin/bash

# Script to generate mocks with Mockery
# All mocks centralized in internal/shared/mocks/

set -e  # Exit on error

echo "🔨 Generating mocks with Mockery (centralized in shared/mocks)..."
echo ""

# Check that .mockery.yaml exists
if [ ! -f ".mockery.yaml" ]; then
    echo "❌ Error: .mockery.yaml not found"
    echo "   Run this script from the project root directory (go/)"
    exit 1
fi

# Clean previous mocks
MOCKS_DIR="internal/shared/mocks"
echo "🧹 Cleaning previous mocks..."
rm -rf "$MOCKS_DIR"
mkdir -p "$MOCKS_DIR"
echo ""

# Generate mocks with Docker (using our matrix:latest image)
echo "📦 Running Mockery in matrix:latest container..."
echo "🔍 Scanning internal/ for interfaces..."
echo ""

docker run --rm \
    -v "$PWD":/app \
    -w /app \
    matrix:latest \
    mockery --config .mockery.yaml

# If you have mockery installed locally, comment the line above and use:
# mockery --config .mockery.yaml

echo ""
echo "════════════════════════════════════════════════════"
echo "✅ Generation completed!"
echo ""
echo "📂 Mocks generated in: $MOCKS_DIR/"
if [ -d "$MOCKS_DIR" ]; then
    mock_count=$(find "$MOCKS_DIR" -name "*.go" 2>/dev/null | wc -l)
    echo "📊 Total mocks: $mock_count"
    echo ""
    echo "📄 Generated files:"
    ls -1 "$MOCKS_DIR" 2>/dev/null | sed 's/^/   /' | sort
else
    echo "⚠️  No mocks generated (no interfaces found)"
fi
echo ""
echo "💡 Usage in tests:"
echo "   import \"github.com/carlosgab83/matrix/go/internal/shared/mocks\""
echo ""
echo "📝 Examples:"
echo "   mockLogger := mocks.NewLogger(t)"
echo "   mockIngestor := mocks.NewIngestor(t)"
echo "   mockFetcher := mocks.NewSymbolFetcher(t)"
