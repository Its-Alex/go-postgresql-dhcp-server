DELETE FROM reservations;
INSERT INTO reservations (mask_subnet, mac, ip) VALUES ('255.255.255.0'::cidr, '08:00:27:8e:15:8a'::macaddr, '192.168.0.11'::inet);
INSERT INTO reservations (mask_subnet, mac, ip) VALUES ('255.255.255.0'::cidr, '1c:ed:0c:0a:88:53'::macaddr, '192.168.0.12'::inet);
INSERT INTO reservations (mask_subnet, mac, ip) VALUES ('255.255.255.0'::cidr, 'ec:b5:0a:fe:a9:62'::macaddr, '192.168.0.13'::inet);
