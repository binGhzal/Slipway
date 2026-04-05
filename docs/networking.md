# Networking Architecture

## Overview

All networking is handled by Cilium, replacing kube-proxy, MetalLB, and traditional ingress controllers.

## Components

### Cilium CNI (kube-proxy replacement)
- eBPF-based packet processing (no iptables)
- 25-40% CPU reduction vs kube-proxy
- Socket-level load balancing

### LB-IPAM + L2 Announcements (replaces MetalLB)
- `CiliumLoadBalancerIPPool` defines available IPs
- `CiliumL2AnnouncementPolicy` advertises IPs via ARP
- LoadBalancer services get IPs from the pool automatically

### Gateway API (ingress)
- Replaces ingress-nginx (archived March 2026)
- `GatewayClass` -> `Gateway` -> `HTTPRoute`/`TLSRoute`
- eBPF-based traffic management

### WireGuard Encryption
- Transparent node-to-node encryption
- No extra configuration needed

### Hubble (observability)
- Network flow visualization
- Policy verdict monitoring
- Accessible via `hubble` CLI or Grafana dashboards

## DNS Automation Flow

```
1. Create HTTPRoute with hostname
2. Cilium Gateway API creates LB Service
3. LB-IPAM assigns IP from CiliumLoadBalancerIPPool
4. L2 Announcement makes IP reachable on network
5. external-dns creates A record in Cloudflare
6. cert-manager issues TLS cert via DNS-01 challenge
7. cloudflared tunnels external traffic from Cloudflare edge
```

## IP Addressing

| Purpose | Range | Notes |
|---------|-------|-------|
| Management nodes | 10.10.0.10-19 | CAPI IPAM managed |
| Homelab nodes | 10.10.0.20-39 | CAPI IPAM managed |
| LoadBalancer pool | 10.10.1.0/24 | Cilium LB-IPAM |
| Pod CIDR (mgmt) | 10.244.0.0/16 | Per-cluster unique |
| Pod CIDR (homelab) | 10.245.0.0/16 | Per-cluster unique |
| Service CIDR (mgmt) | 10.96.0.0/16 | Per-cluster unique |
| Service CIDR (homelab) | 10.97.0.0/16 | Per-cluster unique |

## Network Policies

Default-deny policy is deployed via `platform/cilium-config/default-deny.yaml`. All traffic within the cluster is allowed, external DNS is allowed. Additional policies can be added per namespace.

Start in audit mode to verify before enforcement:
```bash
kubectl annotate ccnp default-deny io.cilium.policy/audit-mode=enabled
```
