# Preparing Proxmox user and API tokens

To interact with the Proxmox API for image building, you need to create a dedicated user and generate API tokens. Execute the following commands on one of your Proxmox nodes:

```bash
# Create a new user for image building
pveum user add image-builder@pve
# Assign the PVEAdmin role to the new user
pveum aclmod / -user image-builder@pve -role PVEAdmin
# Generate an API token for the image-builder user
pveum user token add image-builder@pve capi -privsep 0
```

Make sure to securely store the generated API token, as it will be required for authentication when using the Proxmox API for image building tasks.

The above commands will create a user named `image-builder@pve`, assign it the `PVEAdmin` role, and generate an API token named `capi`. You can use this token to authenticate your API requests when building images.

Also create another user and token to use for cluster-api operations:

```bash
# Create a new user for cluster-api operations
pveum user token add image-builder@pve capi -privsep 0
# Assign the PVEAdmin role to the new user
pveum aclmod / -user capmox@pve -role PVEAdmin
# Generate an API token for the cluster-api user
pveum user token add capmox@pve capi -privsep 0
```

This will create a user named `capmox@pve`, assign it the `PVEAdmin` role, and generate an API token named `capi` for cluster-api operations.

Make sure to securely store the generated API token, as it will be required for authentication when using the Proxmox API for cluster-api tasks.

# Building Images for Proxmox with image-builder

## 1 - Clone the image-builder repository:

```bash
git clone git@github.com:kubernetes-sigs/image-builder.git
cd image-builder/images/capi
```

## 2 - Install dependencies:

```bash
make deps
```

this will install packer and other dependencies required for building images in the folder `image-builder/images/capi/.bin` if theyre not already installed on your system. After installation you can add this folder to your PATH variable to use the installed dependencies globally.

```bash
export PATH=$PWD/.bin:$PATH
```

## 3 - Configure your Proxmox credentials:

Add a new `.env` file in the `image-builder/images/capi/packer/proxmox/` directory with the following content:

```env
export PROXMOX_URL="https://<PROXMOX_NODE_IP>:8006/api2/json"
export PROXMOX_USERNAME=image-builder@pve!capi
export PROXMOX_TOKEN=<IMAGE_BUILD_API_TOKEN>
export PROXMOX_NODE="<PROXMOX_NODE_NAME>"
export PROXMOX_ISO_POOL="local"
export PROXMOX_STORAGE_POOL="<PROXMOX_STORAGE_POOL>"
export DISK_FORMAT="<DISK_FORMAT>"
export PROXMOX_BRIDGE="vmbr0"
export PROXMOX_NIC_MODEL="virtio"
# If you are working with OPNSense, image-builder is not
# assigning a Proxmox MAC address to the VM, and this seems
# to be an issue with DHCP in VXLANs controlled by OPNSense!
# In theory you can configure your MTU and VLAN tag with the
# following variables. For the time being leave them commented.
# export PROXMOX_MTU="1450"
# export PROXMOX_VLAN="23"
export PACKER_FLAGS="--var memory=2048 --var 'kubernetes_rpm_version=1.32.5' --var 'kubernetes_semver=v1.32.5' --var 'kubernetes_series=v1.32' --var 'kubernetes_deb_version=1.32.5-1.1'"
```

Make sure to add the `.env` file to your `.gitignore` to avoid committing sensitive information if you have forked the repository and are using version control.

The environment variable `PROXMOX_ISO_POOL`, `PROXMO_BRIDGE`, `PROXMOX_STORAGE_POOL` are optional, theyre just here to demonstrate how to override default values, because we might need to customize some later on for a multi-node proxmox cluster setup.

If you want to specify a specific Kubernetes version, you can do so by adding the `PACKER_FLAGS` variable as shown above. Adjust the version numbers as needed. You can ignore this, but its not recommended unless you fully understand how image-builder uses Packer to build images behind the scenes.

## 4 - Configure `PROXMOX_STORAGE_POOL` and `DISK_FORMAT` variables:

If you are running on a **single Proxmox node** setup, you need to set the `PROXMOX_STORAGE_POOL` value to `local-lvm`, and because `local-lvm` only supports `raw` disk format, set the `DISK_FORMAT` variable to `raw`.

Otherwise, if you are running on a **multi-node Proxmox cluster** setup, you can set the `PROXMOX_STORAGE_POOL` value to any shared storage pool available in your cluster (e.g., `nfs-shared`, `ceph-shared`, etc.), and you can set the `DISK_FORMAT` variable to `qcow2` or `raw`, depending on your storage backend capabilities.

## 5 - Adjust the Packer template:

We will need to make some adjustments to the Packer template, as it is missing some details that are very handy when working with Proxmox.

Some of these variable are just missing from the template, while others just specify hardware defaults the we want to adjust for our base line image template.

open the file `image-builder/images/capi/packer/proxmox/packer.json.tmpl` and make the following changes:

- Add the following variables to the `variables` section:

```json
"nic_model": "{{env `PROXMOX_NIC_MODEL`}}",
```

- Also in the `variables` section, make sure to adjust the following variables to match the environment variables we set in the `.env` file: (- indicates existing lines, + indicates new lines)

```json
- "disk_format": "qcow2",
+ "disk_format": "{{env `DISK_FORMAT`}}",
- "sockets": "2",
+ "sockets": "1",
```

- In the `network_adapter` section, add a new line to specify the NIC model:

```json
"network_adapters": [
        {
          "model": "{{user `nic_model`}}",
          "bridge": "{{user `bridge`}}",
          "mtu": "{{ user `mtu` }}",
          "vlan_tag": "{{user `vlan_tag`}}"
        }
 ],
```

## 6 - Build the image:

First, make sure you have sourced the environment variables:

```bash
source packer/proxmox/.env
```

Now you are ready to build the image! Run one the following command from the `image-builder/images/capi/` directory to start the build process:

```bash
make build-proxmox-ubuntu-2204
```

This will start the Packer build process, which will create a new VM in Proxmox based on the specified OS template. The build process may take some time(~25-40 minutes), depending on your network speed and Proxmox node performance. Once the build is complete, you should see a new VM template in your Proxmox storage pool.

I recommend monitoring the build process through the Proxmox web interface, as it provides a more user-friendly way to track the VM creation and configuration steps.