#!/bin/bash -xe

# c.f. https://www.oracle.com/database/technologies/instant-client/winx64-64-downloads.html
mkdir -p C:/opt/oracle
curl -OL https://download.oracle.com/otn_software/nt/instantclient/19500/instantclient-basiclite-windows.x64-19.5.0.0.0dbru.zip
curl -OL https://download.oracle.com/otn_software/nt/instantclient/19500/instantclient-sdk-windows.x64-19.5.0.0.0dbru.zip
unzip -q instantclient-basiclite-windows.x64-19.5.0.0.0dbru.zip -d C:/opt/oracle
unzip -q instantclient-sdk-windows.x64-19.5.0.0.0dbru.zip -d C:/opt/oracle
mv C:/opt/oracle/instantclient_19_5 C:/opt/oracle/instantclient
