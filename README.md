# m-azure-kubernetes-service

Epiphany Module: Azure Kubernetes Service

AzKS module is responsible for providing AKS cluster. 

# Basic usage

## Requirements

Requirements are listed in a separate [document](docs/REQUIREMENTS.md).

## Run module

* Create a shared directory:

  ```shell
  mkdir /tmp/shared
  ```

  This 'shared' dir is a place where all configs and states will be stored while working with Epiphany modules.

* Generate ssh keys in: /tmp/shared/vms_rsa.pub

  ```shell
  ssh-keygen -t rsa -b 4096 -f /tmp/shared/vms_rsa -N ''
  ```

* [Optional] Build Docker image if development version is used

  ```shell
  make build
  ```

* Initialize AzKS module:

  This module requires that existing resource group, vnet and subnet are already created (preferably with [AzBI](https://github.com/epiphany-platform/m-azure-basic-infrastructure) module). If those resources are already provided by AzBI module (and there is at least one empty subnet provided by AzBI module), initialization step can be started without any parameters and module will extract configuration from state file.  

  ```shell
  docker run --rm -v /tmp/shared:/shared epiphanyplatform/azks:dev init
  ```
  
  If mentioned resources were created some other way (i.e.: manually) there is possibility to provide additional parameters: 

  ```shell
  docker run --rm -v /tmp/shared:/shared epiphanyplatform/azks:dev init --name=azks-modules-test  --subnet_name=some-subnet --vnet_name=some-vnet --rg_name=some-rg
  ```
  
  For detailed list of init parameters go to [INPUTS](docs/INPUTS.adoc) documentation. 

  :star: Variable values can be passed as docker environment variables as well. In presented example we could use `docker run` command `-e NAME=azks-modules-test` parameter instead of `--name=azks-modules-test` command parameter.

  :warning: Use image's tag according to tag generated in build step.

  This command will create configuration file of AzKS module in /tmp/shared/azks/azks-config.yml. You can investigate what is stored in that file. Available parameters are described in the [inputs](docs/INPUTS.adoc) document.

  :warning: Pay attention to the docker image tag you are using. Command `make build` command uses a specific version
  tag (default `epiphanyplatrofm/azks:dev`).

* Plan and apply AzBI module:

  ```shell
  docker run --rm -v /tmp/shared:/shared -e SUBSCRIPTION_ID=subscriptionId -e CLIENT_ID=appId -e CLIENT_SECRET=password -e TENANT_ID=tenantId epiphanyplatform/azks:dev plan
  docker run --rm -v /tmp/shared:/shared -e SUBSCRIPTION_ID=subscriptionId -e CLIENT_ID=appId -e CLIENT_SECRET=password -e TENANT_ID=tenantId epiphanyplatform/azks:dev apply
  ```
  :star: Variable values can be passed as docker environment variables. I's often more convenient to pass sensitive values as presented.

  Running those commands should create AKS cluster in provided RG. You should verify in Azure Portal.

* Destroy module resources:

  ```shell
  docker run --rm -v /tmp/shared:/shared -e SUBSCRIPTION_ID=subscriptionId -e CLIENT_ID=appId -e CLIENT_SECRET=password -e TENANT_ID=tenantId epiphanyplatform/azks:dev plan --destroy
  docker run --rm -v /tmp/shared:/shared -e SUBSCRIPTION_ID=subscriptionId -e CLIENT_ID=appId -e CLIENT_SECRET=password -e TENANT_ID=tenantId epiphanyplatform/azks:dev destroy
  ```
  :star: Variable values can be passed as docker environment variables. I's often more convenient to pass sensitive values as presented.

  :warning: Running those commands will remove AKS resource and all application deployed on it so be careful. You should verify in Azure Portal.

# AzKS output data

The output from this module is just kubeconfig. To extract it from state one can use [JQ](https://stedolan.github.io/jq) combined with `output` command: `docker run --rm -v $(pwd)/shared:/shared epiphanyplatform/azks:dev output | grep -v DEPRECATION | jq -r '.azks.output.kubeconfig'`. :warning: `grep -v DEPRECATION` part is filtering additional possible notification lines. 

# Examples

For examples running description please have a look into [this document](docs/EXAMPLES.md).

# Development

For development related topics please look into [this document](docs/DEVELOPMENT.md).

# Azure limits

There are subscription or regional limits in Azure Cloud. All of them can be investigated [on this site](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/azure-subscription-service-limits). For this module most important limitation is kubernetes version operating in a region of your choice. To check version `az aks get-versions --location <your region>`. To check your current resources usage please go to Azure Portal > Subscriptions > select subscription > Usage + quotas. 

# Module dependencies

| Component                 | Version | Repo/Website                                          | License                                                           |
| ------------------------- | ------- | ----------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                 | 0.13.2  | https://www.terraform.io/                             | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AzureRM provider | 2.27.0 | https://github.com/terraform-providers/terraform-provider-azurerm | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/LICENSE) |
