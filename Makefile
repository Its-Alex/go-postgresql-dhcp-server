.PHONY: enter
enter:
	docker-compose exec goworkspace bash

.PHONY: deps
deps:
	dep ensure -v --vendor-only

.PHONY: build
build:
	go build -v -o bin/go-postgresql-dhcp-server .

.PHONY: test
test:
	go test ./...

.PHONY: docker-build
docker-build:
	docker build . -t itsalex/go-postgresql-dhcp-server

.PHONY: db-schema
db-schema:
	echo "CREATE TABLE reservations(mask_subnet cidr NOT NULL, mac macaddr NOT NULL UNIQUE, ip inet NOT NULL UNIQUE);" | docker-compose exec -T --user postgres postgres psql dhcp4

.PHONY: db-fixtures
db-fixtures:
	cat sample-data.sql | docker-compose exec -T --user postgres postgres psql dhcp4

.PHONY: clean
clean:
	docker-compose stop | true
	docker-compose rm -f | true
	rm -rf ./data-dhcp4/
	rm -rf bin/
	rm -rf vendor/
	vagrant destroy -f | true
	rm .vagrant -rf
