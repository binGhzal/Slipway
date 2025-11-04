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

# Installing Prerequisites

To maintain a single guide for both MacOS and Linux users, we will use `brew` as the package manager for installing the required dependencies.

## 1 - Install Homebrew:

If you don't have Homebrew installed, you can install it by running the following command in your terminal:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Then follow the on-screen instructions to complete the installation.

## 2 - Install required packages:

Once Homebrew is installed, you can use it to install the necessary packages by running the following commands:

```bash
brew install kind
brew install kubectl
brew install clusterctl
```

The `clusterctl` package includes the `cluster-api` command-line tool, which is essential for managing Kubernetes clusters using Cluster API. It helps you create, scale, and manage clusters across different infrastructure providers, while helping you avoid mis-configurations and it will ensuring best practices are followed.

# Create Configuration for Proxmox

To use `clusterctl` with Proxmox, you need to create a configuration file that specifies the Proxmox provider and other necessary settings.
This could be done by setting environment variables for the [Cluster API Provider for Proxmox VE](https://github.com/ionos-cloud/cluster-api-provider-proxmox/) (CAPMOX) and then generating a default configuration file using `clusterctl`.

First, create a new directory for your cluster configuration and navigate into it:

```bash
mkdir my-proxmox-cluster
cd my-proxmox-cluster
```

Next, set the required environment variables for CAPMOX. You can do this by creating a `capmox.env` file in your cluster configuration directory with the following content:

```env
## -- Controller settings -- ##
export PROXMOX_URL="https://<PROXMOX_NODE_IP>:8006"
export PROXMOX_TOKEN='capmox@pve!capi'
export PROXMOX_SECRET="<CAPMOX_API_TOKEN_SECRET>"
## -- Required workload cluster default settings -- ##
export PROXMOX_SOURCENODE="<PROXMOX_SOURCENODE_NAME>"
export TEMPLATE_VMID="<PROXMOX_BASE_IMAGE_TEMPLATE_ID>"
export ALLOWED_NODES="[<PROXMOX_NODE1_NAME>,<PROXMOX_NODE2_NAME>...]"
export VM_SSH_KEYS="ssh-rsa AAAAB..."
## -- networking configuration-- ##
export CONTROL_PLANE_ENDPOINT_IP="192.168.1.230"
export NODE_IP_RANGES="[192.168.1.220-192.168.1.229]"
export GATEWAY="192.168.1.1"
export IP_PREFIX="24"
export DNS_SERVERS="[192.168.1.1]"
export BRIDGE="vmbr0"
export NETWORK_MODEL="virtio"
## -- xl nodes -- ##
export BOOT_VOLUME_DEVICE="scsi0"
export BOOT_VOLUME_SIZE="20"
export NUM_SOCKETS="1"
export NUM_CORES="2"
export MEMORY_MIB="2048"
export FILE_STORAGE_FORMAT="raw"
export STORAGE_NODE="local-lvm"
export NODE_REPLICAS=1
export CLOUD_INIT_CONFIG="#cloud-config package_update=true packages=- net-tools"
```

Make sure to replace the placeholder values (e.g., `<PROXMOX_NODE_IP>`, `<CAPMOX_API_TOKEN_SECRET>`, etc.) with your actual Proxmox configuration details.
And make sure to replace all the `192.168.1.0/24` network values with the CIDR that matches your bridge interface `vmbr0` configuration, or whatever bridge you chose to use.

# Setting up a Management Cluster

First we will create a management cluster using `kind` (Kubernetes IN Docker). This cluster will host the Cluster API components and manage the lifecycle of your workload clusters on Proxmox. It run one or more Infrastructure Providers and will hold their resources (e.g., machines, clusters, etc.), as well as the Cluster API core components.

## 1 - Create the management cluster:

First create a new kind file named `kind-management-cluster.yaml`, with extra mount points to allow the docker provider to access the Docker socket of the host machine:

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
    ipFamily: dual
nodes:
    - role: control-plane
      extraMounts:
          - hostPath: /var/run/docker.sock
            containerPath: /var/run/docker.sock
```

Now create the management cluster using the kind configuration file:

```bash
kind create cluster --name=mgmt-capmox --config=kind-management-cluster.yaml
```

This command will create a new Kubernetes cluster named `mgmt-capmox` using the specified configuration file. The cluster will be set up with the necessary components to manage your Proxmox workload clusters.

The kubeconfig for the new cluster will be automatically updated, so you can start using `kubectl` to interact with it right away.

## 2 - Convert the kind cluster to a management cluster:

Now that you have your kind cluster up and running, you can convert it into a management cluster using `clusterctl`. Run the following command:

```bash
source capmox.env
clusterctl init --infrastructure proxmox --ipam in-cluster
```

This command initializes the management cluster with the Proxmox infrastructure provider and sets up in-cluster IP address management (IPAM) with Kubeadm as control-plane and bootstrap provider.

This process may take a few minutes as it pulls the necessary images and deploys the required components. Once completed, your management cluster will be ready to manage Proxmox workload clusters.

# Creating a Cluster with Kubeadm CP Provider

With the management cluster set up, you can now create a workload cluster on Proxmox using the Kubeadm Control Plane (KCP) provider. This provider allows you to create and manage Kubernetes clusters using Kubeadm.

The kubeadm control plane provider watches `KubeadmControlPlane` resources and creates the necessary control plane machines for your workload cluster. When you create a `KubeadmControlPlane` resource, the provider will generate a corresponding set of `Machine` objects (defined by a `MachineTemplate` file) and orchestrates `kubeadm init` to bootstrap the first control plane node, and `kubeadm join` for any additional control plane nodes.

KubeadmControlPlane also manages certificate rotation for the control plane nodes, and will handle rolling updates when you change the configuration of the control plane. All configurations (e.g., Kubernetes version, number of replicas, etc.) are defined declaratively in the `KubeadmControlPlane` resource, and the provider ensures that the actual state of the control plane matches the desired state specified in the resource.

When you create a workload cluster using the Kubeadm Control Plane with the Proxmox provider, the provider will create the necessary Proxmox VMs based on the `MachineTemplate` you provide. It will also configure the VMs with the appropriate networking settings, SSH keys, and other configurations specified in your cluster configuration.

## 1 - Generate manifest files:

First, generate the manifest files for your workload cluster using `clusterctl`. Run the following command:

```bash
source capmox.env

clusterctl generate cluster capi \
  --kubernetes-version v1.33.0 \
  --control-plane-machine-count=1 \
  --worker-machine-count=2  > capi.yaml
```

This command generates a set of YAML manifest files that define the resources needed to create a workload cluster named `capi` with one control plane node and one worker node, running Kubernetes version `v1.33.0`. The generated manifests will be saved in a file named `capi.yaml`.

## 2 - Modify the generated manifest files:

You will need to make some adjustments to the generated manifest files to ensure they are correctly configured for your Proxmox environment.

Change the `cidrBlock` to a value that does not overlap with the CIDR block of your Proxmox bridge network (e.g., `vmbr0`). Otherwise your Pods will provide IPs that conflict with the Proxmox host network, causing connectivity issues.
in my case I changed it to a value that is not used anywhere else in my network either in a physical or virtual lan: `10.44.0.0/16`

```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
    name: capi
    namespace: default
spec:
    clusterNetwork:
        pods:
            cidrBlocks:
                - 10.44.0.0/16
    controlPlaneRef:
        apiVersion: controlplane.cluster.x-k8s.io/v1beta1
        kind: KubeadmControlPlane
        name: capi-control-plane
    infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1alpha1
        kind: ProxmoxCluster
        name: capi
```

## 3 - Allow the proxmox provider to overcommit memory:

By default, the Proxmox provider does not allow memory overcommitment, which can lead to issues when creating multiple VMs on a Proxmox node with limited physical memory. To enable memory overcommitment, you need to modify the `ProxmoxCluster` resource in the generated manifest file. If you haven’t explicitly enabled memory over-commit, this check will refuse to schedule any VM once the sum of allocated RAM reaches the node’s nominal capacity, causing new VMs to hang even though Proxmox could safely oversubscribe memory in practice.

The following property needs to be added under the `spec` section of the `ProxmoxCluster` resource:

```yaml
schedulerHints:
    memoryAdjustment: 0
```

Here is an example of how the modified `ProxmoxCluster` resource should look:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha1
kind: ProxmoxCluster
metadata:
    name: capi
    namespace: default
spec:
    allowedNodes:
        - pve-am06
    controlPlaneEndpoint:
        host: 10.0.13.230
        port: 6443
    dnsServers:
        - 10.0.0.1
    ipv4Config:
        addresses:
            - 10.0.13.220-10.0.13.229
        gateway: 10.0.0.1
        prefix: 20
    schedulerHints:
        memoryAdjustment: 0
```

## 4 - Apply the manifest files to create the cluster:

```bash
kubectl apply -f capi.yaml
```

This command applies the generated manifest files to your management cluster, which will create the necessary resources to provision a workload cluster on Proxmox. The Proxmox provider will then create the VMs for the control plane and worker nodes based on the configurations specified in the manifest file.

Sit back and watch the magic happen! as this process may take several minutes to complete, depending on your Proxmox environment and network speed.

To monitor the progress of the cluster creation, you can use the following command to check the status of the machines:

```bash
kubectl get machines
```

and to get more detailed information about the cluster and its components, you can describe the cluster resource:

```bash
kubectl describe cluster capi
```

After the process is done, get the kubeconfig for your new workload cluster by running the following command:

```bash
clusterctl get kubeconfig capi > capi.kubeconfig
```

## 5 - Deploy a CNI plugin:

To enable networking for your workload cluster, you need to deploy a Container Network Interface (CNI) plugin. A popular choice is Calico, which provides robust networking and network policy features which we will use in our guide, but you can choose any CNI plugin that suits your needs.

```bash
kubectl --kubeconfig=./capi.kubeconfig \
  apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml
```

This command deploys Calico to your workload cluster using the kubeconfig file you obtained earlier. Once Calico is deployed, your cluster will have networking capabilities, allowing pods to communicate with each other and with external networks.

after theat we can check if all nodes are running by using this command:

```bash
kubectl --kubeconfig=./capi.kubeconfig get nodes -w
```

## 6 - Move ClusterAPI objects to the workload cluster:

After the workload cluster is up and running, it's a good practice to move the ClusterAPI objects from the management cluster to the workload cluster. This helps in better management and organization of resources.

After all nodes are ready, we will first initialize the management cluster's control over the workload cluster by running:

```bash
clusterctl init --infrastructure proxmox --ipam in-cluster --kubeconfig ./capi.kubeconfig
```

After that, we can move the ClusterAPI objects by executing the following command:

```bash
clusterctl move --to-kubeconfig ./capi.kubeconfig
```

Once the move is complete, the management cluster will no longer have control over the workload cluster, and all ClusterAPI resources will be managed directly within the workload cluster itself.
So we will delete the kind management cluster because we dont need it anymore:

```bash
kind delete cluster --name=mgmt-capmox
```

## 7 - Add bootstrap cluster to your kubeconfig:

To easily manage your workload cluster using `kubectl`, you can add the kubeconfig of the workload cluster to your existing kubeconfig file. This allows you to switch between clusters seamlessly.
Run the following command to merge the workload cluster's kubeconfig into your existing kubeconfig:

```bash
cp ~/.kube/config ~/.kube/config-bck
export KUBECONFIG=~/.kube/config:./capi.kubeconfig
kubectl config view
kubectl config view --flatten > one-config.yaml
mv one-config.yaml ~/.kube/config
```

This command copies your existing kubeconfig to a backup file, merges the workload cluster's kubeconfig with your existing kubeconfig, and then flattens the combined configuration into a single file. The final merged kubeconfig is saved back to `~/.kube/config`.

## 8 - Create additional clusters as needed:

You can repeat the process of generating manifest files, modifying them, and applying them to create additional workload clusters on Proxmox as needed. Just make sure to use unique names for each cluster and adjust the configurations accordingly.

```bash
clusterctl generate cluster <cluster-name> \
--kubernetes-version v1.33.0 \
--control-plane-machine-count=1 \
--worker-machine-count=2 \
> <cluster-name>.yaml
kubetl apply -f <cluster-name>.yaml
```

and whenver you need to remove a cluster, just delete the cluster resource using kubectl:

```bash
kubectl delete cluster <cluster-name>
```
