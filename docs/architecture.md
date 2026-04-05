# Architecture

## Overview

Slipway is a fully GitOps-managed Kubernetes platform that provisions Talos Linux clusters on Proxmox using Cluster API (CAPI), with ArgoCD for continuous delivery and Cilium for all networking.

## Cluster Topology

```
Kind (temporary)  -->  Management Cluster  -->  Workload Clusters
                       (slipway-mgmt)           (homelab, ...)
```

1. **Kind** - Temporary local cluster used only during initial bootstrap
2. **Management Cluster** - Single-node Talos cluster on Proxmox that runs CAPI controllers and ArgoCD
3. **Workload Clusters** - Multi-node Talos clusters provisioned by CAPI from the management cluster

## Component Stack

### Infrastructure
- **Talos Linux** - Immutable, API-driven Kubernetes OS (no SSH, no shell)
- **Cluster API** - Declarative cluster lifecycle management
  - CAPMOX (Proxmox infrastructure provider)
  - CABPT (Talos bootstrap provider)
  - CACPPT (Talos control plane provider)
  - In-Cluster IPAM

### Networking (Cilium)
- CNI with kube-proxy replacement (eBPF)
- LB-IPAM + L2 Announcements (replaces MetalLB)
- Gateway API (ingress controller)
- Hubble (network observability)
- Tetragon (runtime security)
- WireGuard node-to-node encryption

### GitOps
- **ArgoCD** - App of Apps pattern with ApplicationSets
- **external-dns** - Automatic Cloudflare DNS record management
- **cert-manager** - TLS certificates via Cloudflare DNS-01
- **External Secrets Operator** - Automatic secret sync from secret store

### Observability
- kube-prometheus-stack (Prometheus + Grafana + AlertManager)
- Loki (logs) + Tempo (traces)
- Grafana Alloy (unified agent)
- Hubble (network flows)

### Security
- Kyverno (policy engine)
- Trivy Operator (vulnerability scanning)
- Kubescape (runtime security)
- Cilium Network Policies (L3-L7)

## Data Flow

### Deploying a Change
```
Developer -> Git commit -> ArgoCD detects -> Syncs to cluster -> Controllers reconcile
```

### Provisioning a New Cluster
```
Git (infrastructure/clusters/new/) -> ArgoCD -> CAPI on mgmt -> Proxmox VMs -> Talos bootstrap -> Kyverno registers in ArgoCD -> Platform services deploy
```

### DNS Automation
```
HTTPRoute created -> Cilium creates LB Service -> LB-IPAM assigns IP -> L2 announces -> external-dns creates Cloudflare record -> cert-manager issues TLS cert
```

## Directory Structure

```
Slipway/
  .env.example              # Variable template
  Taskfile.yaml             # Bootstrap automation
  bootstrap/                # Kind cluster config
  infrastructure/
    bases/proxmox-talos/    # Reusable CAPI templates
    management/             # Mgmt cluster overlay
    clusters/<name>/        # Workload cluster overlays
  platform/                 # Helm charts (30+ components)
  gitops/
    root-app.yaml           # Self-managing root application
    applications/           # Individual ArgoCD apps
    applicationsets/        # Auto-discovery generators
  apps/<cluster>/           # Per-cluster workload apps
  scripts/                  # Helper scripts
  docs/                     # Documentation
```
