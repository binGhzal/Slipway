# Templates

Ready-made templates to help you structure Kubernetes configuration for Slipway applications. Copy this directory into your own Git repository and customize it for your needs.

```
├── bootstrap
│   ├── helmrepositories
│   │   └── helmrepository-podinfo.yaml
│   └── namespaces
│       └── namespace-podinfo.yaml
└── podinfo
    ├── configmap-podinfo-helm-chart-value-overrides.yaml
    └── helmrelease-podinfo.yaml
```

These templates are generic and not tied to any specific GitOps tool, but they can be consumed by systems like Flux if desired.
