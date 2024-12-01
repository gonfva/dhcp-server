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
	address   string
	allocated bool
	hwAddr    *net.HardwareAddr
}

func (p *AllocationPool) allocateIPAddress(hwAddr *net.HardwareAddr) (ipAddress string, err error) {
	for _, v := range *p.pool {
		if v.hwAddr != nil && slices.Equal(*v.hwAddr, *hwAddr) {
			log.Printf("Address already reserved %s for hwadd %s\n", v.address, v.hwAddr)
			return v.address, nil
		}
	}

	for idx, v := range *p.pool {
		if !v.allocated {
			v.allocated = true
			v.hwAddr = hwAddr
			(*p.pool)[idx] = v
			return v.address, nil
		}
	}
	return net.IP{}.String(), fmt.Errorf("no IP available")
}

func (p *AllocationPool) deallocateIPAddress(ipAddress net.IP) error {
	for idx, v := range *p.pool {
		if v.address == ipAddress.String() {
			log.Printf("Deallocated address %s\n", v.address)
			v.allocated = false
			v.hwAddr = nil
			(*p.pool)[idx] = v
			return nil
		}
	}
	return fmt.Errorf("address not found")
}

func getPoolAddress(poolCidr string) AllocationPool {
	ip, ipnet, err := net.ParseCIDR(poolCidr)
	if err != nil {
		log.Fatalf("Invalid CIDR: %v", err)
	} else {
		log.Printf("Valid CIDR: %v", poolCidr)

	}

	// Create a slice to store IPs
	var ips []Address

	// Iterate through all IPs in the range
	for ip := ip; ipnet.Contains(ip); inc(ip) {
		// Copy the IP to avoid reference issues
		//log.Printf("Adding address %s\n", ip)
		newAddress := Address{
			address:   ip.String(),
			allocated: false,
		}
		ips = append(ips, newAddress)
	}
	pool := AllocationPool{pool: &ips}
	return pool
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
