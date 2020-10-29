data "azurerm_resource_group" "main_rg" {
  name = var.existing_rg_name != "unset" ? var.existing_rg_name : var.rg_name
}

data "azurerm_virtual_network" "vnet" {
  count = var.existing_subnet_id != "unset" ? 0 : 1
  name                = var.vnet_name
  resource_group_name = data.azurerm_resource_group.main_rg.name
}

resource "azurerm_subnet" "subnet" {
  count = var.existing_subnet_id != "unset" ? 0 : 1
  name                 = "${var.name}-snet-aks"
  resource_group_name  = data.azurerm_resource_group.main_rg.name
  virtual_network_name = data.azurerm_virtual_network.vnet[0].name
  address_prefixes     = [ cidrsubnet(var.address_prefix, 8, 16) ]
}

module "aks" {
  source       = "./modules/aks"
  name         = var.name
  rg_name      = data.azurerm_resource_group.main_rg.name
  subnet_id    = var.existing_subnet_id != "unset" ? var.existing_subnet_id : azurerm_subnet.subnet[0].id
  size         = var.size
  min          = var.min
  max          = var.max
  vm_size      = var.vm_size
  disk_size    = var.disk_size
  auto_scaling = var.auto_scaling
  tf_key_path  = var.rsa_pub_path
}
