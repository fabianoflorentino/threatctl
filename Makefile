
.PHONY: help build test tidy coverage check-coverage check-coverage-strict clean bins

bin = bin

# Colors
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m

.DEFAULT_GOAL := help

##@ Geral

help: ## Mostra esta mensagem de ajuda
	@echo ""
	@echo -e "$(BLUE)╔══════════════════════════════════════════════════════════╗$(NC)"
	@echo -e "$(BLUE)║                  threatctl - Comandos disponíveis         ║$(NC)"
	@echo -e "$(BLUE)╚══════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z0-9_.-]+:.*?##/ { printf "  %-20s %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0,5) } ' $(MAKEFILE_LIST)
	@echo ""

##@ Build & Binaries

build: bins ## Constrói os binários do projeto
	@echo -e "$(BLUE)🔨 Building binaries...$(NC)"
	go build -o $(bin)/threatctl .
	go build -o $(bin)/genpcap ./genpcap
	@echo -e "$(GREEN)✓ Build concluído: $(bin)/$(NC)"

bins: ## Cria diretório de binários
	@mkdir -p $(bin)

clean: ## Remove artefatos gerados (bin, coverage)
	@echo -e "$(YELLOW)🧹 Limpando artefatos...$(NC)"
	@rm -rf $(bin) coverage.out
	@echo -e "$(GREEN)✓ Limpo$(NC)"

##@ Testes e qualidade

tidy: ## Roda `go mod tidy`
	go mod tidy

test: ## Executa todos os testes
	go test ./...

coverage: ## Gera o arquivo coverage.out e mostra o resumo
	go test ./... -coverprofile=coverage.out
	@echo "Coverage summary:"
	go tool cover -func=coverage.out | sed -n '/total:/p'




check-coverage: ## Roda o script em scripts/check_coverage.sh que imprime o resumo
	@bash scripts/check_coverage.sh

check-coverage-strict: ## Roda verificação estrita de cobertura (usa COVERAGE_THRESHOLD)
	@bash scripts/check_coverage_strict.sh || printf "$(YELLOW)WARNING: check-coverage-strict failed (coverage below threshold)$(NC)\n"
