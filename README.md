# Slipway

Slipway is a Kubernetes-first PaaS with a tiny core and optional plugins.

## Quickstart with Kind

```
make kind-up && make build && make docker-build && make kind-deploy
kubectl apply -f deploy/examples/01-project-acme.yaml
kubectl apply -f deploy/examples/02-plugins.yaml
kubectl apply -f deploy/examples/03-app-api.yaml
```

TLS requires cert-manager plugin; otherwise apps run HTTP only.
Security defaults: restricted Pod Security and default-deny NetworkPolicy.

## Architecture

```
Project -> App -> Release
         |            
         -> Plugins (capabilities)
```

