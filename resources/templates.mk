define M_METADATA_CONTENT
labels:
  version: $(M_VERSION)
  name: Azure Kubernetes Service
  short: $(M_MODULE_SHORT)
  kind: infrastructure
  provider: azure
endef

define M_CONFIG_CONTENT
kind: $(M_MODULE_SHORT)-config
$(M_MODULE_SHORT):
  name: $(M_NAME)
  rg_name: $(M_RG_NAME)
  vnet_name: $(M_VNET_NAME)
  address_prefix: $(M_ADDRESS_PREFIX)
  subnet_name: $(M_SUBNET_NAME)
  existing_subnet: $(M_EXISTING_SUBNET)

  # default node pool
  size: $(M_SIZE)
  min: $(M_MIN)
  max: $(M_MAX)
  vm_size: $(M_VM_SIZE)
  disk_size: $(M_DISK_SIZE)
  auto_scaling: $(M_AUTO_SCALING)
  rsa_pub_path: "$(M_SHARED)/$(M_VMS_RSA).pub"
endef

define M_STATE_INITIAL
kind: state
$(M_MODULE_SHORT):
  status: initialized
endef
