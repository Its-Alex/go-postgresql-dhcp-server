package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Used by sqlx as postgres driver
	"github.com/sirupsen/logrus"
)

// Db variable used to store global database
var Db *sqlx.DB

// Reservation structure in database
type Reservation struct {
	MaskSubnet string `db:"mask_subnet" json:"mask_subnet"`
	MAC        string `db:"mac" json:"mac"`
	IP         string `db:"ip" json:"ip"`
}

// GetReservationByMAC get a reservation from his mac address
func GetReservationByMAC(mac string) Reservation {
	var reservation Reservation
	err := Db.Get(&reservation, "SELECT * FROM reservations WHERE mac = $1", mac)
	if err != nil {
		logrus.Error(err)
	}
	return reservation
}
