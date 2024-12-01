package main

import (
	"fmt"
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

type Handler struct {
	responsePort int
	subnet       net.IP
	gateway      net.IP
	dns          net.IP
	server       net.IP
	pool         AllocationPool
}

func (h *Handler) handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	// Log the received packet for debugging
	log.Printf("Received DHCPv4 packet from %v: %+v\n", peer, m)

	switch m.MessageType() {
	case dhcpv4.MessageTypeDiscover:
		h.discoveryResponse(conn, peer, m)
	case dhcpv4.MessageTypeRequest:
		h.requestResponse(conn, peer, m)
	case dhcpv4.MessageTypeRelease:
		h.releaseAddress(conn, peer, m)
	default:
		log.Printf("Message type ignored %s\n", m.MessageType())
		return
	}

}

func (h *Handler) discoveryResponse(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Printf("DISCOVERY response %v\n", peer)
	// Allocate an IP address from your pool
	ipAddr, err := h.pool.allocateIPAddress(&m.ClientHWAddr)
	if err != nil {
		log.Println("Error allocating a new address")
		return
	}
	log.Printf("Address allocated %v\n", ipAddr)
	response, err := dhcpv4.New(
		dhcpv4.WithTransactionID(m.TransactionID),
		dhcpv4.WithHwAddr(m.ClientHWAddr),
		dhcpv4.WithYourIP(net.ParseIP(ipAddr)),
		dhcpv4.WithServerIP(h.server),
		dhcpv4.WithBroadcast(true),
		dhcpv4.WithMessageType(dhcpv4.MessageTypeOffer),
		dhcpv4.WithDNS(h.dns),
		dhcpv4.WithGatewayIP(h.gateway),
	)
	if err != nil {
		log.Println("Error generating packet")
		return
	}
	log.Printf("Sending packet %v\n", response)
	// Send the response packet
	if err := sendUDPPacket(peer, response.ToBytes()); err != nil {
		log.Printf("Error sending DHCP response: %v\n", err)
	}
}

func sendUDPPacket(addr net.Addr, message []byte) error {
	conn, err := net.Dial("udp", addr.String())
	if err != nil {
		return fmt.Errorf("error creating UDP connection: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(message)
	if err != nil {

		return fmt.Errorf("error sending UDP packet: %v", err)
	}

	log.Println("UDP packet sent successfully!")
	return nil
}

func (h *Handler) requestResponse(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Println("REQUEST response")
	//Wed don't confirm
	response, err := dhcpv4.New(
		dhcpv4.WithTransactionID(m.TransactionID),
		dhcpv4.WithHwAddr(m.ClientHWAddr),
		dhcpv4.WithServerIP(h.server),
		dhcpv4.WithClientIP(m.ClientIPAddr),
		dhcpv4.WithMessageType(dhcpv4.MessageTypeAck),
	)
	if err != nil {
		log.Println("Error generating packet")
		return
	}

	// Send the response packet
	if err := sendUDPPacket(peer, response.ToBytes()); err != nil {
		log.Printf("Error sending DHCP response: %v\n", err)
	}
	log.Printf("Packet sent %v\n", response)
}

func (h *Handler) releaseAddress(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Println("Deallocating address")
	err := h.pool.deallocateIPAddress(m.ClientIPAddr)
	if err != nil {
		log.Println("Error allocating a new address")
		return
	}
	response, err := dhcpv4.New(
		dhcpv4.WithTransactionID(m.TransactionID),
		dhcpv4.WithHwAddr(m.ClientHWAddr),
		dhcpv4.WithServerIP(h.server),
		dhcpv4.WithClientIP(m.ClientIPAddr),
		dhcpv4.WithMessageType(dhcpv4.MessageTypeAck),
	)
	if err != nil {
		log.Println("Error generating packet")
		return
	}

	// Send the response packet
	if err := sendUDPPacket(peer, response.ToBytes()); err != nil {
		log.Printf("Error sending DHCP response: %v\n", err)
	}
	log.Printf("Packet sent %v\n", response)
}
