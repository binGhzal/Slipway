# Deploying Applications

## Overview

Applications are deployed via GitOps. Place your app manifests in `apps/<cluster-name>/<app-name>/` and ArgoCD deploys them automatically.

## Steps

### 1. Create App Directory

```bash
mkdir -p apps/<cluster-name>/<app-name>
```

### 2. Add Manifests

Add your Kubernetes manifests (Deployment, Service, HTTPRoute, etc.) or a Helm chart with `Chart.yaml` + `values.yaml`.

### 3. Example: Simple App

```bash
mkdir -p apps/homelab/whoami
```

`apps/homelab/whoami/deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: whoami
spec:
  replicas: 1
  selector:
    matchLabels:
      app: whoami
  template:
    metadata:
      labels:
        app: whoami
    spec:
      containers:
        - name: whoami
          image: traefik/whoami:latest
          ports:
            - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: whoami
spec:
  selector:
    app: whoami
  ports:
    - port: 80
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: whoami
spec:
  parentRefs:
    - name: cilium-gateway
      namespace: kube-system
  hostnames:
    - whoami.yourdomain.com
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /
      backendRefs:
        - name: whoami
          port: 80
```

### 4. Commit and Push

```bash
git add apps/homelab/whoami/
git commit -m "feat: deploy whoami to homelab"
git push
```

### 5. Automatic DNS

If you added an HTTPRoute with a hostname:
1. Cilium creates a LoadBalancer service
2. external-dns creates a DNS record in Cloudflare
3. cert-manager issues a TLS certificate
4. Your app is accessible at `https://whoami.yourdomain.com`

## How the ApplicationSet Works

The `apps` ApplicationSet in `gitops/applicationsets/apps.yaml` uses a matrix generator:
- Clusters labeled `slipway.io/apps: "true"`
- Git directories matching `apps/{{cluster-name}}/*`

Each matching directory becomes an ArgoCD Application deployed to the corresponding cluster.
