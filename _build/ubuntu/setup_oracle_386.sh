#!/bin/bash -xe

# c.f. https://www.oracle.com/database/technologies/instant-client/linux-x86-32-downloads.html
mkdir -p /opt/oracle
wget --quiet --tries=0 https://download.oracle.com/otn_software/linux/instantclient/195000/instantclient-basiclite-linux-19.5.0.0.0dbru.zip
wget --quiet --tries=0 https://download.oracle.com/otn_software/linux/instantclient/195000/instantclient-sdk-linux-19.5.0.0.0dbru.zip
unzip -q instantclient-basiclite-linux-19.5.0.0.0dbru.zip -d /opt/oracle
unzip -q instantclient-sdk-linux-19.5.0.0.0dbru.zip -d /opt/oracle
mv /opt/oracle/instantclient_19_5 /opt/oracle/instantclient

# c.f. https://forums.oracle.com/ords/apexds/post/instant-client-on-ubuntu-24-04-noble-numbat-7244
ln -s /usr/lib/$(uname -m)-linux-gnu/libaio.so.1t64 /opt/oracle/instantclient/libaio.so.1
