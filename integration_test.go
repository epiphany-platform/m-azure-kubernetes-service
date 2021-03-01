package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	azks "github.com/epiphany-platform/e-structures/azks/v0"
	"github.com/epiphany-platform/e-structures/utils/to"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/epiphany-platform/m-azure-kubernetes-service/cmd"
	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/docker"
	"golang.org/x/crypto/ssh"
)

func TestMetadata(t *testing.T) {
	tests := []struct {
		name               string
		wantOutputTemplate string
		parameters         map[string]string
	}{
		{
			name: "default metadata",
			wantOutputTemplate: `labels:
  kind: infrastructure
  name: Azure Kubernetes Service
  provider: azure
  provides-pubips: true
  short: azks
  version: %s
`,
			parameters: nil,
		},
		{
			name:               "json metadata",
			wantOutputTemplate: "{\"labels\":{\"kind\":\"infrastructure\",\"name\":\"Azure Kubernetes Service\",\"provider\":\"azure\",\"provides-pubips\":true,\"short\":\"azks\",\"version\":\"%s\"}}",
			parameters:         map[string]string{"--json": "true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := dockerRun(t, "metadata", tt.parameters, nil, "")
			if diff := deep.Equal(gotOutput, fmt.Sprintf(tt.wantOutputTemplate, cmd.Version)); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name                 string
		existingStateContent []byte
		initParams           map[string]string
		wantConfigContent    []byte
	}{
		{
			name:                 "no state no init params",
			existingStateContent: nil,
			initParams:           nil,
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "epiphany",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "epiphany-rg",
    "vnet_name": "epiphany-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "empty state no init params",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "initialized",
    "config": null,
    "output": null
  }
}`),
			initParams: nil,
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "epiphany",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "epiphany-rg",
    "vnet_name": "epiphany-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "existing state no free subnets no init params",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "applied",
    "config": {
      "kind": "azbi",
      "version": "v0.1.1",
      "params": {
        "name": "some-test-name",
        "location": "northeurope",
        "address_space": [
          "10.0.0.0/16"
        ],
        "subnets": [
          {
            "name": "main",
            "address_prefixes": [
              "10.0.1.0/24"
            ]
          }
        ],
        "vm_groups": [
          {
            "name": "vm-group0",
            "vm_count": 1,
            "vm_size": "Standard_DS2_v2",
            "use_public_ip": true,
            "subnet_names": [
              "main"
            ],
            "vm_image": {
              "publisher": "Canonical",
              "offer": "UbuntuServer",
              "sku": "18.04-LTS",
              "version": "18.04.202006101"
            },
            "data_disks": [
              {
                "disk_size_gb": 10
              }
            ]
          }
        ],
        "rsa_pub_path": "/shared/vms_rsa.pub"
      }
    },
    "output": {
      "rg_name": "some-test-name-rg",
      "vnet_name": "some-test-name-vnet",
      "vm_groups": [
        {
          "vm_group_name": "vm-group0",
          "vms": [
            {
              "vm_name": "some-test-name-vm-group0-0",
              "private_ips": [
                "10.0.1.4"
              ],
              "public_ip": "12.34.56.78",
              "data_disks": [
                {
                  "size": 10,
                  "lun": 10
                }
              ]
            }
          ]
        }
      ]
    }
  }
}
`),
			initParams: nil,
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "some-test-name",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "some-test-name-rg",
    "vnet_name": "some-test-name-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "existing state with free subnet no init params",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "applied",
    "config": {
      "kind": "azbi",
      "version": "v0.1.1",
      "params": {
        "name": "some-test-name",
        "location": "northeurope",
        "address_space": [
          "10.0.0.0/16"
        ],
        "subnets": [
          {
            "name": "main",
            "address_prefixes": [
              "10.0.1.0/24"
            ]
          },
          {
            "name": "kubernetes",
            "address_prefixes": [
              "10.0.2.0/24"
            ]
          }
        ],
        "vm_groups": [
          {
            "name": "vm-group0",
            "vm_count": 1,
            "vm_size": "Standard_DS2_v2",
            "use_public_ip": false,
            "subnet_names": [
              "main"
            ],
            "vm_image": {
              "publisher": "Canonical",
              "offer": "UbuntuServer",
              "sku": "18.04-LTS",
              "version": "18.04.202006101"
            },
            "data_disks": [
              {
                "disk_size_gb": 10
              }
            ]
          }
        ],
        "rsa_pub_path": "/shared/vms_rsa.pub"
      }
    },
    "output": {
      "rg_name": "some-test-name-rg",
      "vnet_name": "some-test-name-vnet",
      "vm_groups": [
        {
          "vm_group_name": "vm-group0",
          "vms": [
            {
              "vm_name": "some-test-name-vm-group0-0",
              "private_ips": [
                "10.0.1.4"
              ],
              "public_ip": "",
              "data_disks": [
                {
                  "size": 10,
                  "lun": 10
                }
              ]
            }
          ]
        }
      ]
    }
  }
}`),
			initParams: nil,
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "some-test-name",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "some-test-name-rg",
    "vnet_name": "some-test-name-vnet",
    "subnet_name": "kubernetes",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "existing state with omit_state init param",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "applied",
    "config": {
      "kind": "azbi",
      "version": "v0.1.1",
      "params": {
        "name": "some-test-name",
        "location": "northeurope",
        "address_space": [
          "10.0.0.0/16"
        ],
        "subnets": [
          {
            "name": "main",
            "address_prefixes": [
              "10.0.1.0/24"
            ]
          },
          {
            "name": "kubernetes",
            "address_prefixes": [
              "10.0.2.0/24"
            ]
          }
        ],
        "vm_groups": [
          {
            "name": "vm-group0",
            "vm_count": 1,
            "vm_size": "Standard_DS2_v2",
            "use_public_ip": false,
            "subnet_names": [
              "main"
            ],
            "vm_image": {
              "publisher": "Canonical",
              "offer": "UbuntuServer",
              "sku": "18.04-LTS",
              "version": "18.04.202006101"
            },
            "data_disks": [
              {
                "disk_size_gb": 10
              }
            ]
          }
        ],
        "rsa_pub_path": "/shared/vms_rsa.pub"
      }
    },
    "output": {
      "rg_name": "some-test-name-rg",
      "vnet_name": "some-test-name-vnet",
      "vm_groups": [
        {
          "vm_group_name": "vm-group0",
          "vms": [
            {
              "vm_name": "some-test-name-vm-group0-0",
              "private_ips": [
                "10.0.1.4"
              ],
              "public_ip": "",
              "data_disks": [
                {
                  "size": 10,
                  "lun": 10
                }
              ]
            }
          ]
        }
      ]
    }
  }
}`),
			initParams: map[string]string{"--omit_state": "true"},
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "epiphany",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "epiphany-rg",
    "vnet_name": "epiphany-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "empty state with init params 1",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "initialized",
    "config": null,
    "output": null
  }
}`),
			initParams: map[string]string{"--name": "name-test-1"},
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "name-test-1",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "epiphany-rg",
    "vnet_name": "epiphany-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "empty state with init params 2",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "initialized",
    "config": null,
    "output": null
  }
}`),
			initParams: map[string]string{"--kubernetes_version": "1.2.3"},
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "epiphany",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "epiphany-rg",
    "vnet_name": "epiphany-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.2.3",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "empty state with init params 3",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "initialized",
    "config": null,
    "output": null
  }
}`),
			initParams: map[string]string{"--rg_name": "rg-name-test-1"},
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "epiphany",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "rg-name-test-1",
    "vnet_name": "epiphany-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "empty state with init params 4",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "initialized",
    "config": null,
    "output": null
  }
}`),
			initParams: map[string]string{"--subnet_name": "subnet-name-test-1"},
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "epiphany",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "epiphany-rg",
    "vnet_name": "epiphany-vnet",
    "subnet_name": "subnet-name-test-1",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "empty state with init params 5",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "initialized",
    "config": null,
    "output": null
  }
}`),
			initParams: map[string]string{"--vnet_name": "vnet-name-test-1"},
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "epiphany",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "epiphany-rg",
    "vnet_name": "vnet-name-test-1",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
		{
			name: "empty state with init params 6",
			existingStateContent: []byte(`{
  "kind": "state",
  "version": "v0.0.2",
  "azbi": {
    "status": "initialized",
    "config": null,
    "output": null
  }
}`),
			initParams: map[string]string{"--vms_rsa": "test_rsa"},
			wantConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "epiphany",
    "location": "northeurope",
    "rsa_pub_path": "/shared/test_rsa.pub",
    "rg_name": "epiphany-rg",
    "vnet_name": "epiphany-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, remoteSharedPath, localSharedPath, environments, _ := setup(t, tt.initParams)
			defer cleanup(t, localSharedPath, environments["SUBSCRIPTION_ID"], name)

			if tt.existingStateContent != nil {
				err := ioutil.WriteFile(path.Join(localSharedPath, "state.json"), tt.existingStateContent, 0600)
				if err != nil {
					t.Fatal(err)
				}
			}

			_ = dockerRun(t, "init", tt.initParams, nil, remoteSharedPath)
			gotFileContent, err := ioutil.ReadFile(path.Join(localSharedPath, "azks/azks-config.json"))
			if err != nil {
				t.Errorf("wasnt able to read form output file: %v", err)
			}

			e, err := jsonEqual(t, mustReMarshalJson(t, gotFileContent), mustReMarshalJson(t, tt.wantConfigContent))
			if err != nil {
				t.Error(err)
			}
			if !e {
				t.Errorf("Init result got = \n%v\nwant = \n%v\n", string(gotFileContent), string(tt.wantConfigContent))
			}
		})
	}
}

// In TestPlan it's extremely important to keep location, rg, vnet and subnet names stable, as they are created once and removed once to optimise test time.
func TestPlan(t *testing.T) {
	tests := []struct {
		name                   string
		existingConfigContent  []byte
		planParams             map[string]string
		wantPlanOutputLastLine string
	}{
		{
			name: "initialized without any prior steps",
			existingConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "azks-integration-test",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "azks-integration-test-rg",
    "vnet_name": "azks-integration-test-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
			planParams:             map[string]string{"--debug": "true"},
			wantPlanOutputLastLine: "\tAdd: 1, Change: 0, Destroy: 0",
		},
	}

	config := &azks.Config{}
	err := json.Unmarshal(tests[0].existingConfigContent, config)
	if err != nil {
		t.Fatal(err)
	}
	fakeInitParams := map[string]string{"--name": *config.Params.Name}
	name, remoteSharedPath, localSharedPath, environments, _ := setup(t, fakeInitParams)
	defer cleanup(t, localSharedPath, environments["SUBSCRIPTION_ID"], name)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.existingConfigContent != nil {
				err := os.MkdirAll(path.Join(localSharedPath, "azks"), os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
				err = ioutil.WriteFile(path.Join(localSharedPath, "azks/azks-config.json"), tt.existingConfigContent, 0600)
				if err != nil {
					t.Fatal(err)
				}
			}
			err := ioutil.WriteFile(path.Join(localSharedPath, "state.json"), []byte(`{
	"kind": "state",
	"version": "v0.0.3",
	"azbi": {
		"status": "",
		"config": null,
		"output": null
	},
	"azks": {
		"status": "initialized",
		"config": null,
		"output": null
	}
}`), 0600)
			if err != nil {
				t.Fatal(err)
			}

			prepareAzurePreRequisites(t, environments["SUBSCRIPTION_ID"], config)

			gotPlanOutputLastLine := getLastLineFromMultilineSting(t, dockerRun(t, "plan", tt.planParams, environments, remoteSharedPath))
			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}
		})
	}
}

// CAUTION: every entry in tests table will make this test longer for about 12 - 13 minutes!
func TestApply(t *testing.T) {
	tests := []struct {
		name                    string
		existingConfigContent   []byte
		planParams              map[string]string
		wantPlanOutputLastLine  string
		applyParams             map[string]string
		wantApplyOutputLastLine string
	}{
		{
			name: "initialized without any prior steps",
			existingConfigContent: []byte(`{
  "kind": "azks",
  "version": "v0.0.1",
  "params": {
    "name": "azks-integration-test",
    "location": "northeurope",
    "rsa_pub_path": "/shared/vms_rsa.pub",
    "rg_name": "azks-integration-test-rg",
    "vnet_name": "azks-integration-test-vnet",
    "subnet_name": "azks",
    "kubernetes_version": "1.18.14",
    "enable_node_public_ip": false,
    "enable_rbac": false,
    "default_node_pool": {
      "size": 2,
      "min": 2,
      "max": 5,
      "vm_size": "Standard_DS2_v2",
      "disk_size": "36",
      "auto_scaling": true,
      "type": "VirtualMachineScaleSets"
    },
    "auto_scaler_profile": {
      "balance_similar_node_groups": false,
      "max_graceful_termination_sec": "600",
      "scale_down_delay_after_add": "10m",
      "scale_down_delay_after_delete": "10s",
      "scale_down_delay_after_failure": "10m",
      "scan_interval": "10s",
      "scale_down_unneeded": "10m",
      "scale_down_unready": "10m",
      "scale_down_utilization_threshold": "0.5"
    },
    "azure_ad": null,
    "identity_type": "SystemAssigned",
    "admin_username": "operations"
  }
}`),
			planParams:              map[string]string{"--debug": "true"},
			wantPlanOutputLastLine:  "\tAdd: 1, Change: 0, Destroy: 0",
			applyParams:             map[string]string{"--debug": "true"},
			wantApplyOutputLastLine: "\tAdd: 1, Change: 0, Destroy: 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config := &azks.Config{}
			err := json.Unmarshal(tt.existingConfigContent, config)
			if err != nil {
				t.Fatal(err)
			}
			fakeInitParams := map[string]string{"--name": *config.Params.Name}
			name, remoteSharedPath, localSharedPath, environments, _ := setup(t, fakeInitParams)
			defer cleanup(t, localSharedPath, environments["SUBSCRIPTION_ID"], name)

			if tt.existingConfigContent != nil {
				err := os.MkdirAll(path.Join(localSharedPath, "azks"), os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
				err = ioutil.WriteFile(path.Join(localSharedPath, "azks/azks-config.json"), tt.existingConfigContent, 0600)
				if err != nil {
					t.Fatal(err)
				}
			}
			err = ioutil.WriteFile(path.Join(localSharedPath, "state.json"), []byte(`{
	"kind": "state",
	"version": "v0.0.3",
	"azbi": {
		"status": "",
		"config": null,
		"output": null
	},
	"azks": {
		"status": "initialized",
		"config": null,
		"output": null
	}
}`), 0600)
			if err != nil {
				t.Fatal(err)
			}

			prepareAzurePreRequisites(t, environments["SUBSCRIPTION_ID"], config)

			gotPlanOutputLastLine := getLastLineFromMultilineSting(t, dockerRun(t, "plan", tt.planParams, environments, remoteSharedPath))
			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			gotApplyOutputLastLine := getLastLineFromMultilineSting(t, dockerRun(t, "apply", tt.applyParams, environments, remoteSharedPath))

			if diff := deep.Equal(gotApplyOutputLastLine, tt.wantApplyOutputLastLine); diff != nil {
				t.Error(diff)
			}
		})
	}
}

// dockerRun function wraps docker run operation and returns `docker run` output.
func dockerRun(t *testing.T, command string, parameters map[string]string, environments map[string]string, sharedPath string) string {
	commandWithParameters := []string{command}
	for k, v := range parameters {
		commandWithParameters = append(commandWithParameters, fmt.Sprintf("%s=%s", k, v))
	}

	var opts *docker.RunOptions
	if sharedPath != "" {
		opts = &docker.RunOptions{
			Command: commandWithParameters,
			Remove:  true,
			Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
		}
	} else {
		opts = &docker.RunOptions{
			Command: commandWithParameters,
			Remove:  true,
		}
	}
	var envs []string
	for k, v := range environments {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	opts.EnvironmentVariables = envs

	//in case of error Run function calls FailNow anyways
	return docker.Run(t, fmt.Sprintf("%s:%s", prepareImageTag(t), cmd.Version), opts)
}

// prepareImageTag returns IMAGE_REPOSITORY environment variable
func prepareImageTag(t *testing.T) string {
	imageRepository := os.Getenv("IMAGE_REPOSITORY")
	if len(imageRepository) == 0 {
		t.Fatal("expected IMAGE_REPOSITORY environment variable")
	}
	return imageRepository
}

// setup function ensures that all prerequisites for tests are in place.
func setup(t *testing.T, initParams map[string]string) (string, string, string, map[string]string, ssh.Signer) {
	rsaName := "vms_rsa"
	if value, ok := initParams["--vms_rsa"]; ok {
		rsaName = value
	}
	name := "epiphany-rg"
	if value, ok := initParams["--name"]; ok {
		name = value
	}

	environments := loadEnvironmentVariables(t)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	var remoteSharedPath string
	if v, ok := environments["K8S_VOL_PATH"]; ok && v != "" {
		remoteSharedPath = v
	} else {
		remoteSharedPath = path.Join(wd, "shared")
	}
	var localSharedPath string
	if v, ok := environments["K8S_HOST_PATH"]; ok && v != "" {
		localSharedPath = v
	} else {
		localSharedPath = path.Join(wd, "shared")
	}
	err = os.MkdirAll(localSharedPath, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	privateKey := generateRsaKeyPair(t, localSharedPath, rsaName)
	if isResourceGroupPresent(t, environments["SUBSCRIPTION_ID"], name) {
		removeResourceGroup(t, environments["SUBSCRIPTION_ID"], name)
	}
	return name, remoteSharedPath, localSharedPath, environments, privateKey
}

// prepareAzurePreRequisites creates RG, VNET and Subnet for plan and apply tests.
func prepareAzurePreRequisites(t *testing.T, subscriptionId string, config *azks.Config) {
	t.Log("prepareAzurePreRequisites()")
	groupsClient := resources.NewGroupsClient(subscriptionId)
	vnetClient := network.NewVirtualNetworksClient(subscriptionId)

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		t.Error(err)
	}

	groupsClient.Authorizer = authorizer
	vnetClient.Authorizer = authorizer

	ctx := context.TODO()
	dts := time.Now().Format("2006-01-02 15:04:05")

	g, err := groupsClient.CreateOrUpdate(ctx, *config.Params.RgName,
		resources.Group{
			Location: config.Params.Location,
			Tags: map[string]*string{
				"product":    to.StrPtr("epiphany"),
				"module":     to.StrPtr("azks"),
				"purpose":    to.StrPtr("integration-tests"),
				"created_at": to.StrPtr(dts),
				"creator":    to.StrPtr("automation"),
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = vnetClient.CreateOrUpdate(ctx, *g.Name, *config.Params.VnetName,
		network.VirtualNetwork{
			Location: g.Location,
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{"10.0.0.0/16"},
				},
				Subnets: &[]network.Subnet{
					{
						Name: config.Params.SubnetName,
						SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
							AddressPrefix: to.StrPtr("10.0.1.0/24"),
						},
					},
				},
			},
			Tags: map[string]*string{
				"product":    to.StrPtr("epiphany"),
				"module":     to.StrPtr("azks"),
				"purpose":    to.StrPtr("integration-tests"),
				"created_at": to.StrPtr(dts),
				"creator":    to.StrPtr("automation"),
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

// cleanup function removes directories created during test and ensures that resource
// group gets removed if it was created.
func cleanup(t *testing.T, sharedPath string, subscriptionId string, name string) {
	t.Log("cleanup()")
	_ = os.RemoveAll(sharedPath)
	if isResourceGroupPresent(t, subscriptionId, name) {
		removeResourceGroup(t, subscriptionId, name)
	}
}

// isResourceGroupPresent function checks if resource group with given name exists.
func isResourceGroupPresent(t *testing.T, subscriptionId string, name string) bool {
	t.Log("isResourceGroupPresent()")
	groupsClient := resources.NewGroupsClient(subscriptionId)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		t.Error(err)
	}
	groupsClient.Authorizer = authorizer
	rgName := fmt.Sprintf("%s-rg", name)
	_, err = groupsClient.Get(context.TODO(), rgName)
	if err != nil {
		return false
	} else {
		return true
	}
}

// removeResourceGroup function invokes Delete operation on provided resource
// group name and waits for operation completion.
func removeResourceGroup(t *testing.T, subscriptionId string, name string) {
	t.Log("Will prepare new az groups client")
	ctx := context.TODO()
	groupsClient := resources.NewGroupsClient(subscriptionId)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	groupsClient.Authorizer = authorizer
	rgName := fmt.Sprintf("%s-rg", name)
	t.Log("Will perform delete RG operation")
	gdf, err := groupsClient.Delete(ctx, rgName)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	now := time.Now()

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		defer close(done)
		t.Log("Will start waiting for RG deletion finish.")
		err = gdf.WaitForCompletionRef(ctx, groupsClient.BaseClient.Client)
		t.Log("Finished RG deletion.")
		if err != nil {
			t.Fatal(err)
		}
	}()

	for {
		select {
		case <-ticker.C:
			t.Logf("Waiting for deletion to complete: %v", time.Since(now).Round(time.Second))
		case <-done:
			t.Logf("Finished waiting for RG deletion.")
			ticker.Stop()
			return
		}
	}
}

// generateRsaKeyPair function generates RSA public and private keys and returns
// ssh.Signer that can create signatures that verify against a public key.
func generateRsaKeyPair(t *testing.T, directory string, name string) ssh.Signer {
	privateRsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err)
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRsaKey)}
	privateKeyBytes := pem.EncodeToMemory(pemBlock)

	publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	err = ioutil.WriteFile(path.Join(directory, name), privateKeyBytes, 0600)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join(directory, fmt.Sprintf("%s.pub", name)), publicKeyBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		t.Fatal(err)
	}
	return signer
}

// getLastLineFromMultilineSting is helper function to extract just last line
// from multiline string.
func getLastLineFromMultilineSting(t *testing.T, s string) string {
	in := strings.NewReader(s)
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if err == io.EOF {
			return string(line)
		}
	}
}

// loadEnvironmentVariables obtains 6 variables from environment.
// Two of them (K8S_VOL_PATH and K8S_HOST_PATH) are optional and
// are not checked but another four (AZURE_CLIENT_ID, AZURE_CLIENT_SECRET
// AZURE_SUBSCRIPTION_ID and AZURE_TENANT_ID) are required and if missing
// will cause test to fail.
func loadEnvironmentVariables(t *testing.T) map[string]string {
	result := make(map[string]string)
	result["CLIENT_ID"] = os.Getenv("AZURE_CLIENT_ID")
	if len(result["CLIENT_ID"]) == 0 {
		t.Fatalf("expected AZURE_CLIENT_ID environment variable")
	}
	result["CLIENT_SECRET"] = os.Getenv("AZURE_CLIENT_SECRET")
	if len(result["CLIENT_SECRET"]) == 0 {
		t.Fatalf("expected AZURE_CLIENT_SECRET environment variable")
	}
	result["SUBSCRIPTION_ID"] = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(result["SUBSCRIPTION_ID"]) == 0 {
		t.Fatalf("expected AZURE_SUBSCRIPTION_ID environment variable")
	}
	result["TENANT_ID"] = os.Getenv("AZURE_TENANT_ID")
	if len(result["TENANT_ID"]) == 0 {
		t.Fatalf("expected AZURE_TENANT_ID environment variable")
	}
	result["K8S_VOL_PATH"] = os.Getenv("K8S_VOL_PATH")
	result["K8S_HOST_PATH"] = os.Getenv("K8S_HOST_PATH")
	return result
}

// jsonEqual decodes two JSONs into interface{}s and then compares them
func jsonEqual(t *testing.T, got, want io.Reader) (bool, error) {
	var dGot, dWant interface{}
	d := json.NewDecoder(got)
	if err := d.Decode(&dGot); err != nil {
		t.Errorf("wasn't able to Decode: %v", err)
	}
	d = json.NewDecoder(want)
	if err := d.Decode(&dWant); err != nil {
		t.Errorf("wasn't able to Decode: %v", err)
	}
	t.Logf("compared got: \n%v\nand want: \n%v\n", dGot, dWant)
	return reflect.DeepEqual(dGot, dWant), nil
}

// mustReMarshalJson Unmashals and then Marshals JSON to make sure that it is ordered in predictable way
func mustReMarshalJson(t *testing.T, b []byte) io.Reader {
	var j interface{}
	err := json.Unmarshal(b, &j)
	if err != nil {
		t.Errorf("wasn't able to Unmarshal: %v\n\njson:\n%s\n", err, string(b))
	}
	output, err := json.Marshal(j)
	if err != nil {
		t.Errorf("wasn't able to Marshal: %v", err)
	}
	return bytes.NewReader(output)
}
