#!/usr/bin/env bash
set -euo pipefail

# Gera coverage.out (não falha se os testes falharem)
if ! go test ./... -coverprofile=coverage.out >/dev/null 2>&1; then
  echo "go test returned non-zero or no tests; continuing"
fi

if [ ! -f coverage.out ]; then
  echo "cannot determine coverage: coverage.out missing"
  exit 0
fi

# Imprime o sumário total de coverage
go tool cover -func=coverage.out | sed -n '/total:/p' || true
