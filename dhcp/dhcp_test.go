package dhcp_test

import (
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Its-Alex/go-postgresql-dhcp-server/dhcp"
	"github.com/krolaw/dhcp4"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	log.Println("MAIN TEST")
	viper.Set("port", "10000")
	viper.Set("interface", "en0")
	viper.Set("server_ip", "192.168.0.254")
	go dhcp.Start()
	time.Sleep(5 * time.Second)
	os.Exit(m.Run())
}

func TestClientBroadcast(t *testing.T) {
	log.Println("CLIENT BROACAST")
	mac, err := net.ParseMAC("7b:31:15:6c:80:29")
	if err != nil {
		t.Error(err)
	}

	packet := dhcp4.RequestPacket(
		dhcp4.Discover,
		mac,
		net.IP{192, 168, 0, 11},
		dhcp4.Packet{},
		true,
		[]dhcp4.Option{},
	)

	conn, err := net.Dial("udp4", "localhost:10000")
	if err != nil {
		t.Error(err)
	}

	conn.Write(packet)
	response := make([]byte, 4096)

	_, err = conn.Read(response)
	if err != nil {
		t.Error(err)
	}
	t.Error(response)
}
