variable "existing_subnet_id" {
	description = "ID of the existing subnet to deploy AKS cluster in"
	type        = string
}

variable "existing_rg_name" {
	description = "Name of the existing resource group to deploy AKS cluster in"
	type        = string
}
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
  type = string
}
