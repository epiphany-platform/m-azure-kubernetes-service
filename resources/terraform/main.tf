data "azurerm_resource_group" "main_rg" {
  name = var.resource_group
}

data "azurerm_virtual_network" "vnet" {
  name                = var.vnet
  resource_group_name = data.azurerm_resource_group.main_rg.name
}

resource "azurerm_subnet" "subnet" {
  name                 = "${var.name}-snet-aks"
  resource_group_name  = data.azurerm_resource_group.main_rg.name
  virtual_network_name = data.azurerm_virtual_network.vnet.name
  address_prefixes     = [ cidrsubnet(var.address_prefix, 8, 16) ]
}

module "aks" {
  source         = "./modules/aks"
  name           = var.name
  resource_group = data.azurerm_resource_group.main_rg.name
  subnet_id      = azurerm_subnet.subnet.id
  size           = var.size
  min            = var.min
  max            = var.max
  vm_size        = var.vm_size
  disk_size      = var.disk_size
  auto_scaling   = var.auto_scaling
  tf_key_path    = var.rsa_pub_path
}
