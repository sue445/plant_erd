#!/bin/bash -xe

# c.f. https://www.oracle.com/database/technologies/instant-client/macos-intel-x86-downloads.html
mkdir -p /opt/oracle
wget --quiet --tries=0 https://download.oracle.com/otn_software/mac/instantclient/193000/instantclient-basiclite-macos.x64-19.3.0.0.0dbru.zip
wget --quiet --tries=0 https://download.oracle.com/otn_software/mac/instantclient/193000/instantclient-sdk-macos.x64-19.3.0.0.0dbru.zip
unzip -q instantclient-basiclite-macos.x64-19.3.0.0.0dbru.zip -d /opt/oracle
unzip -q instantclient-sdk-macos.x64-19.3.0.0.0dbru.zip -d /opt/oracle
mv /opt/oracle/instantclient_19_3 /opt/oracle/instantclient
