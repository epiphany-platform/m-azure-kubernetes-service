variable "name" {
  description = "Prefix for resource names"
  type        = string
}

variable "resource_group" {
  description = "Name of an resource group to deploy AKS cluster in"
  type        = string
  default     = ""
}

variable "vnet" {
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
  type = string
}
