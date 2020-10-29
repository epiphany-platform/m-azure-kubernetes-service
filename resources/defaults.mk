M_NAME ?= epiphany

# create aks in existing subnet
M_EXISTING_RG_NAME ?= unset
M_EXISTING_SUBNET_ID ?= unset

# create aks in azbi
ifeq ($(M_EXISTING_RG_NAME),unset)
M_RG_NAME ?= $(M_NAME)-rg
M_VNET_NAME ?= $(M_NAME)-vnet
M_ADDRESS_PREFIX ?= 10.0.0.0/16
else
M_RG_NAME ?= unset
M_VNET_NAME ?= unset
M_ADDRESS_PREFIX ?= unset
endif

# default node pool
M_SIZE ?= 2
M_MIN ?= 2
M_MAX ?= 5
M_VM_SIZE ?= Standard_DS2_v2
M_DISK_SIZE ?= 36
M_AUTO_SCALING ?= true
M_VMS_RSA ?= vms_rsa

# azure credentials
M_ARM_CLIENT_ID ?= unset
M_ARM_CLIENT_SECRET ?= unset
M_ARM_SUBSCRIPTION_ID ?= unset
M_ARM_TENANT_ID ?= unset
