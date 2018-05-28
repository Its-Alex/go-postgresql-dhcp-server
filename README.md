[![Docker Automated build](https://img.shields.io/docker/automated/itsalex/go-postgresql-dhcp-server.svg)](https://hub.docker.com/r/itsalex/go-postgresql-dhcp-server/)

# go-postgresql-dhcp-server

This project aim to launch a dhcp4 server with only reservation capacibility.
A postgres database contain reservation info, all database entry will
be allocated with an ip and other will be ignored.


## Work in docker container

```
$ docker-compose up -d
$ make db-schema db-fixtures
$ make enter
# make deps
# make build
# make test
```

## Test in Vagrant environment

```
$ make build
$ vagrant up pxe_server
$ vagrant ssh pxe_server
$ sudo su
# export DHCP4_INTERFACE=enp0s8
# export DHCP4_PSQL_ADDR=10.0.2.2
# /vagrant/bin/go-postgresql-dhcp-server
```

In another terminal:

```
$ vagrant destroy blank_server -f && vagrant up blank_server
```
