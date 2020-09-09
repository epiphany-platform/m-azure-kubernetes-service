output "kubeconfig" {
  value = azurerm_kubernetes_cluster.aks.kube_config_raw
}
