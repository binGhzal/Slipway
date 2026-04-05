# Secret Management

## Philosophy

**No manual encryption. No SOPS. Fully hands-off.**

Secrets are stored in an external secret store (1Password, Infisical, Doppler, or Vault). The External Secrets Operator (ESO) automatically syncs them into Kubernetes Secrets. You commit `ExternalSecret` CRDs to Git (which reference secrets by name, not value), and ESO keeps them in sync.

## How It Works

```
Secret Store (1Password/Infisical/Vault)
    |
    v
External Secrets Operator (runs in cluster)
    |
    v
Kubernetes Secret (auto-created, auto-refreshed)
    |
    v
Pod (consumes secret as env var or volume mount)
```

## Setup

### 1. Choose a Backend

Edit `platform/external-secrets/secret-store.yaml` and uncomment your chosen backend.

### 2. Store Secrets in Your Backend

Create entries in your secret store for:
- **cloudflare** - `api-token` field
- **proxmox** - `url`, `token`, `secret` fields
- **cloudflare-tunnel** - `token` field
- Any app-specific secrets

### 3. Bootstrap Credential

ESO needs one credential to connect to the secret store. This is created during bootstrap by `task mgmt:eso` from your `.env` values. It's the only secret not managed by ESO itself.

### 4. ExternalSecret CRDs

For each secret needed, create an `ExternalSecret` in `platform/external-secrets/secrets/`:

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: my-app-secret
  namespace: my-app
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: onepassword  # Match your ClusterSecretStore name
    kind: ClusterSecretStore
  target:
    name: my-app-secret
    creationPolicy: Owner
  data:
    - secretKey: database-url
      remoteRef:
        key: my-app         # Item name in secret store
        property: db-url    # Field name in the item
```

### 5. Verify

```bash
kubectl get externalsecrets -A    # All should show "SecretSynced"
kubectl get secrets -n my-app     # ESO-created secrets appear here
```

## Rotating Secrets

1. Update the secret in your secret store
2. ESO automatically picks up the change within `refreshInterval`
3. Reloader detects the Secret change and restarts affected pods

No Git commits needed for secret rotation.
