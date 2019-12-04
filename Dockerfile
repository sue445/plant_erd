FROM golang:1.13

# RUN apt-get update \
#  && apt-get install -y default-mysql-client

COPY . /app

WORKDIR /app
