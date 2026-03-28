#!/usr/bin/env bash
set -euo pipefail

# COVERAGE_THRESHOLD can be provided as env var (default 80)
THRESHOLD=${COVERAGE_THRESHOLD:-80}

echo "Running tests and generating coverage... (threshold=${THRESHOLD}%)"
if ! go test ./... -coverprofile=coverage.out; then
  echo "go test failed" >&2
  exit 1
fi

if [ ! -f coverage.out ]; then
  echo "cannot determine coverage: coverage.out missing" >&2
  exit 2
fi

PERCENT=$(go tool cover -func=coverage.out | awk '/total:/ {gsub(/%/,"",$3); print $3; exit}') || true
if [ -z "${PERCENT}" ]; then
  echo "cannot determine coverage" >&2
  exit 2
fi

printf "total coverage: %s%%\n" "${PERCENT}"

awk -v cov="${PERCENT}" -v thr="${THRESHOLD}" 'BEGIN { if (cov+0 < thr+0) { printf("coverage too low: %s%% (threshold %s%%)\n", cov, thr); exit 1 } else { printf("coverage OK: %s%% (threshold %s%%)\n", cov, thr); exit 0 } }'
