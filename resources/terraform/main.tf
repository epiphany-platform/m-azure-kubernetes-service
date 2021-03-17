data "azurerm_resource_group" "main_rg" {
  name = var.rg_name
}

data "azurerm_virtual_network" "vnet" {
  name                = var.vnet_name
  resource_group_name = data.azurerm_resource_group.main_rg.name
}

data "azurerm_subnet" "subnet" {
  name                 = var.subnet_name
  virtual_network_name = data.azurerm_virtual_network.vnet.name
  resource_group_name  = data.azurerm_resource_group.main_rg.name
}

resource "azurerm_kubernetes_cluster" "aks" {
  name                = "${var.name}-aks"
  resource_group_name = data.azurerm_resource_group.main_rg.name
  location            = data.azurerm_resource_group.main_rg.location
  dns_prefix          = var.name
  node_resource_group = "${var.name}-rg-worker"
  kubernetes_version  = var.kubernetes_version

  default_node_pool {
    name                  = "default"
    node_count            = var.default_node_pool.size
    vm_size               = var.default_node_pool.vm_size
    vnet_subnet_id        = data.azurerm_subnet.subnet.id
    orchestrator_version  = var.kubernetes_version
    os_disk_size_gb       = var.default_node_pool.disk_gb_size
    enable_node_public_ip = var.enable_node_public_ip
    type                  = var.default_node_pool.type
    enable_auto_scaling   = var.default_node_pool.auto_scaling
    min_count             = var.default_node_pool.min
    max_count             = var.default_node_pool.max
  }

  identity {
    type = var.identity_type
  }

  linux_profile {
    admin_username = var.admin_username
    ssh_key {
      key_data = file(var.rsa_pub_path)
    }
  }

  network_profile {
    network_plugin     = var.network_plugin
    network_policy     = var.network_policy
    service_cidr       = "10.96.0.0/16"
    dns_service_ip     = "10.96.0.10"
    docker_bridge_cidr = "172.17.0.1/16"
  }

  role_based_access_control {
    enabled = var.enable_rbac
    dynamic "azure_active_directory" {
      for_each = var.azure_ad == null ? [] : [var.azure_ad]
      content {
        managed                = azure_active_directory.managed
        tenant_id              = azure_active_directory.tenant_id
        admin_group_object_ids = azure_active_directory.admin_group_object_ids
      }
    }
  }

  auto_scaler_profile {
    balance_similar_node_groups      = var.auto_scaler_profile.balance_similar_node_groups
    max_graceful_termination_sec     = var.auto_scaler_profile.max_graceful_termination_sec
    scale_down_delay_after_add       = var.auto_scaler_profile.scale_down_delay_after_add
    scale_down_delay_after_delete    = var.auto_scaler_profile.scale_down_delay_after_delete
    scale_down_delay_after_failure   = var.auto_scaler_profile.scale_down_delay_after_failure
    scan_interval                    = var.auto_scaler_profile.scan_interval
    scale_down_unneeded              = var.auto_scaler_profile.scale_down_unneeded
    scale_down_unready               = var.auto_scaler_profile.scale_down_unready
    scale_down_utilization_threshold = var.auto_scaler_profile.scale_down_utilization_threshold
  }

  tags = {
    Environment = var.name
  }
}
