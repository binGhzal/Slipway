# Slipway

Infrastructure boilerplates and automation to build images, provision VMs/LXCs on Proxmox, and stand up a minimal K3s cluster. Batteries included for Packer, Terraform (BPG Proxmox provider), and Ansible.

## Overview

Slipway helps you:

- Build reusable Proxmox VM templates with Packer (Ubuntu 24.04 a.k.a. Noble).
- Download/prepare cloud images on Proxmox for VMs or LXC templates via Terraform modules.
- Clone templates into VMs, optionally with cloud-init and static IPs, or create LXCs with mounts and VLANs.
- Bootstrap a small K3s cluster (1 controller + N agents) on Proxmox VMs.
- Automate VM prep and Kubernetes install steps with Ansible playbooks.

The repo is organized so you can pick-and-choose: use only Packer, only Terraform provisioning, or run the full flow.

## Repository structure

Top-level folders and what they contain:

- `packer/`
  - `packer.pkr.hcl` and plugin setup
  - `ubuntu-server-noble/` Packer template, cloud-init http files, and `files/99-pve.cfg` for Proxmox integration
  - `variables.pkrvars.hcl.template` for local secrets/connection
- `terraform/`
  - `proxmox/`
    - `image-download/` module usage to fetch images (.img/.iso and LXC `vztmpl`) to Proxmox storage
    - `vm-clone/` module usage to clone a VM from a template with cloud-init
    - `lxc/` module usage to create LXCs (including static IPs and mount points)
    - `k3s-cluster/` example that clones multiple VMs and installs K3s
  - `templates/` cloud provider examples (Cloudflare, Civo, Kubernetes automation)
- `ansible/`
  - `ubuntu/` small utility plays for apt, qemu-guest-agent, zsh, reboot, ssh keys
  - `kubernetes/` a guided multi-play install for Kubernetes prerequisites, control-plane init, and node join
- `.github/` policies and lint workflow

## Prerequisites

- Proxmox VE with API access and a storage pool for ISOs/templates (e.g., `local`, `local-lvm`, or `bigdisk`).
- macOS or Linux workstation with:
  - Packer 1.9+
  - Terraform 1.5+ (or OpenTofu 1.6+)
  - Ansible 2.14+
  - SSH key configured for VM access
- Proxmox API token with sufficient permissions for image upload and VM/LXC operations.

## Quick start

1. Clone the repo and create your local variables files/secrets (never commit secrets):

- For Packer: copy `packer/variables.pkrvars.hcl.template` to `packer/secrets.pkrvars.hcl` and fill in values.
- For Terraform: create `*.auto.tfvars` files alongside the examples containing your Proxmox endpoint and token.

1. Choose a workflow below and run it end-to-end.

## Workflow A: Build a Proxmox VM template with Packer

Files:

- `packer/packer.pkr.hcl`
- `packer/ubuntu-server-noble/ubuntu-server-noble.pkr.hcl`
- `packer/ubuntu-server-noble/http/{meta-data,user-data}`

Steps:

- Initialize the Proxmox Packer plugin via required_plugins or install it globally.
- Edit `ubuntu-server-noble.pkr.hcl` to set your Proxmox node name, storage pools, and ISO source (local or download).
- Optionally adjust `http/user-data` for cloud-init autoinstall settings.
- Build:
  - From `packer/ubuntu-server-noble/`: run Packer with your vars file.

Note: The template cleans machine-id, enables qemu-guest-agent, and drops `files/99-pve.cfg` for cloud-init in Proxmox.

## Workflow B: Download cloud images to Proxmox storage (Terraform)

Folder: `terraform/proxmox/image-download/`

What it does:

- Uses the BPG Proxmox provider to upload/convert VM images (.qcow2 to .img) or fetch LXC templates (`vztmpl`).

Important variables (see `variables.tf`):

- `pve_api_url`, `pve_token_id`, `pve_token_secret` — Proxmox API auth.
- For LXC templates provide `image_content_type = "vztmpl"`.

Usage:

- Create a `*.auto.tfvars` with your credentials and desired image URL/checksum and apply.

## Workflow C: Clone a VM from a template (Terraform)

Folder: `terraform/proxmox/vm-clone/`

What it does:

- Clones a Proxmox template into one or many VMs.
- Supports cloud-init: user, SSH key, DNS, static IPs, vendor/user/meta data.

Key inputs:

- `template_id`, `vm_id`, `node`, `vcpu`, `memory`, `disks[]`
- Cloud-init inputs like `ci_ssh_key`, `ci_ipv4_cidr`, `ci_ipv4_gateway`, `ci_vendor_data`.

Usage patterns in `main.tf` include single VM, multiple VMs, UEFI, and cross-node template cloning.

## Workflow D: Create LXCs (Terraform)

Folder: `terraform/proxmox/lxc/`

What it does:

- Provisions unprivileged LXCs from templates; supports static IPs, VLAN tags, mount points (including bind mounts).

Notes:

- Bind mounts to host storage require root SSH access configured in the provider.

## Workflow E: Stand up a K3s cluster (Terraform + cloud-init + remote-exec)

Folder: `terraform/proxmox/k3s-cluster/`

What it does:

- Uploads `files/vendor-data.yaml` as a Proxmox snippet.
- Clones VMs for a controller and agents from a template.
- Installs K3s on controller and agents using a shared token via remote-exec.

Outputs:

- Controller IP and agent IPs to help you connect and verify.

## Ansible plays

Folders: `ansible/ubuntu/` and `ansible/kubernetes/`

- Ubuntu utilities: apt updates, qemu-guest-agent, zsh, reboots, SSH key deploy.
- Kubernetes guided install: prepares containerd, disables swap, fetches keys, installs kubeadm/kubelet/kubectl, initializes cluster, applies Flannel CNI, generates and fetches the join command, and joins workers.

Inventory expectations for `ansible/kubernetes` are described in its README.

## Configuration and secrets

- Do not commit secrets; `.gitignore` already excludes `**/secrets.*` and `.env` files.
- Example Packer vars live at `packer/variables.pkrvars.hcl.template`; create a private copy and keep it local.
- For Terraform, store tokens and endpoints in `*.auto.tfvars` or environment variables.

## Contributing

See `.github/CONTRIBUTING.md`, `.github/CODE_OF_CONDUCT.md`, and `.github/SECURITY.md`. MIT licensed.

## References

- [BPG Proxmox Provider](https://github.com/bpg/terraform-provider-proxmox)
- [Linux Containers | LXC Images](https://images.linuxcontainers.org/images/) (prefer cloud-init enabled `<DISTRO>/amd64/cloud/rootfs.tar.xz`)
- [OpenStack cloud image sources](https://docs.openstack.org/image-guide/obtain-images.html)
- [Ubuntu Cloud Images](https://cloud-images.ubuntu.com/releases/) (user: `ubuntu`)
- [Debian Cloud Images](https://cloud.debian.org/images/cloud/) (user: `debian`)
- [Fedora Cloud Images](https://fedoraproject.org/cloud/download) (user: `fedora`)
- [CentOS Stream](https://cloud.centos.org/) (user: `centos`)
- Related repos/templates:
  - [trfore/packer-proxmox-templates](https://github.com/trfore/packer-proxmox-templates)
  - [trfore/proxmox-template-scripts](https://github.com/trfore/proxmox-template-scripts)
  - [trfore/terraform-bpg-proxmox](https://github.com/trfore/terraform-bpg-proxmox)
  - [ChristianLempa/boilerplates](https://github.com/ChristianLempa/boilerplates)
