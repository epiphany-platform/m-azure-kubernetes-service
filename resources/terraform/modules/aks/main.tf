data "azurerm_resource_group" "main_rg" {
  name = var.resource_group
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
    node_count            = var.size
    vm_size               = var.vm_size
    vnet_subnet_id        = var.subnet_id
    orchestrator_version  = var.kubernetes_version
    os_disk_size_gb       = var.disk_size
    enable_node_public_ip = var.enable_node_public_ip
    type                  = var.default_node_pool_type
    enable_auto_scaling   = var.auto_scaling
    min_count             = var.min
    max_count             = var.max
  }

  identity {
    type = var.identity_type
  }

  linux_profile {
    admin_username = var.admin_username
    ssh_key {
      key_data = file(var.tf_key_path)
    }
  }

  network_profile {
    network_plugin     = var.network_plugin
    network_policy     = var.network_policy
    service_cidr       = "10.96.0.0/16"
    dns_service_ip     = "10.96.0.10"
    docker_bridge_cidr = "172.17.0.1/16"
  }

  addon_profile {
    kube_dashboard {
      enabled = var.kube_dashboard_enabled
    }
  }

  auto_scaler_profile {
    balance_similar_node_groups      = var.balance_similar_node_groups
    max_graceful_termination_sec     = var.max_graceful_termination_sec
    scale_down_delay_after_add       = var.scale_down_delay_after_add
    scale_down_delay_after_delete    = var.scale_down_delay_after_delete
    scale_down_delay_after_failure   = var.scale_down_delay_after_failure
    scan_interval                    = var.scan_interval
    scale_down_unneeded              = var.scale_down_unneeded
    scale_down_unready               = var.scale_down_unready
    scale_down_utilization_threshold = var.scale_down_utilization_threshold
  }

  tags = {
    Environment = var.name
  }
}
