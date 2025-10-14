# Slipway Copilot Instructions

## Orientation

- This repo is currently planning-heavy; core guidance lives in `docs/guidelines.md`, `docs/plan.md`, and `docs/research.md`. Read them first to understand the five-layer platform strategy.
- Capture any material decision or deviation as an ADR before coding; the guidelines treat ADRs as gatekeepers for implementation work.
- Keep deliverables declarative and GitOps-friendly—scripts are acceptable only when they feed automation back into Git-managed workflows.

## Architecture Blueprint

- Slipway is organized as five layers (Proxmox infrastructure → Kubernetes orchestration → GitOps management → Crossplane multi-cloud → platform services). Align new work with the layer definitions and exit criteria in `docs/plan.md`.
- Cluster lifecycle must flow through Cluster API with the Proxmox provider (CAPMOX); Terraform is reserved for bootstrapping gaps and must reconcile back into Git manifests.
- Management plane expectations: hardened K3s VM hosting Cluster API endpoints, Cilium CNI, kube-vip for HA. Any deviation requires an ADR and updated runbooks.

## Delivery Workflow

- Phase gates in `docs/plan.md` dictate sequencing—avoid skipping layers. Each phase needs validation evidence (e.g., conformance tests, latency benchmarks) committed alongside configuration.
- Prefer Argo CD for GitOps; Flux is acceptable fallback. Use app-of-apps/ApplicationSets patterns, and segregate environment values per the GitOps structure described in `docs/guidelines.md`.
- Secrets must flow through External Secrets Operator references; never commit raw `Secret` manifests.

## Implementation Conventions

- When authoring automation, favor Go or Python; include lint/format hooks (`gofmt`, `golangci-lint`, or `black`/`mypy`). Shell scripts require `set -euo pipefail` and `shellcheck` coverage.
- Kubernetes manifests use two-space indentation and should pass `kubeconform`/OPA checks. Helm values follow snake_case keys and environment-specific overlays.
- Terraform stacks are for environment wiring only; run `terraform fmt`, `validate`, and security scanners (tfsec/checkov) before opening PRs.

## Verification & Evidence

- Document provisioning targets (cluster creation <5 min, Crossplane reconciliation <2 min) and attach test outputs or dashboards when improvements affect SLOs.
- For new cloud integrations, capture latency/cost benchmarks and Crossplane composition examples under a clearly-named folder (e.g., `docs/multi-cloud/<provider>`).
- Update onboarding references, runbooks, or TechDocs sources when workflows change; this repo is the canonical source for platform knowledge.

## Collaboration

- Branch names: `feature/<ticket>-<slug>` or `fix/<ticket>-<slug>`; commits should be conventional or reference work items explicitly.
- Every change needs a peer review plus passing automation results. Emergency fixes must get retroactive docs within 24 hours per `docs/guidelines.md`.
- Surface open questions in PR descriptions when requirements are unclear—platform assumptions are strict, and undocumented shortcuts are rejected.
