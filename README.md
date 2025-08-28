# Refrences

OpenTofu/Terraform:

- [GitHub: BPG/Terraform-Provider-Proxmox](https://github.com/bpg/terraform-provider-proxmox)

Linux Container Images:

- [Linux Containers | LXC Images](https://images.linuxcontainers.org/images/)
  - We suggesting using the cloud-init enabled images: `<DISTRO>/amd64/cloud/rootfs.tar.xz`.

Linux VM Cloud Images:

- [OpenStack: Cloud Images], collection of image links.
- [CentOS Cloud Images]
  - Default User: `centos`
  - Use `CentOS-Stream-GenericCloud-X-latest.x86_64.qcow2`
- [Debian Cloud Images]
  - Default User: `debian`
  - Use `debian-1x-generic-amd64.qcow2`
  - Avoid:
    - `genericcloud` images fail to run `cloud-init`
    - `nocloud` images do not have `cloud-init` installed and defaults to a password-less `root` user.
- [Fedora Cloud Images]
  - Default User: `fedora`
  - Use `Fedora-Cloud-Base-Generic.x86_64-XX-XX.qcow2`
- [Ubuntu Cloud Images]
  - Default User: `ubuntu`
  - Use `ubuntu-2x.04-server-cloudimg-amd64.img`

[Terraform]: https://github.com/hashicorp/terraform
[OpenTofu]: https://opentofu.org/
[Proxmox]: https://www.proxmox.com/
[BPG Proxmox]: https://github.com/bpg/terraform-provider-proxmox
[GitHub: BPG/Terraform-Provider-Proxmox]: https://github.com/bpg/terraform-provider-proxmox
[CentOS Cloud Images]: https://cloud.centos.org/
[Debian Cloud Images]: https://cloud.debian.org/images/cloud/
[Fedora Cloud Images]: https://fedoraproject.org/cloud/download
[Ubuntu Cloud Images]: https://cloud-images.ubuntu.com/releases/
[OpenStack: Cloud Images]: https://docs.openstack.org/image-guide/obtain-images.html
[packer-proxmox-templates]: https://github.com/trfore/packer-proxmox-templates
[proxmox-template-scripts]: https://github.com/trfore/proxmox-template-scripts
[terraform-bpg-proxmox]: https://github.com/trfore/terraform-bpg-proxmox
[terraform-telmate-proxmox]: https://github.com/trfore/terraform-telmate-proxmox
[terraform-bpg-proxmox]: https://github.com/trfore/terraform-bpg-proxmox.git
[boilerplates]: https://github.com/ChristianLempa/boilerplates.git