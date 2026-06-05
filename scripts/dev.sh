#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"

if [ ! -d "$ROOT/ui/node_modules" ]; then
	echo "Installing UI dependencies..."
	(cd "$ROOT/ui" && pnpm install)
fi

echo "Starting UI dev server (Vite on http://127.0.0.1:5173)..."
(cd "$ROOT/ui" && pnpm run dev) &
UI_PID=$!

cleanup() {
	echo ""
	echo "Shutting down dev servers..."
	kill "$UI_PID" 2>/dev/null || true
	wait "$UI_PID" 2>/dev/null || true
}

trap cleanup INT TERM EXIT

echo "Starting API (air with --dev proxy)..."
echo "Browse at http://0.0.0.0:9081 (API + proxied UI)"
cd "$ROOT"
air
