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
$ sudo su
# export DHCP4_INTERFACE=enp0s8
# export DHCP4_PSQL_ADDR=10.0.2.2
# /vagrant/bin/dhcp-server
root@pxe-server:/home/vagrant# /vagrant/bin/dhcp-server
{"level":"info","msg":"dhcp4 start on interface 67 and on port enp0s8","time":"2018-05-28T13:25:09Z"}
```

In another terminal:

```
$ vagrant up blank_server
```

To retry:

```
$ vagrant destroy blank_server -f
$ vagrant up blank_server
```
