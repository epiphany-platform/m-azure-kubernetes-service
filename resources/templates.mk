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
  kubernetes_version: $(M_K8S_VERSION)
  enable_node_public_ip: $(M_PUBLIC_IP_ENABLED)
  enable_rbac: $(M_RBAC_ENABLED)
  azure_ad: $(M_AZURE_AD)
  default_node_pool: $(M_DEFAULT_NODE_POOL)
  auto_scaler_profile: $(M_AUTO_SCALER_PROFILE)
  rsa_pub_path: "$(M_SHARED)/$(M_VMS_RSA).pub"
  identity_type: $(M_IDENTITY_TYPE)
  kube_dashboard_enabled: $(M_K8S_DASHBOARD_ENABLED)
  admin_username: $(M_ADMIN_USERNAME)
endef

define M_STATE_INITIAL
kind: state
$(M_MODULE_SHORT):
  status: initialized
endef
