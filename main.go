package main

import (
	"flag"
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

func main() {
	var listeningIp string
	var listeningPort int
	var pool string
	var subnet_cidr string
	var gateway_address string
	var dns_address string
	var server_address string

	// flags declaration using flag package
	flag.StringVar(&listeningIp, "listenIp", "0.0.0.0", "Listening address")
	flag.IntVar(&listeningPort, "listenPort", dhcpv4.ServerPort, "Listening port")
	flag.StringVar(&pool, "pool", "192.168.1.64/25", "Pool range. A CIDR with the IPs that can be used to give out address. For example from 192.168.1.64 to 192.168.1.127")
	flag.StringVar(&subnet_cidr, "subnet", "192.168.1.0/24", "Subnet range. A CIDR the subnet range")
	flag.StringVar(&gateway_address, "gateway", "192.168.1.1", "The gateway address")
	flag.StringVar(&server_address, "server", "192.168.1.2", "The server address")
	flag.StringVar(&dns_address, "dns", "8.8.8.8", "The DNS entry")
	flag.Parse()

	laddr := &net.UDPAddr{
		IP:   net.ParseIP(listeningIp),
		Port: listeningPort,
	}
	poolArray := getPoolAddress(pool)

	h := &Handler{
		subnet:  net.ParseIP(subnet_cidr),
		gateway: net.ParseIP(gateway_address),
		dns:     net.ParseIP(dns_address),
		server:  net.ParseIP(server_address),
		pool:    poolArray,
	}
	log.Printf("listening on %s and port %d", laddr.IP, laddr.Port)
	log.Printf("parameters %v", h)
	server, err := server4.NewServer("", laddr, h.handler)
	if err != nil {
		log.Fatal(err)
	}

	server.Serve()
}
