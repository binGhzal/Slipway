# Talos Image Management

## Overview

Talos Linux uses **nocloud** images for Proxmox deployments. Images are built via [Talos Image Factory](https://factory.talos.dev) with a deterministic schematic ID that ensures reproducible builds.

## Required Extensions

| Extension | Why |
|-----------|-----|
| `siderolabs/qemu-guest-agent` | **Mandatory.** CAPMOX uses this to detect VM IP addresses and readiness |
| `siderolabs/iscsi-tools` | Required for Longhorn distributed storage |
| `siderolabs/util-linux-tools` | Useful system utilities |

## Creating a Schematic

### Via Web UI
1. Go to [factory.talos.dev](https://factory.talos.dev)
2. Select your Talos version
3. Select platform: **nocloud**
4. Add the extensions listed above
5. Note the **Schematic ID**
6. Download `nocloud-amd64.raw.xz`

### Via API
```bash
SCHEMATIC_ID=$(curl -sX POST --data-binary @- https://factory.talos.dev/schematics <<'EOF'
customization:
  systemExtensions:
    officialExtensions:
      - siderolabs/qemu-guest-agent
      - siderolabs/iscsi-tools
      - siderolabs/util-linux-tools
EOF
)
echo "Schematic ID: ${SCHEMATIC_ID}"
```

## Creating the Proxmox Template

Use the helper script:
```bash
./scripts/create-talos-template.sh
```

Or manually on the Proxmox host:
```bash
# Download image
wget https://factory.talos.dev/image/<SCHEMATIC_ID>/<TALOS_VERSION>/nocloud-amd64.raw.xz
xz -d nocloud-amd64.raw.xz

# Create template VM
qm create 9000 --name talos-nocloud-template --memory 4096 --cores 2 --net0 virtio,bridge=vmbr0
qm importdisk 9000 nocloud-amd64.raw local-lvm
qm set 9000 --scsihw virtio-scsi-pci --scsi0 local-lvm:vm-9000-disk-0
qm set 9000 --boot order=scsi0
qm set 9000 --ide2 local-lvm:cloudinit
qm set 9000 --serial0 socket --vga serial0
qm set 9000 --agent enabled=1
qm template 9000
```

## Upgrading Talos

When a new Talos version is released:
1. Generate a new image using the same schematic ID (extensions stay the same)
2. Create a new template on Proxmox (or update existing)
3. Update `TALOS_VERSION` in `.env`
4. Update `version` fields in CAPI TalosControlPlane and MachineDeployment manifests
5. Commit - CAPI performs rolling updates

## Current Configuration

Save your schematic ID here for reference:
- **Schematic ID:** `<fill in after creation>`
- **Talos Version:** See `.env`
- **Template VMID:** 9000
