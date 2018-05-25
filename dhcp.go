package main

import (
	"fmt"
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Handler is struct containing server information
type Handler struct {
	ip            net.IP
	leaseDuration time.Duration
	options       dhcp.Options
}

// ServeDHCP is a state machine for dhcp response
func (h *Handler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	mac := p.CHAddr().String()

	switch msgType {
	case dhcp.Discover:
		reservation := GetReservationByMAC(mac)
		if reservation.IP == "" {
			logrus.Debugf("Unknown server : %s", mac)
			return nil
		}
		return dhcp.ReplyPacket(p, dhcp.Offer, h.ip, net.ParseIP(reservation.IP), h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
	case dhcp.Request:
		if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.ip) {
			return nil
		}
		reqIP := net.IP(options[dhcp.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = net.IP(p.CIAddr())
		}
		reservation := GetReservationByMAC(mac)
		if reservation.IP == "" {
			logrus.Debugf("Unknown server : %s", mac)
			return dhcp.ReplyPacket(p, dhcp.NAK, h.ip, nil, 0, nil)
		}
		return dhcp.ReplyPacket(p, dhcp.ACK, h.ip, net.ParseIP(reservation.IP), h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
	}
	return nil
}

// dhcpStart is used by cobra to launch programm
func dhcpStart(cmd *cobra.Command, args []string) {
	conn, err := conn.NewUDP4FilterListener(viper.GetString("interface"), fmt.Sprintf(":%s",
		viper.GetString("port"),
	))
	if err != nil {
		logrus.Fatal(err)
	}
	handler := &Handler{
		ip:            net.ParseIP(viper.GetString("server_ip")),
		leaseDuration: 2 * time.Hour,
		options: dhcp.Options{
			dhcp.OptionRouter:           []byte(viper.GetString("server_ip")),
			dhcp.OptionDomainNameServer: []byte(viper.GetString("server_ip")),
			dhcp.OptionBootFileName:     []byte("pxelinux.0"),
		},
	}
	logrus.Info("dhcp4 start")
	logrus.Fatal(dhcp.Serve(conn, handler))
}
