# BPG LXC Module

## Requirements

| Name        | Version  |
| ----------- | -------- |
| [terraform] | >= 1.5.0 |

## Providers

| Name          | Version   |
| ------------- | --------- |
| [bpg proxmox] | >= 0.53.1 |

## Inputs

### LXC Variables

| Variable            | Default     | Type         | Description                                                                                      | Required |
| ------------------- | ----------- | ------------ | ------------------------------------------------------------------------------------------------ | -------- |
| node                |             | String       | Name of Proxmox node to provision LXC on, e.g. `pve`                                             | **Yes**  |
| lxc_id              |             | Number       | ID number for new LXC                                                                            | **Yes**  |
| lxc_name            | `null`      | String       | Defaults to using PVE naming, e.g. `CT<LXC_ID>`                                                  | no       |
| description         | `null`      | String       | LXC description                                                                                  | no       |
| tags                | `null`      | List(String) | Proxmox tags                                                                                     | no       |
| unprivileged        | `true`      | Boolean      | Set container to unprivileged                                                                    | no       |
| os_template         |             | String       | Template for LXC, e.g. `local:vztmpl/ubuntu.tar.gz`                                              | **Yes**  |
| os_type             | `unmanaged` | String       | Container OS specific setup, uses setup scripts in `/usr/share/lxc/config/<os_type>.common.conf` | no       |
| vcpu                | `1`         | Number       | Number of CPU cores                                                                              | no       |
| memory              | `512`       | Number       | Memory size in `MiB`                                                                             | no       |
| memory_swap         | `512`       | Number       | Memory swap size in `MiB`                                                                        | no       |
| disk_storage        | `local-lvm` | String       | Disk storage location                                                                            | no       |
| disk_size           | `8`         | Number       | Disk size                                                                                        | no       |
| user_ssh_key_public | `null`      | String       | File path to public SSH key for LXC user, e.g. `~/.ssh/id_ed25519.pub`                           | no       |
| user_password       | `null`      | String       | Password for LXC user                                                                            | no       |
| start_on_boot       | `false`     | Boolean      | Start container on PVE boot                                                                      | no       |
| start_order         | `1`         | Number       | Start order                                                                                      | no       |
| start_delay         | `null`      | Number       | Startup delay in seconds                                                                         | no       |
| shutdown_delay      | `null`      | Number       | Shutdown delay in seconds                                                                        | no       |
| start_on_create     | `true`      | Boolean      | Start container after creation                                                                   | no       |

### Network Variables

| Variable     | Default | Type         | Description                                                           | Required |
| ------------ | ------- | ------------ | --------------------------------------------------------------------- | -------- |
| vnic_name    | `eth0`  | String       | Networking adapter name                                               | no       |
| vnic_bridge  | `vmbr0` | String       | Networking adapter bridge                                             | no       |
| vlan_tag     | `null`  | Number       | Network adapter VLAN tag                                              | no       |
| ipv4         |         | List(Object) | Defaults to DHCP, see example below for setting static IP and Gateway | no       |
| ipv4_address | `dhcp`  | String       | Defaults to DHCP, for static IPv4 address set CIDR                    | no       |
| ipv4_gateway | `null`  | String       | Defaults to DHCP, for static IPv4 gateway set IP address              | no       |

Example:

```HCL
module "lxc_static_ip_config" {
  source = "github.com/trfore/terraform-bpg-proxmox//modules/lxc"
  ...

  ipv4 = [
    {
      ipv4_address = "192.168.1.103/24"
      ipv4_gateway = "192.168.1.1"
    },
  ]
}
```

### Mount Point Object

| Variable     | Default | Type         | Description                                                                                                   | Required |
| ------------ | ------- | ------------ | ------------------------------------------------------------------------------------------------------------- | -------- |
| mountpoint   |         | List(Object) | Default will not create mount point, see example below for creating ones                                      | no       |
| mp_volume    | `null`  | String       | Storage name or host directory path ([bind mounts]) - e.g. `local-lvm`, `local-zfs`, `/mnt/pve/MY_NAS_SHARE`. | no       |
| mp_size      | `null`  | Number       | Size of the drive in GB. For [bind mounts] to local storage leave the default `null` value.                   | no       |
| mp_path      | `null`  | String       | Path within the container to mount the drive, e.g. `/mnt/storage`                                             | no       |
| mp_backup    | `false` | Boolean      |                                                                                                               | no       |
| mp_read_only | `false` | Boolean      |                                                                                                               | no       |

To add [bind mounts] for local drive(s), you must use root SSH access.

```HCL
# Required for Bind Mounts
provider "proxmox" {
  endpoint = var.pve_api_url
  insecure = true
  username = "root@pam"
  password = "MyRootPVEPassword"

  ssh {
    agent       = false
    private_key = file("${var.ssh_key_path}")
  }
}
```

Mount Example:

```HCL
# Create a new 4Gb drive using local storage
module "lxc_mountpoint_config" {
  source = "github.com/trfore/terraform-bpg-proxmox//modules/lxc"
  ...

  mountpoint = [
    {
      mp_volume    = "local-lvm"
      mp_size      = 4
      mp_path      = "/mnt/local"
      mp_backup    = true
    },
  ]
}

# Attach network storage (REQUIRES ROOT SSH ACCESS)
module "lxc_bind_network_storage" {
  source = "github.com/trfore/terraform-bpg-proxmox//modules/lxc"
  ...

  mountpoint = [
    {
      mp_volume = "/mnt/pve/nas-storage"
      mp_path   = "/mnt/nas-storage"
    },
  ]
}
```

## Outputs

## Examples

- [See example LXC configurations](../../examples/lxc/main.tf)

[terraform]: https://github.com/hashicorp/terraform
[bpg proxmox]: https://github.com/bpg/terraform-provider-proxmox
[bind mounts]: https://pve.proxmox.com/wiki/Unprivileged_LXC_containers#Using_local_directory_bind_mount_points
