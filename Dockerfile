FROM golang:1.18

RUN apt-get update \
 && apt-get install -y libaio1 unzip

# c.f. https://github.com/kubo/ruby-oci8/blob/ruby-oci8-2.2.7/docs/install-instant-client.md#install-oracle-instant-client-packages
RUN mkdir -p /opt/oracle \
 && wget --quiet https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-basiclite-linux.x64-19.3.0.0.0dbru.zip \
 && wget --quiet https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-sdk-linux.x64-19.3.0.0.0dbru.zip \
 && unzip -q instantclient-basiclite-linux.x64-19.3.0.0.0dbru.zip -d /opt/oracle \
 && unzip -q instantclient-sdk-linux.x64-19.3.0.0.0dbru.zip -d /opt/oracle \
 && mv /opt/oracle/instantclient_19_3 /opt/oracle/instantclient \
 && rm *.zip

COPY _build/ubuntu/oci8.pc /usr/local/lib/pkgconfig/oci8.pc

COPY . /app

WORKDIR /app
