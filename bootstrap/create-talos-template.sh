#!/bin/bash
set -euo pipefail

# ═══════════════════════════════════════════════════════════════
# Talos VM Template Builder for Proxmox
# Creates a Proxmox VM template from a Talos nocloud disk image
# ═══════════════════════════════════════════════════════════════

TALOS_VERSION="${TALOS_VERSION:-v1.9.5}"
VMID="${VMID:-9000}"
STORAGE="${STORAGE:-bigdisk}"
EFI_STORAGE="${EFI_STORAGE:-local-lvm}"
BRIDGE="${BRIDGE:-vmbr0}"
VLAN="${VLAN:-2}"

# Talos Image Factory schematic ID
# Default includes: qemu-guest-agent, iscsi-tools
# Generate custom at https://factory.talos.dev/
SCHEMATIC_ID="${SCHEMATIC_ID:-376567988ad370138ad8b2698212367b8edcb69b5fd68c80be1f2ec7d603b4ba}"

IMAGE_URL="https://factory.talos.dev/image/${SCHEMATIC_ID}/${TALOS_VERSION}/nocloud-amd64.raw.xz"

echo "═══════════════════════════════════════════════════════════"
echo "  Talos Template Builder"
echo "  Version: ${TALOS_VERSION}"
echo "  VMID:    ${VMID}"
echo "  Storage: ${STORAGE}"
echo "  Network: ${BRIDGE} (VLAN ${VLAN})"
echo "═══════════════════════════════════════════════════════════"

# 1. Download the Talos nocloud image
echo "==> Downloading Talos ${TALOS_VERSION} nocloud image..."
wget -q --show-progress -O /tmp/talos-nocloud.raw.xz "${IMAGE_URL}"
echo "==> Decompressing..."
xz -d /tmp/talos-nocloud.raw.xz

# 2. Destroy existing VM if present
echo "==> Removing existing VM ${VMID} if present..."
qm destroy ${VMID} --purge 2>/dev/null || true

# 3. Create a new VM shell
echo "==> Creating VM ${VMID}..."
qm create ${VMID} \
  --name "talos-${TALOS_VERSION}" \
  --memory 4096 \
  --cores 2 \
  --sockets 1 \
  --cpu host \
  --net0 "virtio,bridge=${BRIDGE},tag=${VLAN}" \
  --bios ovmf \
  --machine q35 \
  --scsihw virtio-scsi-single \
  --agent enabled=1 \
  --ostype l26 \
  --serial0 socket \
  --vga serial0

# 4. Add EFI disk (UEFI without Secure Boot)
echo "==> Adding EFI disk..."
qm set ${VMID} --efidisk0 "${EFI_STORAGE}:0,efitype=4m,pre-enrolled-keys=0"

# 5. Import the nocloud disk image as the boot disk
echo "==> Importing nocloud disk to ${STORAGE}..."
qm importdisk ${VMID} /tmp/talos-nocloud.raw ${STORAGE}

# 6. Attach the imported disk as scsi0
echo "==> Attaching disk as scsi0..."
qm set ${VMID} \
  --scsi0 "${STORAGE}:vm-${VMID}-disk-0,discard=on,iothread=1,ssd=1"

# 7. Add a cloud-init drive on ide2
# CAPMOX will replace this with its generated ISO containing the Talos machine config
echo "==> Adding cloud-init drive on ide2..."
qm set ${VMID} --ide2 "${EFI_STORAGE}:cloudinit"

# 8. Set boot order: disk first (after Talos installs), then cloud-init is just data
echo "==> Setting boot order..."
qm set ${VMID} --boot order=scsi0

# 9. Convert to template
echo "==> Converting to template..."
qm template ${VMID}

# 10. Clean up
echo "==> Cleaning up..."
rm -f /tmp/talos-nocloud.raw

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Template ready!"
echo "  VMID:    ${VMID}"
echo "  Name:    talos-${TALOS_VERSION}"
echo "  Storage: ${STORAGE}"
echo "  Network: ${BRIDGE} (VLAN ${VLAN})"
echo "  UEFI:    Yes (no Secure Boot)"
echo "  Cloud-init drive: ide2 (CAPMOX will replace with machine config)"
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Next: run 'task mgmt:apply' to provision the management cluster"
