variable "name" {
  description = "Prefix for resource names"
  type        = string
}

variable "rg_name" {
  description = "Name of an resource group to deploy AKS cluster in"
  type        = string
  default     = ""
}

variable "vnet_name" {
  description = "Main Virtual Network's name"
  type        = string
}

variable "address_prefix" {
  description = "SubNetwork address space for AKS"
  type        = string
}

# DEFAULT NODE-POOL

variable "size" {
  type = number
}

variable "min" {
  type = number
}

variable "max" {
  type = number
}

variable "vm_size" {
  type = string
}

variable "disk_size" {
  type = string
}

variable "auto_scaling" {
  type = bool
}

variable "rsa_pub_path" {
  description = "Public ssh key path"
  type        = string
  default     = "/shared/vms_rsa.pub"
}

# AUTO-SCALER PROFILE

variable "balance_similar_node_groups" {
  type    = bool
  default = false
}

variable "max_graceful_termination_sec" {
  type    = string
  default = "600"
}

variable "scale_down_delay_after_add" {
  type    = string
  default = "10m"
}

variable "scale_down_delay_after_delete" {
  type    = string
  default = "10s"
}

variable "scale_down_delay_after_failure" {
  type    = string
  default = "10m"
}

variable "scan_interval" {
  type    = string
  default = "10s"
}

variable "scale_down_unneeded" {
  type    = string
  default = "10m"
}

variable "scale_down_unready" {
  type    = string
  default = "10m"
}

variable "scale_down_utilization_threshold" {
  type    = string
  default = "0.5"
}

# RANDOM DEFAULTS

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.18.6"
}

variable "enable_node_public_ip" {
  type    = bool
  default = false
}

variable "default_node_pool_type" {
  type    = string
  default = "VirtualMachineScaleSets"
}

variable "identity_type" {
  type    = string
  default = "SystemAssigned"
}

variable "network_plugin" {
  type    = string
  default = "azure"
}

variable "network_policy" {
  type    = string
  default = "azure"
}

variable "kube_dashboard_enabled" {
  type    = bool
  default = true
}

variable "admin_username" {
  description = "Admin user on Linux OS"
  type        = string
  default     = "operations"
}
