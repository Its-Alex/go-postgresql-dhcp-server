// Example of minimal DHCP server:
package main

import (
	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"

	"log"
	"net"
	"time"
)

// DHCPHandler is struct containing needed information
type DHCPHandler struct {
	ip            net.IP
	leaseDuration time.Duration // Lease period
	options       dhcp.Options  // Options to send to DHCP Clients
}

// Example using DHCP with a single network interface device
func main() {
	serverIP := net.IP{192, 168, 0, 254}
	conn, err := conn.NewUDP4FilterListener("enp0s8", ":67")
	if err != nil {
		log.Fatal("Can't setup interface", err)
	}
	log.Fatal(dhcp.Serve(conn, &DHCPHandler{
		ip:            serverIP,
		leaseDuration: 2 * time.Hour,
		options: dhcp.Options{
			dhcp.OptionSubnetMask:       net.IP{255, 255, 255, 0},
			dhcp.OptionRouter:           []byte(serverIP),
			dhcp.OptionDomainNameServer: []byte(serverIP),
			dhcp.OptionBootFileName:     []byte("pxelinux.0"),
		},
	}))
}

// ServeDHCP is a state machine for dhcp response
func (h *DHCPHandler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	switch msgType {

	case dhcp.Discover:
		return dhcp.ReplyPacket(p, dhcp.Offer, h.ip, net.IP{192, 168, 0, 11}, h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))

	case dhcp.Request:
		// if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.ip) {
		// 	return nil // Message not for this dhcp server
		// }
		// reqIP := net.IP(options[dhcp.OptionRequestedIPAddress])
		// if reqIP == nil {
		// 	reqIP = net.IP(p.CIAddr())
		// }

		// if nic := p.CHAddr().String(); nic == "" /*Address connu*/ {
		return dhcp.ReplyPacket(p, dhcp.ACK, h.ip, net.IP{192, 168, 0, 11}, h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
		// }
		// return dhcp.ReplyPacket(p, dhcp.NAK, h.ip, nil, 0, nil)
	}
	return nil
}
