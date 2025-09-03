# Packer Proxmox

## Installing Proxmox plugin

You have two options:

- You can add this config block to your packer.pkr.hcl file and run `packer init`.

```hcl
packer {
  required_plugins {
    name = {
      version = "~> 1"
      source  = "github.com/hashicorp/proxmox"
    }
  }
}
```

- Run `packer plugins install github.com/hashicorp/proxmox` to install the plugin globally in packer.

## Running Packer

1. Fill in your Proxmox API credentials in the `variables.pkrvars.hcl` file.
2. Navigate into the folder you want to create a template with, and edit the Packer template file (e.g., `ubuntu-server-noble.pkr.hcl`) and `http/user-data` file.
3. Run `packer build -var-file=../variables.pkrvars.hcl .`

## Troubleshooting

- If you have tailscale installed, be aware that packer could grab the IP of your tailscale adapter rather than your LAN. You can either hard code the IP in the boot command or try setting the `http_interface` option
- Sometimes the boot command is typed too fast and can cause issues. You can increase the time between types by using the `boot_key_interval` option.
