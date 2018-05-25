package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Used by sqlx as postgres driver
)

// Db variable used to store global database
var Db *sqlx.DB

// Device structure in database
type Device struct {
	MaskSubnet string
	IP         string
	MAC        string
}

// GetDeviceByMAC get a device from his mac address
func GetDeviceByMAC(mac string) Device {
	var device Device
	Db.Get(&device, "SELECT * FROM devices WHERE mac == ?", mac)
	return device
}
