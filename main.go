package main

import (
	"flag"
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

func handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	// this function will just print the received dhcpv4 message, without replying
	log.Print(m.Summary())
}

func main() {
	var listeningIp string
	var listeningPort int
	var cidr string

	// flags declaration using flag package
	flag.StringVar(&listeningIp, "l", "0.0.0.0", "Listening address")
	flag.IntVar(&listeningPort, "p", dhcpv4.ServerPort, "Listening port")
	flag.StringVar(&cidr, "s", "192.168.1.0/24", "Server range. A CIDR with the IPs that can be used to give out address")
	flag.Parse()

	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatalf("Invalid CIDR: %v", err)
	} else {
		log.Printf("Valid CIDR: %v", ipnet)
	}

	laddr := &net.UDPAddr{
		IP:   net.ParseIP(listeningIp),
		Port: listeningPort,
	}
	log.Printf("listening on %s and port %d", laddr.IP, laddr.Port)
	server, err := server4.NewServer("", laddr, handler)
	if err != nil {
		log.Fatal(err)
	}

	server.Serve()
}
