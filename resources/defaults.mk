M_NAME ?= epiphany

ifneq ($(M_EXISTING_SUBNET),true)
M_RG_NAME ?= $(M_NAME)-rg
M_VNET_NAME ?=$(M_NAME)-vnet
M_ADDRESS_PREFIX ?= 10.0.0.0/16
M_SUBNET_NAME ?= $(M_NAME)-snet-aks
M_EXISTING_SUBNET ?= false
else
M_RG_NAME ?= unset
M_VNET_NAME ?= unset
M_ADDRESS_PREFIX ?= unknown
M_SUBNET_NAME ?= unset
M_EXISTING_SUBNET ?= true
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
