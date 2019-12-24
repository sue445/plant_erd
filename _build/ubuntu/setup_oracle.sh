#!/bin/bash -xe

# c.f. https://github.com/kubo/ruby-oci8/blob/ruby-oci8-2.2.7/docs/install-instant-client.md#install-oracle-instant-client-packages
mkdir -p /opt/oracle
wget --quiet --tries=0 https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-basiclite-linux.x64-19.3.0.0.0dbru.zip
wget --quiet --tries=0 https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-sdk-linux.x64-19.3.0.0.0dbru.zip
unzip -q instantclient-basiclite-linux.x64-19.3.0.0.0dbru.zip -d /opt/oracle
unzip -q instantclient-sdk-linux.x64-19.3.0.0.0dbru.zip -d /opt/oracle
mv /opt/oracle/instantclient_19_3 /opt/oracle/instantclient
