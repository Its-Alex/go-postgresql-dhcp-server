[![Docker Automated build](https://img.shields.io/docker/automated/itsalex/go-postgresql-dhcp-server.svg)](https://hub.docker.com/r/itsalex/go-postgresql-dhcp-server/)

# go-postgresql-dhcp-server

This project aim to launch a dhcp4 server with only reservation capacibility.
A postgres database contain reservation info, all database entry will
be allocated with an ip and other will be ignored.


## Work in docker container

### Initialisation

```
$ docker-compose up -d
$ make db-schema db-fixtures
$ make enter
```

### Commands

You can now use commands

Get dependencies

```
# make deps
```

Build project

```
# make build
```

Run tests

```
# make test
```

## Test in Vagrant environment

You must at least run postgres container and populate it with schema and datas

```
$ make build
$ vagrant up pxe_server
$ vagrant ssh pxe_server
$ sudo su
# /vagrant/bin/go-postgresql-dhcp-server
```

Default env variable are set in vagrantfile with provisions [here](/Vagrantfile#L15)

In another terminal:

```
$ vagrant destroy blank_server -f && vagrant up blank_server
```

## Tips and bugs

Sometimes vagrant don't want to sync your files you can run this in terminal:

```
$ vagrant rsync-auto
```

## Source and documentation

* dhcp 
    * https://godoc.org/go.universe.tf/netboot/dhcp4
    * https://github.com/krolaw/dhcp4

## License

[AGPL](https://fr.wikipedia.org/wiki/GNU_Affero_General_Public_License)