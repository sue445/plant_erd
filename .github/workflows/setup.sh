#!/bin/bash -xe

if [ -z "${TARGET_ARCH}" ]; then
  TARGET_ARCH="amd64"
fi

case "${RUNNER_OS}" in
"Linux")
  sudo apt-get update
  sudo apt-get install -y gcc-multilib g++-multilib

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
  if [ ! -e tmp/pacman ]; then
    # Save only pacman files to cache
    git clone --depth=1 --quiet https://github.com/git-for-windows/git-sdk-64

    mkdir -p tmp/pacman/usr/bin
    cp -R git-sdk-64/usr/bin/pacman* tmp/pacman/usr/bin

    mkdir -p tmp/pacman/etc
    cp -R git-sdk-64/etc/pacman.* tmp/pacman/etc

    mkdir -p tmp/pacman/var/lib
    cp -R git-sdk-64/var/lib/pacman tmp/pacman/var/lib

    mkdir -p tmp/pacman/usr/share/makepkg
    cp -R git-sdk-64/usr/share/makepkg/util* tmp/pacman/usr/share/makepkg

    du -sm tmp/pacman
  fi

  # Install pacman
  cp -R tmp/pacman/usr/bin/* /usr/bin/
  cp -R tmp/pacman/etc/* /etc/
  mkdir -p /var/lib/
  cp -R tmp/pacman/var/lib/* /var/lib/
  cp -R tmp/pacman/usr/share/makepkg/* /usr/share/makepkg/

  pacman --database --check

  curl -L https://raw.githubusercontent.com/git-for-windows/build-extra/master/git-for-windows-keyring/git-for-windows.gpg | pacman-key --add -
  pacman-key --lsign-key 1A9F3986

  case "${TARGET_ARCH}" in
  "amd64")
    pacman -S --noconfirm mingw-w64-x86_64-pkg-config
    _build/windows/setup_oracle_x64.sh
    ;;

  "386")
    pacman -S --noconfirm mingw-w64-i686-pkg-config
    # _build/windows/setup_oracle_386.sh
    ;;

  *)
    echo "Uknown TARGET_ARCH: ${TARGET_ARCH}"
    exit 1
  ;;
  esac

  mkdir -p /usr/local/lib/pkgconfig/
  cp _build/windows/oci8.pc /usr/local/lib/pkgconfig/oci8.pc
  ;;

*)
  echo "Uknown OS: ${RUNNER_OS}"
  exit 1
  ;;
esac
