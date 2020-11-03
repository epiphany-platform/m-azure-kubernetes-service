#!/usr/bin/env bash

function usage() {
  echo "usage:
    $0 cleanup
    $0 setup
    $0 generate_junit_report
    $0 test-default-config-suite [image_name]
    $0 test-config-with-variables-suite [image_name]
    $0 test-plan-suite [image_name] [ARM_CLIENT_ID] [ARM_CLIENT_SECRET] [ARM_SUBSCRIPTION_ID] [ARM_TENANT_ID]
    $0 test-apply-suite [image_name] [ARM_CLIENT_ID] [ARM_CLIENT_SECRET] [ARM_SUBSCRIPTION_ID] [ARM_TENANT_ID]
    "
}

function test-default-config-suite() {
  #$1 is IMAGE_NAME
  start_suite test-default-config

  r=0
  run_test init-default-config "$r" "$1"
  r=$?
  run_test check-default-config-content "$r" "$1"
  r=$?

  stop_suite test-default-config "$r"
}

function test-config-with-variables-suite() {
  #$1 is IMAGE_NAME
  start_suite test-config-with-variables

  r=0
  run_test init-2-5-autoscaled-aks "$r" "$1"
  r=$?
  run_test check-2-5-autoscaled-aks-config-content "$r" "$1"
  r=$?

  stop_suite test-config-with-variables "$r"
}

function test-plan-suite() {
  #$1 is IMAGE_NAME
  #$2 is ARM_CLIENT_ID
  #$3 is ARM_CLIENT_SECRET
  #$4 is ARM_SUBSCRIPTION_ID
  #$5 is ARM_TENANT_ID
  start_suite test-plan
  r=0
  run_test init-2-5-autoscaled-aks "$r" "$1"
  r=$?
  run_test check-2-5-autoscaled-aks-config-content "$r" "$1"
  r=$?
  run_test prepare-azks-module-tests-rg "$r" "$1 $2 $3 $4 $5"
  r=$?
  run_test plan-2-5-autoscaled-aks "$r" "$1 $2 $3 $4 $5"
  r=$?
  run_test check-2-5-autoscaled-aks-plan "$r" "$1"
  r=0
  run_test teardown-azks-module-tests-rg "$r" "$1 $2 $3 $4 $5"
  r=$?

  stop_suite test-plan "$r"
}

function test-apply-suite() {
  #$1 is IMAGE_NAME
  #$2 is ARM_CLIENT_ID
  #$3 is ARM_CLIENT_SECRET
  #$4 is ARM_SUBSCRIPTION_ID
  #$5 is ARM_TENANT_ID
  start_suite test-plan

  r=0
  run_test init-2-5-autoscaled-aks "$r" "$1"
  r=$?
  run_test check-2-5-autoscaled-aks-config-content "$r" "$1"
  r=$?
  run_test prepare-azks-module-tests-rg "$r" "$1 $2 $3 $4 $5"
  r=$?
  run_test plan-2-5-autoscaled-aks "$r" "$1 $2 $3 $4 $5"
  r=$?
  run_test check-2-5-autoscaled-aks-plan "$r" "$1"
  r=$?
  run_test apply-2-5-autoscaled-aks "$r" "$1 $2 $3 $4 $5"
  r=$?
  run_test check-2-5-autoscaled-aks-apply "$r" "$1"
  r=$?
  run_test validate-azure-resources-presence "$r" "$1 $2 $3 $4 $5"
  r=0
  run_test cleanup-after-apply "$r" "$1 $2 $3 $4 $5"
  r=0
  run_test teardown-azks-module-tests-rg "$r" "$1 $2 $3 $4 $5"
  r=$?

  stop_suite test-plan "$r"
}

function init-default-config() {
  echo "# prepare test state file"
  cp "$TESTS_DIR"/tests/mocks/default-config/state.yml "$TESTS_DIR"/shared/
  echo "# will initialize config with \"docker run ... init\" command"
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    init
}

function check-default-config-content() {
  echo "# will test if file ./shared/azks/azks-config.yml exists"
  if ! test -f "$TESTS_DIR"/shared/azks/azks-config.yml; then exit 1; fi
  echo "# will test if file ./shared/azks/azks-config.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/azks/azks-config.yml "$TESTS_DIR"/mocks/default-config/config.yml
}

function init-2-5-autoscaled-aks() {
  echo "# prepare test state file"
  cp "$TESTS_DIR"/tests/mocks/config-with-variables/state.yml "$TESTS_DIR"/shared/
  echo "#	will initialize config with \"docker run ... init M_NAME=azks-module-tests M_VMS_RSA=test_vms_rsa M_ADDRESS_PREFIX=10.0.0.0/16 M_SIZE=2 M_MIN=2 M_MAX=5 M_VM_SIZE=Standard_DS2_v2 M_DISK_SIZE=36 M_AUTO_SCALING=true command\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    init \
    M_NAME=azks-module-tests \
    M_VMS_RSA=test_vms_rsa \
    M_ADDRESS_PREFIX=10.0.0.0/16 \
    M_DEFAULT_NODE_POOL="{ size: 2, min: 2, max: 5, vm_size: Standard_DS2_v2, disk_size: 36, auto_scaling: true, type: VirtualMachineScaleSets }"
}

function check-2-5-autoscaled-aks-config-content() {
  echo "#	will test if file ./shared/azks/azks-config.yml exists"
  if ! test -f "$TESTS_DIR"/shared/azks/azks-config.yml; then exit 1; fi
  echo "#	will test if file ./shared/azks/azks-config.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/azks/azks-config.yml "$TESTS_DIR"/mocks/config-with-variables/config.yml
}

function prepare-azks-module-tests-rg() {
  echo "#	will do az login"
  az login --service-principal --username "$2" --password "$3" --tenant "$5" -o none
  echo "#	will create resource group azks-module-tests-rg"
  az group create --subscription "$4" --location francecentral --name azks-module-tests-rg
  echo "#	will create vnet azks-module-tests-vnet"
  az network vnet create --subscription "$4" --resource-group azks-module-tests-rg --name azks-module-tests-vnet --address-prefix 10.0.0.0/16
}

function plan-2-5-autoscaled-aks() {
  echo "#	will plan with \"docker run ... plan M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    plan \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
}

function check-2-5-autoscaled-aks-plan() {
  echo "#	will test if file ./shared/state.yml exists"
  if ! test -f "$TESTS_DIR"/shared/state.yml; then exit 1; fi
  echo "#	will test if file ./shared/state.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/state.yml "$TESTS_DIR"/mocks/plan/state.yml
  echo "#	will test if file ./shared/azks/terraform-apply.tfplan exists"
  if ! test -f "$TESTS_DIR"/shared/azks/terraform-apply.tfplan; then exit 1; fi
  echo "#	will test if file ./shared/azks/terraform-apply.tfplan size is greater than 0"
  filesize=$(du "$TESTS_DIR"/shared/azks/terraform-apply.tfplan | cut -f1)
  if [[ ! $filesize -gt 0 ]]; then exit 1; fi
}

function teardown-azks-module-tests-rg() {
  echo "#	will do az login"
  az login --service-principal --username "$2" --password "$3" --tenant "$5" -o none
  echo "#	will delete vnet azks-module-tests-vnet"
  az network vnet delete --subscription "$4" --resource-group azks-module-tests-rg --name azks-module-tests-vnet
  echo "#	will delete resource group azks-module-tests-rg"
  az group delete --subscription "$4" --name azks-module-tests-rg --yes
}

function apply-2-5-autoscaled-aks() {
  echo "#	will apply with \"docker run ... apply M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    apply \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
}

function check-2-5-autoscaled-aks-apply() {
  echo "#	will test if file ./shared/state.yml exists"
  if ! test -f "$TESTS_DIR"/shared/state.yml; then exit 1; fi
  echo "#	will test if file ./shared/state.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/state.yml "$TESTS_DIR"/mocks/apply/state.yml
  echo "#	will test if file ./shared/azks/terraform.tfstate exists"
  if ! test -f "$TESTS_DIR"/shared/azks/terraform.tfstate; then exit 1; fi
  echo "#	will test if file ./shared/azks/terraform.tfstate size is greater than 0"
  filesize=$(du "$TESTS_DIR"/shared/azks/terraform.tfstate | cut -f1)
  if [[ ! $filesize -gt 0 ]]; then exit 1; fi
}

function validate-azure-resources-presence() {
  echo "#	will do az login"
  az login --service-principal --username "$2" --password "$3" --tenant "$5" -o none
  echo "#	will test if there is expected resource group in subscription"
  group_id=$(az group show --subscription "$4" --name azks-module-tests-rg --query id)
  if [[ -z $group_id ]]; then exit 1; fi
  echo "#	will test if there is expected amount of machines in resource group"
  aks_count=$(az aks list --subscription "$4" --resource-group azks-module-tests-rg -o yaml | yq r - --length)
  if [[ $aks_count -ne 1 ]]; then
    echo "expected 1 but got $aks_count AKSes"
    exit 1
  fi
}

function cleanup-after-apply() {
  echo "#	will apply with \"docker run ... plan-destroy M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    plan-destroy \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
  echo "#	will apply with \"docker run ... destroy M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    destroy \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
}

# K8S_VOL_PATH and K8S_HOST_PATH are variables to set up when kubernetes based build agents are in use ('docker in docker')
# K8S_VOL_PATH - volume's mount path
# K8S_HOST_PATH - shared folder located on kubernetes host (this location is used to mount in container as share)
K8S_VOL_PATH=${K8S_VOL_PATH:=""}
K8S_HOST_PATH=${K8S_HOST_PATH:=""}
TESTS_DIR_TMP="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
TESTS_DIR=${K8S_VOL_PATH:=${TESTS_DIR_TMP}}
MOUNT_DIR=${K8S_HOST_PATH:=${TESTS_DIR_TMP}}

# Create folder structure inside volume
if [[ "$K8S_VOL_PATH" == \/* ]]; then
  mkdir -p "$K8S_VOL_PATH"/shared && cp -r "$TESTS_DIR_TMP"/mocks/ "$K8S_VOL_PATH"
fi

# shellcheck disable=SC1090
source "$(dirname "$0")/suite.sh"

case $1 in
test-default-config-suite)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  test-default-config-suite "$2"
  ;;
test-config-with-variables-suite)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  test-config-with-variables-suite "$2"
  ;;
test-plan-suite)
  if [[ $# -ne 6 ]]; then
    usage
    exit 1
  fi
  test-plan-suite "$2" "$3" "$4" "$5" "$6"
  ;;
test-apply-suite)
  if [[ $# -ne 6 ]]; then
    usage
    exit 1
  fi
  test-apply-suite "$2" "$3" "$4" "$5" "$6"
  ;;
cleanup)
  if [[ $# -ne 1 ]]; then
    usage
    exit 1
  fi
  cleanup
  ;;
setup)
  if [[ $# -ne 1 ]]; then
    usage
    exit 1
  fi
  setup
  ;;
generate_junit_report)
  if [[ $# -ne 1 ]]; then
    usage
    exit 1
  fi
  generate_junit_report
  ;;
*)
  usage
  exit 1
  ;;
esac
