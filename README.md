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

## Templates

Ready-made Kubernetes manifests are available in the [templates](templates/) directory to help you bootstrap apps quickly. The structure keeps configuration organized and can be adopted by GitOps tools like Flux in the future.

```
templates
├── bootstrap
│   ├── helmrepositories
│   │   └── helmrepository-podinfo.yaml
│   └── namespaces
│       └── namespace-podinfo.yaml
└── podinfo
    ├── configmap-podinfo-helm-chart-value-overrides.yaml
    └── helmrelease-podinfo.yaml
```

Copy these templates and adjust them for your own applications.

## Development Plan and Roadmap

See [PROJECT_PLAN.md](docs/PROJECT_PLAN.md) for the repository structure, development guidelines, and future roadmap.

