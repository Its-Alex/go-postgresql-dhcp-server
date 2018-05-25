package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Used by sqlx as postgres driver
)

// Db variable used to store global database
var Db *sqlx.DB

// Reservation structure in database
type Reservation struct {
	MaskSubnet string `json:"subnet_mask"`
	IP         string `json:"ip"`
	MAC        string `json:"mac"`
}

// GetReservationByMAC get a reservation from his mac address
func GetReservationByMAC(mac string) Reservation {
	var reservation Reservation
	Db.Get(&reservation, "SELECT host(mask_subnet), text(mac), text(ip) FROM reservations WHERE mac == ?", mac)
	return reservation
}
