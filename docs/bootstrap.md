# Bootstrap Guide

## Prerequisites

### Tools
- `kind` - Local Kubernetes clusters
- `kubectl` - Kubernetes CLI
- `clusterctl` - Cluster API CLI
- `talosctl` - Talos Linux CLI
- `helm` - Helm package manager
- `cilium` - Cilium CLI

### Infrastructure
- Proxmox VE with API access
- API token with PVEVMAdmin role
- Talos VM template created (see `docs/talos-image.md`)

### Accounts
- Cloudflare account with API token (DNS edit permissions)
- Secret store (1Password, Infisical, Doppler, or Vault)

## Quick Start

```bash
# 1. Clone and configure
git clone https://github.com/binGhzal/Slipway.git
cd Slipway
cp .env.example .env
# Edit .env with your values

# 2. Full automated bootstrap
task full-bootstrap
```

## Step-by-Step Bootstrap

### 1. Initialize
```bash
task init    # Validates .env and checks tools
```

### 2. Create Kind Bootstrap Cluster
```bash
task bootstrap:up    # Creates Kind + installs Cilium
```

### 3. Install CAPI Providers
```bash
task bootstrap:capi  # Installs CAPI Operator + Proxmox/Talos providers
```

### 4. Deploy Management Cluster
```bash
task mgmt:secrets    # Creates Proxmox credentials from .env
task mgmt:apply      # Applies CAPI resources
task mgmt:wait       # Waits for cluster (up to 20min)
```

### 5. Configure Management Cluster
```bash
task mgmt:kubeconfig # Extracts kubeconfig
task mgmt:cilium     # Installs Cilium + L2/Gateway config
task mgmt:argocd     # Installs ArgoCD
task mgmt:eso        # Installs External Secrets Operator
```

### 6. Pivot
```bash
task pivot           # Moves CAPI control to mgmt cluster
task bootstrap:down  # Deletes Kind (no longer needed)
```

### 7. GitOps Bootstrap
```bash
task gitops:bootstrap  # ArgoCD takes over from Git
```

## After Bootstrap

Everything is now managed by ArgoCD. All future changes go through Git:

1. Edit files in the repo
2. Commit and push
3. ArgoCD syncs automatically
4. Controllers reconcile the desired state

### Useful Commands
```bash
task status          # Show clusters, machines, ArgoCD apps
task mgmt:dashboard  # Get ArgoCD URL and admin password
```

## Rebuilding From Scratch

The platform is designed to be fully reproducible:

```bash
# Destroy everything, then:
task full-bootstrap
```

All cluster configs, platform services, and app deployments are in Git and will be automatically restored.
