# Slipway Platform Guidelines

## Architectural Tenets

- Build a five-layer platform (infrastructure -> platform services) and validate each layer before advancing.
- Use Cluster API (CAPMOX on Proxmox) as the authoritative interface for cluster lifecycle; keep Terraform limited to foundational Proxmox setup and external clouds.
- Maintain provider neutrality by abstracting cloud resources behind Crossplane compositions and GitOps workflows.
- Prefer declarative, reconciled systems; avoid one-off scripts unless encoded back into automation.
- Treat Proxmox as the golden environment for experimentation, but ensure all patterns translate to multi-cloud targets with minimal drift.
- Capture every material architectural or process decision in an Architecture Decision Record (ADR) before implementation.

## Layer Guidance

### Layer 1 - Core Infrastructure

- Standardize Proxmox nodes with API tokens, VLAN-aware networking, and Ceph (multi-node) or ZFS (single-node) storage profiles.
- Manage VM templates exclusively via Kubernetes Image Builder pipelines; manual templates are prohibited.
- Version all base images with semantic tags (os-version_k8s-version_build-id) and publish provenance metadata.
- Enforce hardware baselines (CPU, RAM, disk layout) per cluster size class; document exceptions in runbooks.

### Layer 2 - Container Orchestration

- Host the management plane on a hardened K3s VM; restrict SSH and expose only Cluster API endpoints and GitOps ingress.
- Pin supported Kubernetes versions per environment, track End-of-Life dates, and document upgrade playbooks.
- Adopt Cilium as the default CNI, enabling kube-vip (or equivalent) for HA control planes and enforcing network policies by default.
- Require automated conformance testing for every new cluster template before promotion to production use.

### Layer 3 - Management Platform

- Enforce GitOps-first delivery through Argo CD (primary) or Flux (fallback) with an app-of-apps or ApplicationSet hierarchy.
- Package platform components with Helm charts or OCI artifacts; keep values files scoped per environment and stored in Git.
- Gate policy changes through OPA (Gatekeeper/Conftest) with automated tests in CI; document exceptions with expiry dates.
- Manage secrets through External Secrets Operator; direct Kubernetes Secret manifests in Git are disallowed.

### Layer 4 - Multi-Cloud Abstraction

- Model infrastructure with Crossplane compositions and managed resources; keep cloud credentials in External Secrets-backed providers.
- Use Terraform only to bootstrap providers lacking Crossplane coverage, then reconcile resulting state into Git manifests.
- Extend the service mesh (Istio) only after traffic, identity, and security requirements are reviewed and documented.
- Validate latency, failover, and cost profiles for every new cloud footprint before onboarding production workloads.

### Layer 5 - Platform Services

- Treat Backstage as the single portal for developer self-service; expose golden paths via software templates and catalog entries.
- Baseline observability with kube-prometheus-stack; require dashboards and alerts for every critical workload and shared component.
- Integrate CI/CD pipelines (Tekton) with signed artifacts, automated promotion gates, and policy checks prior to deployment.
- Embed security scanning (Trivy, Falco, image signing) into platform services; surface findings in a central reporting channel.

## Governance and Process

- Operate in phase gates aligned to the implementation plan; do not promote to the next phase until exit criteria are met.
- Establish an architecture review board (ARB) to approve ADRs, major dependency changes, and cross-layer impacts.
- Run bi-weekly platform demos and retrospectives; capture action items in the shared backlog.
- Maintain environment parity matrices (Proxmox, AWS, Azure, GCP) and review drift monthly.
- Record all production changes through tracked Git PRs with linked tickets; emergency fixes require retroactive documentation within 24 hours.

## Security and Compliance

- Apply zero-trust principles: identity-aware access to all APIs, TLS in transit, disk encryption at rest, and mandatory network policies.
- Store secrets in External Secrets Operator backed by Vault or cloud secret managers; enforce automated rotation every 90 days or sooner for high-risk credentials.
- Run continuous security scans (Trivy for images, Falco for runtime, Snyk or Grype for dependencies) and triage high findings within 48 hours.
- Maintain audit logs for Proxmox, Kubernetes APIs, GitOps reconciliations, and Crossplane actions; archive logs in immutable storage for 365 days.
- Require threat modeling for new platform services; document mitigations and integrate controls into policies.

## Operations and Reliability

- Define and publish SLOs: cluster provisioning (<5 minutes), deployment success (>99%), and Crossplane reconciliation (<2 minutes) as baselines.
- Implement automated backup and recovery for etcd, Proxmox storage, Git repositories, and Crossplane state; execute recovery drills quarterly.
- Maintain incident response playbooks, on-call rotations, and escalation paths; run post-incident reviews within five business days.
- Track cost, utilization, and performance metrics per environment; feed optimization recommendations (spot usage, right-sizing, storage tiering) into the roadmap.
- Use feature flags or canary environments for risky changes; never deploy unvalidated configurations directly to production clusters.

## Documentation and Knowledge Management

- Keep ADRs, runbooks, playbooks, and design docs versioned in Git; publish final versions to Backstage TechDocs.
- Provide onboarding checklists for platform engineers, security, and application teams; update them after every phase gate.
- Document GitOps repository structures, branching strategies, and promotion workflows with practical examples.
- Maintain an up-to-date dependency catalog (cluster API providers, Helm charts, Terraform modules) with ownership and support tiers.
- Record operational learnings and best practices in the internal wiki; link to relevant runbooks and ADRs.

## Coding and Configuration Standards

### General

- Favor declarative configuration; imperative scripts must be idempotent and wrapped in automation.
- Keep code and configuration DRY; extract shared logic into modules, libraries, or reusable manifests.
- Document non-obvious intent with succinct comments; avoid duplicating what is already clear from the code.
- Validate every change with automated tests or linting before requesting review.

### Git Workflow

- Name branches using `feature/<ticket-id>-<description>` or `fix/<ticket-id>-<description>`; use hyphens for word separators.
- Write conventional commits (type(scope): summary) or descriptive messages that reference associated work items.
- Require at least one peer review for all changes; reviewers must check lint results, tests, and policy impacts before approval.
- Squash merge by default to keep history clean; preserve important details in the commit message.

### YAML, Helm, and Kubernetes Manifests

- Use two-space indentation, quoted strings only when necessary, and snake_case for keys in values files.
- Place Kubernetes manifests under declarative GitOps repositories with environment overlays; avoid using kubectl apply manually.
- Keep Helm chart defaults minimal; use values files per environment and document required overrides.
- Validate manifests with kubeconform or kubeval and enforce policies with OPA/Conftest in CI.

### Terraform and Infrastructure as Code

- Structure Terraform repositories as modules (reusable) and stacks (environment-specific); version modules semantically.
- Indent with two spaces, align equals signs, and order blocks (provider, data, resource, output) logically.
- Use lower*snake_case for resource names, camelCase for variables, and prefix sensitive outputs with sensitive* and mark them as sensitive.
- Store remote state in a backend with locking (e.g., S3 + DynamoDB, GCS + locking); never check state files into Git.
- Run terraform fmt, validate, and security scans (tfsec, checkov) on every change; block merges on failures.

### Scripting and Automation

- Prefer Go or Python for complex automation; keep bash scripts for simple orchestration only.
- Include set -euo pipefail at the top of bash scripts and shellcheck them in CI.
- Package Go code with gofmt/goimports and run golangci-lint; pin dependencies via go.mod and go.sum.
- For Python, follow PEP 8, format with black, and manage dependencies with poetry or pip-tools; run mypy where practical.

### Testing and Validation

- Provide unit tests for modules, integration tests for provisioning workflows, and smoke tests for cluster bootstrap flows.
- Automate conformance tests (e.g., sonobuoy) for Kubernetes upgrades and template changes.
- Run policy-as-code tests (OPA, Conftest) as part of CI; include regression coverage for known incidents.
- Require success criteria and rollback steps in every change proposal; document validation evidence in the PR.

## Tooling Rules

- Use Argo CD, Flux, and Crossplane CLIs only in read-only mode in production; apply changes via Git commits.
- Manage secrets with External Secrets Operator and backed managers; sharing credentials via chat or email is forbidden.
- Track dependencies with Renovate or Dependabot; review automatic upgrades within five business days.
- Maintain parity between local development environments (kind, tilt) and management clusters; document setup scripts.

## Enforcement

- The platform engineering team owns these guidelines and audits compliance quarterly.
- Violations must be remediated promptly, with follow-up actions tracked in the backlog and documented in retrospectives.
