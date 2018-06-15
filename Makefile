.PHONY: enter
enter:
	docker-compose exec goworkspace bash -c "export COLUMNS=`tput cols`; export LINES=`tput lines`; exec bash"

.PHONY: deps
deps:
	dep ensure -v --vendor-only

.PHONY: build
build:
	go build -v -o bin/go-postgresql-dhcp-server .

.PHONY: test
test:
	go test ./...

.PHONY: watch-bin
watch-bin:
	reflex -r "^bin/go-postgresql-dhcp-server$$" -s -- /vagrant/bin/go-postgresql-dhcp-server

.PHONY: watch-go
watch-go:
	reflex -r "\.go$$" -s -- make build

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
	vagrant destroy -f | true
	rm -rf ./data-dhcp4/ bin/ /vendor .vagrant
