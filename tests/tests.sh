#!/usr/bin/env bash

function usage() {
  echo "usage:
    $0 cleanup
    $0 setup
    $0 generate_junit_report
    $0 test-default-config-suite [image_name]
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

function init-default-config() {
  echo "# prepare test state file"
  cp "$TESTS_DIR"/tests/mocks/default-config/state.yml "$TESTS_DIR"/shared/
  echo "# will initialize config with \"docker run ... init\" command"
  docker run --rm \
    -v "$TESTS_DIR"/shared:/shared \
    -t "$1" \
    init
}

function check-default-config-content() {
  echo "# will test if file ./shared/azks/azks-config.yml exists"
  if ! test -f "$TESTS_DIR"/shared/azks/azks-config.yml; then exit 1; fi
  echo "# will test if file ./shared/azks/azks-config.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/azks/azks-config.yml "$TESTS_DIR"/mocks/default-config/config.yml
}


TESTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

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
