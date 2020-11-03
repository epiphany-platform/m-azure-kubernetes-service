# m-azure-kubernetes-service

Epiphany Module: Azure Kubernetes Service

## Prepare service principal

Have a look [here](https://www.terraform.io/docs/providers/azurerm/guides/service_principal_client_secret.html).

```shell
az login
az account list #get subscription from id field
az account set --subscription="SUBSCRIPTION_ID"
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/SUBSCRIPTION_ID" --name="SOME_MEANINGFUL_NAME" #get appID, password, tenant, name and displayName
```

## Run module

The AzKS cluster and new subnet will be created in the resource group and vnet from the [AzBI Module](https://github.com/epiphany-platform/m-azure-basic-infrastructure) or you can create the AzKS cluster in an already existing subnet.

* Initialize the AzKS module in [AzBI Module](https://github.com/epiphany-platform/m-azure-basic-infrastructure) (without parameters it will extract some knowledge from the status file):
  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/azks:latest init
  ```
  or initialize the AzKS module with some parameters (ie.: in already existing subnet by setting `M_SUBNET_NAME`):
  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/azks:latest init M_RG_NAME="demo-ropu-rg" M_VNET_NAME="demo-ropu-vnet" M_SUBNET_NAME="demo-ropu-kubernetes-master-subnet-0"
  ```
  The previous command created a configuration file of the AzKS module in `/tmp/shared/azks/azks-config.yml`. You can investigate what is stored in that file and change it at will.
* Plan and apply AzKS module:
  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/azks:latest plan M_ARM_CLIENT_ID=appId M_ARM_CLIENT_SECRET=password M_ARM_SUBSCRIPTION_ID=subscriptionId M_ARM_TENANT_ID=tenantId
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/azks:latest apply M_ARM_CLIENT_ID=appId M_ARM_CLIENT_SECRET=password M_ARM_SUBSCRIPTION_ID=subscriptionId M_ARM_TENANT_ID=tenantId
  ```
  Running those commands should create the AzKS service. You should verify in Azure Portal.
* Extract kubeconfig following way:
  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/azks:latest kubeconfig
  ```
  This command will create file `/tmp/shared/kubeconfig`. You will need to move this file manually to `/tmp/shared/build/your-cluster-name/kubeconfig`.

## Build image

In main directory run:

```shell
make build
```

## Run example

You can run example of AzBI and AzKS module using files in the `examples` directory.

```shell
cd examples/basic_flow
ARM_CLIENT_ID="appId field" ARM_CLIENT_SECRET="password field" ARM_SUBSCRIPTION_ID="id field" ARM_TENANT_ID="tenant field" make all
```

Or use config file with credentials:

```shell
cd examples/basic_flow
cat >azure.mk <<'EOF'
ARM_CLIENT_ID ?= "appId field"
ARM_CLIENT_SECRET ?= "password field"
ARM_SUBSCRIPTION_ID ?= "id field"
ARM_TENANT_ID ?= "tenant field"
EOF
make all
```

## Run example in existing subnet

You can run example of AzKS in existing subnet using files in `examples` directory.

```shell
cd examples/create_in_existing_subnet
ARM_CLIENT_ID="appId field" ARM_CLIENT_SECRET="password field" ARM_SUBSCRIPTION_ID="id field" ARM_TENANT_ID="tenant field" M_RG_NAME="existing rg name" M_SUBNET_NAME="existing subnet name" M_VNET_NAME="existing vnet name" EXISTING_SUBNET="true" make all
```

Or use config file with credentials:

```shell
cd examples/create_in_existing_subnet
cat >azure.mk <<'EOF'
ARM_CLIENT_ID ?= "appId field"
ARM_CLIENT_SECRET ?= "password field"
ARM_SUBSCRIPTION_ID ?= "id field"
ARM_TENANT_ID ?= "tenant field"
M_RG_NAME ?= "existing rg name"
M_SUBNET_NAME ?= "existing subnet name"
M_VNET_NAME ?= "existing vnet name"
EOF
make apply
```

If You want to destroy the AzKS, execute above instruction in the same way using `destroy` command instead of `all`.

## Release module

```shell
make release
```

or if you want to set different version number:

```shell
make release VERSION=number_of_your_choice
```

## Run tests

```
make test
```

## Run tests in Kubernetes based build system

Kubernetes based build system means that build agents work inside Kubernetes cluster. During testing process application runs inside docker container. This means that we've got "docker inside docker (DiD)". This kind of environment requires a bit different configuration of mount shared storage to docker container than with standard one-layer configuration.

With DiD configuration shared volume needs to be created on host machine and this volume is shared with application container as Kubernetes volume.
Configuration steps:

1.  Create volume  (host path). In deployment.yaml add this config to create kubernetes volume:

```
volumes:
- name: tests-share
  hostPath:
    path: /tmp/tests-share
```

See manual for more details: https://kubernetes.io/docs/concepts/storage/volumes/#hostpath

2. Add mount point for kubernetes pod (agent). In deployment.yaml add this config to define volume's mount point:

```
volumeMounts:
- mountPath: /tests-share
  name: tests-share
```

3. Inside pod where tests will run set two variables to indicate host path and mount point:

```
export K8S_HOST_PATH=/tests-share
export K8S_VOL_PATH=/tmp/tests-share  ##modify paths according your needs, but they need to match paths from steps 1 and 2.
```

4. Go to location where you downloaded repository and run:

```
make test
```

5. Test results will be availabe inside ```/tests-share``` on pod on which tests are running and is mapped to ```/tmp/tests-share``` on kubernetes node.

## Input parameters

To check supported module parameters list navigate to [inputs](docs/INPUTS.md) document.

## Windows users

This module is designed for Linux/Unix development/usage only. If you need to develop from Windows you can use the included [devcontainer setup for VScode](https://code.visualstudio.com/docs/remote/containers-tutorial) and run the examples the same way but then from then ```examples/basic_flow_devcontainer``` folder.

## Module dependencies

| Component                 | Version | Repo/Website                                          | License                                                           |
| ------------------------- | ------- | ----------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                 | 0.13.2  | https://www.terraform.io/                             | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AzureRM provider | 2.27.0 | https://github.com/terraform-providers/terraform-provider-azurerm | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/LICENSE) |
| Make                      | 4.3     | https://www.gnu.org/software/make/                    | [ GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html) |
| yq                        | 3.3.4   | https://github.com/mikefarah/yq/                      | [ MIT License](https://github.com/mikefarah/yq/blob/master/LICENSE) |
