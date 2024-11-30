package main

import (
	"fmt"
	"log"
	"net"
	"slices"
)

type AllocationPool struct {
	pool *[]Address
}

type Address struct {
	address   net.IP
	allocated bool
}

func (p *AllocationPool) allocateIPAddress() (ipAddress net.IP, err error) {
	for _, v := range *p.pool {
		if !v.allocated {
			v.allocated = true
			return v.address, nil
		}
	}
	return net.IP{}, fmt.Errorf("Not implemented")
}

func (p *AllocationPool) deallocateIPAddress(ipAddress net.IP) error {
	for _, v := range *p.pool {
		if slices.Equal(v.address, ipAddress) {
			v.allocated = false
			return nil
		}
	}
	return fmt.Errorf("Address not found")
}

func getPoolAddress(poolCidr string) AllocationPool {
	ip, ipnet, err := net.ParseCIDR(poolCidr)
	if err != nil {
		log.Fatalf("Invalid CIDR: %v", err)
	} else {
		log.Printf("Valid CIDR: %v", ipnet)

	}

	// Create a slice to store IPs
	var ips []Address

	// Iterate through all IPs in the range
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		// Copy the IP to avoid reference issues
		ips = append(ips, Address{
			address:   ip,
			allocated: false,
		})
	}

	return AllocationPool{pool: &ips}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
