output "kubeconfig" {
  value     = module.aks.kubeconfig
  sensitive = true
}
