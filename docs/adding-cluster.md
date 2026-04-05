# Adding a New Workload Cluster

## Overview

Adding a new cluster is a GitOps operation. You create a directory with Kustomize patches, commit, and ArgoCD + CAPI handle the rest.

## Steps

### 1. Copy the Template

```bash
cp -r infrastructure/clusters/homelab infrastructure/clusters/<your-cluster-name>
```

### 2. Edit the Patches

Update each file in `patches/` with your cluster's values:

**`patches/cluster-patch.yaml`**
- Cluster name
- Pod and service CIDRs (must not overlap with other clusters)

**`patches/proxmox-cluster-patch.yaml`**
- Control plane VIP
- IP address range for nodes
- Gateway, prefix
- Allowed Proxmox nodes

**`patches/control-plane-patch.yaml`**
- Number of control plane replicas (1 or 3)
- VIP address (must match ProxmoxCluster)
- ConfigPatches for Talos

**`patches/machine-template-patch.yaml`**
- Control plane VM specs (CPU, RAM, disk)

**`patches/worker-machine-template-patch.yaml`**
- Worker VM specs

**`patches/machine-deployment-patch.yaml`**
- Number of worker replicas

**`patches/worker-config-patch.yaml`**
- Worker Talos config patches

### 3. Update Namespace

Edit `namespace.yaml` and `kustomization.yaml` with the new cluster name.

### 4. Commit and Push

```bash
git add infrastructure/clusters/<your-cluster-name>
git commit -m "feat: add <your-cluster-name> workload cluster"
git push
```

### 5. What Happens Automatically

1. ArgoCD's `workload-clusters` ApplicationSet discovers the new directory
2. An ArgoCD Application is created for the cluster
3. CAPI resources are applied to the management cluster
4. CAPI creates VMs on Proxmox from the Talos template
5. Talos bootstraps Kubernetes on the VMs
6. Kyverno policy auto-registers the cluster in ArgoCD
7. Platform services deploy based on cluster labels

## Network Planning

Ensure no CIDR overlaps between clusters:

| Cluster | Pod CIDR | Service CIDR | Node IPs | VIP |
|---------|----------|--------------|----------|-----|
| slipway-mgmt | 10.244.0.0/16 | 10.96.0.0/16 | 10.10.0.11-19 | 10.10.0.10 |
| homelab | 10.245.0.0/16 | 10.97.0.0/16 | 10.10.0.21-39 | 10.10.0.20 |
| *new-cluster* | 10.246.0.0/16 | 10.98.0.0/16 | 10.10.0.41-59 | 10.10.0.40 |
