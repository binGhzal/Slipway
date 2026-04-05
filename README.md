# Slipway

A fully GitOps-managed Kubernetes platform for homelabs and beyond.

Provisions **Talos Linux** clusters on **Proxmox** using **Cluster API**, managed by **ArgoCD**, with **Cilium** for all networking. Designed to be portable, reproducible, and hands-off.

## Features

- **100% GitOps** - Every change flows through Git. ArgoCD syncs, controllers reconcile.
- **Immutable infrastructure** - Talos Linux nodes are API-driven, no SSH, no drift.
- **Declarative clusters** - Define clusters as YAML, CAPI provisions them automatically.
- **Atomic configuration** - Change networking, kubelet args, anything - commit and it rolls out.
- **Automatic DNS** - Create an HTTPRoute, get a Cloudflare subdomain and TLS cert automatically.
- **No manual secrets** - External Secrets Operator syncs from your secret store. No SOPS.
- **Portable** - Clone, fill in `.env`, run `task full-bootstrap`. Works on any Proxmox setup.
- **Multi-cloud ready** - Architecture supports adding cloud CAPI providers without restructuring.

## Quick Start

```bash
# 1. Install prerequisites
# kind, kubectl, clusterctl, talosctl, helm, cilium CLI

# 2. Clone and configure
git clone https://github.com/binGhzal/Slipway.git
cd Slipway
cp .env.example .env
# Edit .env with your Proxmox, Cloudflare, and secret store values

# 3. Create Talos VM template on Proxmox
./scripts/create-talos-template.sh

# 4. Bootstrap everything
task full-bootstrap
```

## Architecture

```
Kind (temp) -> Management Cluster (Talos/Proxmox) -> Workload Clusters
                     |                                      |
                     +-- CAPI controllers                   +-- Apps (homelab)
                     +-- ArgoCD (self-managed)              +-- Platform services
                     +-- Cilium + Hubble                    +-- Cilium networking
                     +-- Observability stack                +-- Storage (Longhorn)
                     +-- External Secrets Operator
                     +-- cert-manager + external-dns
```

## Adding a Workload Cluster

```bash
cp -r infrastructure/clusters/homelab infrastructure/clusters/my-cluster
# Edit patches with your values
git add . && git commit -m "add my-cluster" && git push
# ArgoCD + CAPI handle the rest
```

## Deploying an App

```bash
mkdir -p apps/homelab/my-app
# Add K8s manifests or Helm chart
git add . && git commit -m "deploy my-app" && git push
# ArgoCD deploys it. external-dns creates DNS. cert-manager issues TLS.
```

## Documentation

| Doc | Description |
|-----|-------------|
| [Architecture](docs/architecture.md) | System design and component interactions |
| [Bootstrap Guide](docs/bootstrap.md) | Step-by-step setup instructions |
| [Talos Image](docs/talos-image.md) | Building and managing Talos VM templates |
| [Adding Clusters](docs/adding-cluster.md) | How to provision new workload clusters |
| [Deploying Apps](docs/adding-app.md) | How to deploy applications |
| [Secret Management](docs/secrets.md) | ESO setup and secret flow |
| [Networking](docs/networking.md) | Cilium, DNS, ingress architecture |
| [Troubleshooting](docs/troubleshooting.md) | Common issues and debugging |

## Platform Components

### Core
Cilium (CNI + LB + Gateway API), cert-manager, external-dns, External Secrets Operator, Cloudflare Tunnel, Reloader

### Observability
Prometheus + Grafana, Loki, Tempo, Hubble, Grafana Alloy, OpenCost, Gatus, Pyrra

### Security
Kyverno, Trivy Operator, Kubescape, Tetragon, Cilium Network Policies

### Operations
ArgoCD, CAPI Operator, Velero, Descheduler, Node Problem Detector, Spegel, Headlamp, Goldilocks

### Storage
Proxmox CSI, Longhorn, CloudNativePG
