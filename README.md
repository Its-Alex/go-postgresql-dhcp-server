# DHCP4 reservation

This project aim to launch a dhcp4 server with only reservation capacibility.
A postgres database contain reservation info, all database entry will
be allocated and other will be ignored.


## Quick start

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
$ vagrant up pxe_server
$ vagrant ssh pxe_server
$ sudo /vagrant/bin/dhcp-server --help
```
