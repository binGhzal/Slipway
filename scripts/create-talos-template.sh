#!/bin/bash
# ============================================================
# Create Talos Linux VM Template on Proxmox
# ============================================================
# This script downloads a Talos nocloud image from factory.talos.dev
# with required extensions and creates a Proxmox VM template.
#
# Prerequisites:
#   - SSH access to the Proxmox host
#   - .env file with TALOS_VERSION, TALOS_SCHEMATIC_ID, PROXMOX_TEMPLATE_ID
#
# Usage:
#   ./scripts/create-talos-template.sh
# ============================================================

set -euo pipefail

# Load environment
if [ -f .env ]; then
  set -a; source .env; set +a
else
  echo "ERROR: .env file not found. Run 'cp .env.example .env' and fill in values."
  exit 1
fi

TALOS_VERSION="${TALOS_VERSION:-v1.11.0}"
SCHEMATIC_ID="${TALOS_SCHEMATIC_ID}"
TEMPLATE_ID="${PROXMOX_TEMPLATE_ID:-9000}"
PROXMOX_NODE="${PROXMOX_NODE:-pve}"

if [ -z "${SCHEMATIC_ID}" ]; then
  echo "ERROR: TALOS_SCHEMATIC_ID not set."
  echo ""
  echo "Generate one at https://factory.talos.dev with these extensions:"
  echo "  - siderolabs/qemu-guest-agent  (required for CAPMOX)"
  echo "  - siderolabs/iscsi-tools       (required for Longhorn)"
  echo "  - siderolabs/util-linux-tools   (optional utilities)"
  echo ""
  echo "Or use the API:"
  echo '  curl -X POST --data-binary @- https://factory.talos.dev/schematics <<EOF'
  echo '  customization:'
  echo '    systemExtensions:'
  echo '      officialExtensions:'
  echo '        - siderolabs/qemu-guest-agent'
  echo '        - siderolabs/iscsi-tools'
  echo '        - siderolabs/util-linux-tools'
  echo '  EOF'
  exit 1
fi

IMAGE_URL="https://factory.talos.dev/image/${SCHEMATIC_ID}/${TALOS_VERSION}/nocloud-amd64.raw.xz"
IMAGE_FILE="nocloud-amd64.raw"

echo "=== Talos Template Creator ==="
echo "Talos Version:  ${TALOS_VERSION}"
echo "Schematic ID:   ${SCHEMATIC_ID}"
echo "Template VMID:  ${TEMPLATE_ID}"
echo "Proxmox Node:   ${PROXMOX_NODE}"
echo ""

echo "Downloading Talos nocloud image..."
wget -q --show-progress "${IMAGE_URL}" -O "${IMAGE_FILE}.xz"
echo "Decompressing..."
xz -d "${IMAGE_FILE}.xz"

echo ""
echo "=== Run the following commands on your Proxmox host ==="
echo ""
echo "# Upload the image to Proxmox"
echo "scp ${IMAGE_FILE} root@${PROXMOX_NODE}:/tmp/"
echo ""
echo "# SSH to Proxmox and create template"
echo "ssh root@${PROXMOX_NODE} << 'PROXMOX_EOF'"
echo "  qm create ${TEMPLATE_ID} --name talos-nocloud-template --memory 4096 --cores 2 --net0 virtio,bridge=vmbr0"
echo "  qm importdisk ${TEMPLATE_ID} /tmp/${IMAGE_FILE} local-lvm"
echo "  qm set ${TEMPLATE_ID} --scsihw virtio-scsi-pci --scsi0 local-lvm:vm-${TEMPLATE_ID}-disk-0"
echo "  qm set ${TEMPLATE_ID} --boot order=scsi0"
echo "  qm set ${TEMPLATE_ID} --ide2 local-lvm:cloudinit"
echo "  qm set ${TEMPLATE_ID} --serial0 socket --vga serial0"
echo "  qm set ${TEMPLATE_ID} --agent enabled=1"
echo "  qm template ${TEMPLATE_ID}"
echo "  rm /tmp/${IMAGE_FILE}"
echo "PROXMOX_EOF"
echo ""
echo "After running the above, your template is ready for CAPI."
