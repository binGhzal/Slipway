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
- **Multi-cloud ready** - Architecture supports adding cloud CAPI providers without restructuring.

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
