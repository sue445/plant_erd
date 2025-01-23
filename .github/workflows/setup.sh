#!/bin/bash -xe

if [ -z "${TARGET_ARCH}" ]; then
  TARGET_ARCH="amd64"
fi

# ref. https://gist.github.com/sj26/88e1c6584397bb7c13bd11108a579746
function retry {
  local retries=$1
  shift

  local count=0
  until "$@"; do
    exit=$?
    wait=$((2 ** $count))
    count=$(($count + 1))
    if [ $count -lt $retries ]; then
      echo "Retry $count/$retries exited $exit, retrying in $wait seconds..."
      sleep $wait
    else
      echo "Retry $count/$retries exited $exit, no more retries left."
      return $exit
    fi
  done
  return 0
}

case "${RUNNER_OS}" in
"Linux")
  sudo apt-get update
  sudo apt-get install -y gcc-multilib g++-multilib libaio1t64

  case "${TARGET_ARCH}" in
  "amd64")
    sudo _build/ubuntu/setup_oracle_x64.sh
    ;;

  "386")
    sudo _build/ubuntu/setup_oracle_386.sh
    ;;

  *)
    echo "Uknown TARGET_ARCH: ${TARGET_ARCH}"
    exit 1
  ;;
  esac

  sudo mkdir -p /usr/local/lib/pkgconfig/
  sudo cp _build/ubuntu/oci8.pc /usr/local/lib/pkgconfig/oci8.pc
  ;;

"macOS")
  brew install pkg-config
  sudo _build/macos/setup_oracle.sh
  sudo mkdir -p /usr/local/lib/pkgconfig/
  sudo cp _build/macos/oci8.pc /usr/local/lib/pkgconfig/oci8.pc
  ;;

"Windows")
  retry 5 choco install --yes --allow-empty-checksums pkgconfiglite zip
  _build/windows/setup_oracle_x64.sh

  mkdir -p /usr/local/lib/pkgconfig/
  cp _build/windows/oci8.pc /usr/local/lib/pkgconfig/oci8.pc
  ;;

*)
  echo "Uknown OS: ${RUNNER_OS}"
  exit 1
  ;;
esac
