package main

import (
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

type Handler struct {
	subnet  net.IP
	gateway net.IP
	dns     net.IP
	server  net.IP
	pool    AllocationPool
}

func (h *Handler) handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	// Log the received packet for debugging
	log.Printf("Received DHCPv4 packet from %v: %+v\n", peer, m)

	switch m.MessageType() {
	case dhcpv4.MessageTypeDiscover:
		h.proposeAddress(conn, peer, m)
	case dhcpv4.MessageTypeRequest:
		h.confirmAddress(conn, peer, m)
	case dhcpv4.MessageTypeRelease:
		h.releaseAddress(conn, peer, m)
	default:
		log.Printf("Message type ignored %s\n", m.MessageType())
		return
	}

}

func (h *Handler) proposeAddress(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Println("Proposing address")
	// Allocate an IP address from your pool
	ipAddr, err := h.pool.allocateIPAddress()
	if err != nil {
		log.Println("Error allocating a new address")
		return
	}
	response, err := dhcpv4.New(
		dhcpv4.WithTransactionID(m.TransactionID),
		dhcpv4.WithHwAddr(m.ClientHWAddr),
		dhcpv4.WithYourIP(ipAddr),
		dhcpv4.WithServerIP(h.server),
		dhcpv4.WithMessageType(dhcpv4.MessageTypeOffer),
		dhcpv4.WithDNS(h.dns),
		dhcpv4.WithGatewayIP(h.gateway),
	)
	if err != nil {
		log.Println("Error generating packet")
		return
	}

	// Send the response packet
	if _, err := conn.WriteTo(response.ToBytes(), peer); err != nil {
		log.Printf("Error sending DHCP response: %v\n", err)
	}
	log.Printf("Packet sent %v\n", response)
}

func (h *Handler) confirmAddress(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Println("Confirming address")
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
	if _, err := conn.WriteTo(response.ToBytes(), peer); err != nil {
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
	if _, err := conn.WriteTo(response.ToBytes(), peer); err != nil {
		log.Printf("Error sending DHCP response: %v\n", err)
	}
	log.Printf("Packet sent %v\n", response)
}
