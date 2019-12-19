#!/bin/bash -xe

case "${RUNNER_OS}" in
"Linux")
  sudo _build/ubuntu/setup_oracle.sh
  sudo mkdir -p /usr/local/lib/pkgconfig/
  sudo cp _build/ubuntu/oci8.pc /usr/local/lib/pkgconfig/oci8.pc
  ;;

"macOS")
  brew install pkg-config
  sudo _build/macos/setup_oracle.sh
  sudo mkdir -p /usr/local/lib/pkgconfig/
  sudo cp _build/macos/oci8.pc /usr/local/lib/pkgconfig/oci8.pc
  ;;

*)
  echo "Uknown OS: ${RUNNER_OS}"
  exit 1
  ;;
esac
