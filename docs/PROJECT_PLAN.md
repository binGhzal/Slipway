# Slipway Development Plan and Roadmap

## Overview
Slipway is a Kubernetes-first Platform as a Service with a small core and optional plugins. This document outlines the project's structure, development rules, and roadmap.

## Repository Structure
- **core/**: Kubernetes controller and APIs that manage projects, apps, and releases.
- **cli/**: Command-line interface for interacting with Slipway clusters.
- **plugins/**: Optional extensions that provide additional capabilities (e.g. cert-manager integration).
- **ui/**: Experimental web interface for managing Slipway resources.
- **templates/**: Predefined Kubernetes manifests to bootstrap applications quickly.
- **deploy/**: Deployment manifests for running Slipway on Kubernetes.
- **examples/**: Sample configurations that demonstrate core features.
- **docs/**: Project documentation, including this plan and legal notices.

## Development Guidelines
- Write code in Go and ensure it is formatted with `gofmt` (`make lint`).
- Execute `go test ./...` before submitting changes to keep the codebase stable.
- Keep commits focused and include a `Signed-off-by` line to satisfy the DCO.
- Update or add documentation when introducing new features or behaviour.
- Discuss large or breaking changes in an issue or design proposal before implementation.

## Contribution Workflow
1. Fork the repository and create a feature branch.
2. Make your changes, keeping commits concise and well described.
3. Run `make lint` and `make test` to verify formatting and unit tests.
4. Open a pull request linking to any related issues and describing the change.

## Roadmap
### Near Term
- Stabilise the core controller and CLI.
- Expand documentation and usage examples.
- Publish initial plugin templates.

### Mid Term
- Enhance the plugin system with discovery and configuration helpers.
- Integrate with GitOps tooling such as Flux.
- Improve the web UI for managing projects and releases.

### Long Term
- Provide multi-tenant support and a plugin marketplace.
- Introduce advanced security features and policy enforcement.
- Formalise release processes and packaging for multiple platforms.

This plan will evolve as the project grows; community contributions and feedback are welcome.
