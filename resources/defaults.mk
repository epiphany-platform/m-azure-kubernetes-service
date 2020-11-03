M_NAME ?= epiphany

ifeq ($(M_SUBNET_NAME),)
M_RG_NAME ?= $(M_NAME)-rg
M_VNET_NAME ?= $(M_NAME)-vnet
M_ADDRESS_PREFIX ?= 10.0.0.0/16
M_SUBNET_NAME ?= unset
else
M_ADDRESS_PREFIX ?= unset
endif

define _M_DEFAULT_NODE_POOL
{
  size: 2,
  min: 2,
  max: 5,
  vm_size: Standard_DS2_v2,
  disk_size: 36,
  auto_scaling: true,
  type: VirtualMachineScaleSets
}
endef

define _M_AUTO_SCALER_PROFILE
{
  balance_similar_node_groups: false,
  max_graceful_termination_sec: 600,
  scale_down_delay_after_add: 10m,
  scale_down_delay_after_delete: 10s,
  scale_down_delay_after_failure: 10m,
  scan_interval: 10s,
  scale_down_unneeded: 10m,
  scale_down_unready: 10m,
  scale_down_utilization_threshold: 0.5
}
endef

M_DEFAULT_NODE_POOL ?= $(_M_DEFAULT_NODE_POOL)
M_AUTO_SCALER_PROFILE ?= $(_M_AUTO_SCALER_PROFILE)

# azure credentials
M_ARM_CLIENT_ID ?= unset
M_ARM_CLIENT_SECRET ?= unset
M_ARM_SUBSCRIPTION_ID ?= unset
M_ARM_TENANT_ID ?= unset

# other parameters
M_K8S_VERSION ?= 1.18.8
M_PUBLIC_IP_ENABLED ?= false
M_RBAC_ENABLED ?= false
M_AZURE_AD ?= null
M_K8S_DASHBOARD_ENABLED ?= true
M_ADMIN_USERNAME ?= operations
M_VMS_RSA ?= vms_rsa
