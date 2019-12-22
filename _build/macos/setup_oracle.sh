#!/bin/bash -xe

# c.f. https://github.com/kubo/ruby-oci8/blob/ruby-oci8-2.2.7/docs/install-on-osx.md#install-oracle-instant-client-manually
mkdir -p /opt/oracle
wget --quiet --tries=0 https://download.oracle.com/otn_software/mac/instantclient/193000/instantclient-basiclite-macos.x64-19.3.0.0.0dbru.zip
wget --quiet --tries=0 https://download.oracle.com/otn_software/mac/instantclient/193000/instantclient-sdk-macos.x64-19.3.0.0.0dbru.zip
unzip -q instantclient-basiclite-macos.x64-19.3.0.0.0dbru.zip -d /opt/oracle
unzip -q instantclient-sdk-macos.x64-19.3.0.0.0dbru.zip -d /opt/oracle
