package dhcp

import (
	"fmt"
	"net"
	"time"

	"github.com/Its-Alex/go-postgresql-dhcp-server/database"
	"github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Handler is struct containing server information
type Handler struct {
	ip            net.IP
	leaseDuration time.Duration
	options       dhcp4.Options
}

// ServeDHCP is a state machine for dhcp response
func (h *Handler) ServeDHCP(p dhcp4.Packet, msgType dhcp4.MessageType, options dhcp4.Options) (d dhcp4.Packet) {
	mac := p.CHAddr().String()

	switch msgType {
	case dhcp4.Discover:
		logrus.Infof("Discover request from: %s", mac)
		reservation := database.GetReservationByMAC(mac)
		if reservation.IP == "" {
			logrus.Infof("Unknown server : %s", mac)
			return nil
		}
		logrus.Infof("Data retrieve from database: %s", reservation)
		h.options[dhcp4.OptionSubnetMask] = net.ParseIP(reservation.MaskSubnet)
		return dhcp4.ReplyPacket(p, dhcp4.Offer, h.ip, net.ParseIP(reservation.IP), h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp4.OptionParameterRequestList]))
	case dhcp4.Request:
		logrus.Infof("Request from: %s", mac)
		// if server, ok := options[dhcp4.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.ip) {
		// 	return nil
		// }
		reqIP := net.IP(options[dhcp4.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = net.IP(p.CIAddr())
		}
		reservation := database.GetReservationByMAC(mac)
		if reservation.IP == "" {
			logrus.Infof("Unknown server : %s", mac)
			return dhcp4.ReplyPacket(p, dhcp4.NAK, h.ip, nil, 0, nil)
		}
		logrus.Infof("Data retrieve from database: %s", reservation)
		h.options[dhcp4.OptionSubnetMask] = net.ParseIP(reservation.MaskSubnet)
		return dhcp4.ReplyPacket(p, dhcp4.ACK, h.ip, net.ParseIP(reservation.IP), h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp4.OptionParameterRequestList]))
	default:
		logrus.Infof("Not handled request from: %s", mac)
	}
	return nil
}

// Start is used by cobra to launch programm
func Start() {
	conn, err := conn.NewUDP4FilterListener(viper.GetString("interface"), fmt.Sprintf(":%s",
		viper.GetString("port"),
	))
	if err != nil {
		logrus.Fatal("NewUDP4FilterListener: ", err)
	}
	handler := &Handler{
		ip:            net.ParseIP(viper.GetString("server_ip")),
		leaseDuration: 2 * time.Hour,
		options: dhcp4.Options{
			dhcp4.OptionRouter:           []byte(viper.GetString("server_ip")),
			dhcp4.OptionDomainNameServer: []byte(viper.GetString("server_ip")),
			dhcp4.OptionBootFileName:     []byte("pxelinux.0"),
		},
	}
	logrus.Infof("dhcp4 start on interface %s and on port %s", viper.GetString("port"), viper.GetString("interface"))
	logrus.Fatal(dhcp4.Serve(conn, handler))
}
