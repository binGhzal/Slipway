# Slipway - Multi-Cloud K8s Platform

CLUSTER_NAME := slipway-mgmt
KIND_CONFIG := clusters/management/kind-config.yaml

.PHONY: help bootstrap up down clean deps

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Check dependencies
	@command -v kind >/dev/null 2>&1 || { echo >&2 "kind is not installed. Aborting."; exit 1; }
	@command -v kubectl >/dev/null 2>&1 || { echo >&2 "kubectl is not installed. Aborting."; exit 1; }
	@command -v clusterctl >/dev/null 2>&1 || { echo >&2 "clusterctl is not installed. Aborting."; exit 1; }
	@echo "All dependencies check passed."

up: ## Spin up local Kind cluster (without CNI/KubeProxy)
	kind create cluster --config $(KIND_CONFIG)

down: ## Destroy local cluster
	kind delete cluster --name $(CLUSTER_NAME)

bootstrap: deps up ## Bootstrap the entire platform (Kind + CAPI + ArgoCD)
	@echo "Bootstrapping Slipway..."
	# 1. Initialize CAPI with Docker provider
	clusterctl init --infrastructure docker

	# 2. Install ArgoCD
	kubectl create namespace argocd || true
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
	@echo "Waiting for ArgoCD server..."
	kubectl wait --for=condition=available deployment/argocd-server -n argocd --timeout=300s

	# 3. Apply Root App
	kubectl apply -f argocd/bootstrap/root.yaml

	@echo "Bootstrap complete!"
	@echo "IMPORTANT: Push your changes to GitHub so ArgoCD can sync/install Cilium."
	@echo "ArgoCD UI: https://localhost:8080 (Forward port: kubectl port-forward svc/argocd-server -n argocd 8080:443)"
	@echo "Password: kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d"

clean: down ## Clean up everything
