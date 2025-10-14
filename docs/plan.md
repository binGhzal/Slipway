# Slipway Implementation Plan

## Overview

- Objective: deliver a Proxmox-first Kubernetes platform with multi-cloud expansion, using Cluster API + Crossplane and GitOps automation.
- Horizon: 14-week roadmap in five phases aligned with the bottom-up architecture.
- Success criteria: production-ready base platform, documented runbooks, and repeatable GitOps workflows.

## Phase Breakdown

### Phase 1 – Foundation (Weeks 1-2)

- Outcomes: production-grade Proxmox cluster, standardized VM templates, baseline monitoring.
- Key Tasks:
    - Harden Proxmox nodes, enable API tokens, configure VLAN networking.
    - Deploy Ceph (multi-node) or ZFS (single-node) storage; document capacity plans.
    - Implement Image Builder pipeline for Ubuntu 24.04 + Kubernetes images; store artifacts in registry.
    - Stand up initial observability (node metrics, logging) and validate alerting hooks.
- Deliverables: infrastructure ADRs, automated template pipeline, baseline monitoring dashboards.
- Exit Criteria: template build runs end-to-end; all nodes reachable; runbook for provisioning new hardware.

### Phase 2 – Kubernetes Bootstrap (Weeks 3-5)

- Outcomes: management plane operational, first workload cluster delivered via Cluster API.
- Key Tasks:
    - Provision hardened K3s management VM; lock down SSH and open required ports only.
    - Install Cluster API with CAPMOX; seed cluster templates in Git.
    - Configure Cilium and kube-vip; validate HA control plane story.
    - Provision and smoke-test first workload cluster (conformance + e2e app deployment).
- Deliverables: reusable cluster templates, networking validation report, upgrade playbook draft.
- Exit Criteria: clusters provision in <5 minutes; automated tests green; rollback procedure documented.

### Phase 3 – GitOps Enablement (Weeks 6-7)

- Outcomes: fully declarative delivery via Argo CD/Flux, policy guardrails enforced.
- Key Tasks:
    - Deploy Argo CD (primary) or Flux (fallback) with multi-cluster targets.
    - Structure Git repos (infrastructure, clusters, applications) with app-of-apps pattern.
    - Integrate OPA/Gatekeeper for policy-as-code; add security validation to CI.
    - Introduce External Secrets Operator; connect to secret backends.
- Deliverables: GitOps repo templates, policy rule catalog, secret management runbook.
- Exit Criteria: platform changes flow through Git; policy violations block merges; secrets synced automatically.

### Phase 4 – Multi-Cloud Expansion (Weeks 8-11)

- Outcomes: Crossplane control plane live with at least one public cloud integration; cross-cluster connectivity defined.
- Key Tasks:
    - Install Crossplane with Upbound providers; secure credentials via External Secrets.
    - Design and implement compositions for core services (networking, database, cluster provisioning) across AWS/Azure/GCP.
    - Integrate Terraform where Crossplane coverage gaps exist; reconcile state back into Git.
    - Plan and pilot Istio service mesh for cross-cloud traffic with documented SLOs.
- Deliverables: Crossplane package catalog, multi-cloud reference architecture, network latency benchmarks.
- Exit Criteria: managed resources reconcile successfully; failover test executed; documentation reviewed.

### Phase 5 – Platform Services (Weeks 12-14)

- Outcomes: developer self-service portal, CI/CD integration, security instrumentation.
- Key Tasks:
    - Deploy Backstage with service catalog, TechDocs, and golden path templates.
    - Roll out kube-prometheus-stack dashboards, alert routing, and cost reporting.
    - Implement Tekton pipelines with artifact signing and environment promotion strategy.
    - Integrate Trivy and Falco; establish incident response workflow.
- Deliverables: Backstage portal MVP, CI/CD reference pipelines, security posture report.
- Exit Criteria: developers can provision sandbox clusters via templates; monitoring/alerts validated; security tooling active.

## Cross-Cutting Workstreams

- **Automation & Testing**: maintain CI for Image Builder, Cluster API templates, policy checks, and Crossplane compositions.
- **Documentation**: capture ADRs, runbooks, onboarding guides per phase; publish in Backstage TechDocs.
- **Governance**: hold phase-gate reviews, track risk mitigations, update dependency map.
- **Cost & Performance**: collect metrics on provisioning time, utilization, and cross-cloud spend; feed into optimization backlog.

## Resource & Role Assumptions

- Platform Engineering (2-3 FTE): owns automation, cluster management, Crossplane compositions.
- Infrastructure Engineering (1-2 FTE): manages Proxmox, networking, storage.
- Security/Compliance (0.5 FTE): policies, secrets, audit integration.
- Developer Experience (1 FTE): Backstage, CI/CD, templates.

## Risks & Mitigations

- **CAPMOX maturity gaps**: monitor upstream issues, maintain hotfix branch, and design fallback Terraform modules.
- **Crossplane provider drift**: pin versions, schedule monthly upgrade reviews, add canary environment.
- **Secret sprawl**: enforce External Secrets usage, audit quarterly, automate rotation policy.
- **Multi-cloud latency**: baseline network SLAs early, consider regional workloads, enable mesh traffic policy.

## Metrics & Reporting

- Provisioning duration (target <5 min), deployment success rate (>99%), policy compliance (100%), security findings SLA (<48h).
- Weekly status updates, bi-weekly demos per phase, and retro after each phase gate.

## Next Steps

1. Confirm team assignments and availability for each phase.
2. Launch infrastructure ADR cycle to lock reference architecture.
3. Kick off Phase 1 with Proxmox hardening checklist and Image Builder pipeline setup.
