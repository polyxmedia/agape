#!/bin/sh
# Convenience wrapper: load .env then run the agape binary with all args
# passed through. The .env file is gitignored and chmod 600.
#
# Examples:
#   ./run.sh -smoke
#   ./run.sh -config config.yaml
#   ./run.sh                       # uses default config.yaml

set -e

DIR="$(cd "$(dirname "$0")" && pwd)"

if [ ! -f "$DIR/.env" ]; then
    echo "missing $DIR/.env (expected: ANTHROPIC_API_KEY, OPENAI_API_KEY)" >&2
    exit 1
fi

# shellcheck disable=SC1091
. "$DIR/.env"

cd "$DIR"
exec go run ./cmd/agape "$@"
