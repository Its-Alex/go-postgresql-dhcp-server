.PHONY: deps
deps:
	dep ensure -v --vendor-only

.PHONY: build
build:
	go build .

.PHONY: test
test:
	go test ./...

.PHONY: docker-build
docker-build:
	docker build . -t itsalex/dhcp4-reservation

.PHONY: db-schema
db-schema:
	echo "CREATE TABLE reservations(mask_subnet cidr NOT NULL, mac macaddr NOT NULL UNIQUE, ip inet NOT NULL UNIQUE);" | docker-compose exec -T --user postgres postgres psql dhcp4
	
.PHONY: db-fixtures
db-fixtures:
	echo "INSERT INTO reservations (mask_subnet, mac, ip) VALUES ('255.255.255.0'::cidr, '7b:31:15:6c:80:29'::macaddr, '192.168.0.11'::inet);" | docker-compose exec -T --user postgres postgres psql dhcp4
	echo "INSERT INTO reservations (mask_subnet, mac, ip) VALUES ('255.255.255.0'::cidr, '1c:ed:0c:0a:88:53'::macaddr, '192.168.0.12'::inet);" | docker-compose exec -T --user postgres postgres psql dhcp4
	echo "INSERT INTO reservations (mask_subnet, mac, ip) VALUES ('255.255.255.0'::cidr, 'ec:b5:0a:fe:a9:62'::macaddr, '192.168.0.13'::inet);" | docker-compose exec -T --user postgres postgres psql dhcp4