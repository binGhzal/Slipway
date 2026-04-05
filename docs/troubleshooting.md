# Troubleshooting

## Talos Debugging (No SSH)

Talos has no SSH. Use `talosctl` for all debugging:

```bash
# Dashboard (live view of node status)
talosctl --talosconfig <path> dashboard

# View kernel messages
talosctl --talosconfig <path> dmesg

# View service logs
talosctl --talosconfig <path> logs kubelet
talosctl --talosconfig <path> logs containerd

# Get machine config
talosctl --talosconfig <path> get machineconfig

# Check etcd status
talosctl --talosconfig <path> etcd status
```

## CAPI Troubleshooting

```bash
# Cluster status overview
clusterctl describe cluster <name> -n <namespace>

# Watch machine provisioning
kubectl get machines -A -w

# Check CAPI controller logs
kubectl logs -n capmox-system deploy/capmox-controller-manager
kubectl logs -n cabpt-system deploy/cabpt-controller-manager
kubectl logs -n cacppt-system deploy/cacppt-controller-manager

# Check events for errors
kubectl get events -A --sort-by='.lastTimestamp' | tail -50
```

## Common Issues

### Machine stuck in Provisioning
- Check `skipCloudInitStatus: true` is set in ProxmoxMachineTemplate
- Verify QEMU guest agent extension is in the Talos image
- Verify guest agent is enabled on the Proxmox template (`qm set <id> --agent enabled=1`)

### IP Address Not Assigned
- Check IPAM pool has available addresses
- Verify the IP range doesn't overlap with DHCP
- Check `kubectl get ipaddresses -A` and `kubectl get ipaddressclaims -A`

### Nodes Not Ready
- Check if CNI (Cilium) is installed: `cilium status`
- Check cloud-provider-external is set in kubelet args
- Verify Talos Cloud Controller Manager is running

### ArgoCD Not Syncing
- Check ArgoCD app status: `kubectl get applications -n argocd`
- Check ArgoCD logs: `kubectl logs -n argocd deploy/argocd-server`
- Verify Git repo is accessible from the cluster

### External DNS Not Creating Records
- Check external-dns logs: `kubectl logs -n external-dns deploy/external-dns`
- Verify Cloudflare API token has DNS edit permissions
- Check `txtOwnerId` doesn't conflict with other external-dns instances

### Cilium Issues
- Run `cilium status` for component health
- Run `cilium connectivity test` for network validation
- Check Hubble: `hubble observe --verdict DROPPED`
