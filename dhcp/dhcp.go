package dhcp

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Its-Alex/go-postgresql-dhcp-server/database"
	log "github.com/Its-Alex/go-postgresql-dhcp-server/log"
	"github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
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
	logger := log.Logger.WithFields(logrus.Fields{
		"incoming_request_mac":  mac,
		"incoming_request_type": fmt.Sprintf("%s", msgType),
	})
	logger.Infof("New request")

	switch msgType {
	case dhcp4.Discover:
		reservation := database.GetReservationByMAC(mac)
		if reservation.IP == "" {
			logger.Debug("Server not found in database")
			return nil
		}
		logger.Debugf("Data retrieve from database: %s", reservation)
		// h.options[dhcp4.OptionSubnetMask] = net.ParseIP(reservation.MaskSubnet)
		return dhcp4.ReplyPacket(p, dhcp4.Offer, h.ip, net.ParseIP(reservation.IP), h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp4.OptionParameterRequestList]))
	case dhcp4.Request:
		// Get ip requested
		reqIP := net.IP(options[dhcp4.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = net.IP(p.CIAddr())
		}

		// Ensure that this request is for this dhcp
		if server, ok := options[dhcp4.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.ip) && bytes.Compare(server, []byte{0, 0, 0, 0}) != 0 {
			logger.Debugf("Wrong option server identifier: %s != %s", net.IP(server), net.IP(h.ip))
			return nil
		}

		// Ensure that request ip is a true ip (no 0.0.0.0 or 192.168.0)
		if len(reqIP) != 4 && reqIP.Equal(net.IPv4zero) {
			logger.Debugf("Wrong ip asked by server: %s", reqIP)
			return dhcp4.ReplyPacket(p, dhcp4.NAK, h.ip, nil, 0, nil)
		}

		// Get datas from database
		reservation := database.GetReservationByMAC(mac)
		// Ensure that this server is in our database
		if reservation.IP == "" {
			logger.Debug("Server not found in database")
			return dhcp4.ReplyPacket(p, dhcp4.NAK, h.ip, nil, 0, nil)
		}

		logger.Debugf("Data retrieve from database: %s", reservation)
		// h.options[dhcp4.OptionSubnetMask] = net.ParseIP(reservation.MaskSubnet)
		return dhcp4.ReplyPacket(p, dhcp4.ACK, h.ip, net.ParseIP(reservation.IP), h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp4.OptionParameterRequestList]))
	default:
		logger.Infof("Not handled request")
	}
	return dhcp4.ReplyPacket(p, dhcp4.NAK, h.ip, nil, 0, nil)
}

// Start is used by cobra to launch programm
func Start() {
	conn, err := conn.NewUDP4FilterListener(viper.GetString("interface"), fmt.Sprintf(":%s",
		viper.GetString("port"),
	))
	if err != nil {
		log.Logger.Fatal("NewUDP4FilterListener: ", err)
	}
	handler := &Handler{
		ip:            net.ParseIP(viper.GetString("server_ip")),
		leaseDuration: 2 * time.Hour,
		options: dhcp4.Options{
			dhcp4.OptionSubnetMask:       []byte{255, 255, 255, 0},
			dhcp4.OptionRouter:           []byte(viper.GetString("server_ip")),
			dhcp4.OptionDomainNameServer: []byte(viper.GetString("server_ip")),
			dhcp4.OptionServerIdentifier: []byte{192, 168, 0, 254},
			dhcp4.OptionBootFileName:     []byte("pxelinux.0"),
		},
	}
	log.Logger.Infof("dhcp4 start on interface %s and on port %s", viper.GetString("port"), viper.GetString("interface"))
	log.Logger.Fatal(dhcp4.Serve(conn, handler))
}
