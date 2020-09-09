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
  resource_group: $(M_RESOURCE_GROUP)
  vnet: $(M_VNET)
  address_prefix: $(M_ADDRESS_PREFIX)
  rsa_pub_path: "$(M_SHARED)/$(M_VMS_RSA).pub"

  # default node pool
  size: $(M_SIZE)
  min: $(M_MIN)
  max: $(M_MAX)
  vm_size: $(M_VM_SIZE)
  disk_size: $(M_DISK_SIZE)
  auto_scaling: $(M_AUTO_SCALING)
endef

define M_STATE_INITIAL
kind: state
$(M_MODULE_SHORT):
  status: initialized
endef
