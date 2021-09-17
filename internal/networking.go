package wololo

import (
	"errors"
	"net"
)

// Byte sequence for WOL, global variable for handler and main
type WolPacket struct {
	header  [6]byte
	macAddr [16]MACAddress
}

// Extract the IPv4 address from a specified interface
func InterfaceToIp(iface string) (*net.UDPAddr, error) {
	ifc, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}

	// Get associated IPv4/IPv6 addresses
	addrs, err := ifc.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		switch at := addr.(type) {
		case *net.IPNet:
			// Only IPv4 addresses will have a default mask
			if at.IP.DefaultMask() != nil {
				return &net.UDPAddr{
					IP: at.IP,
				}, nil
			}
		}
	}

	// Couldn't find anything..
	return nil, errors.New("No suitable IP address associated to interface " + iface)
}

// Construct the raw WOL packet
func BuildWolPacket(macAddr MACAddress) *WolPacket {
	// Build the WOL packet
	var wolPacket WolPacket
	wolPacket.header = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	for i := 0; i < 16; i++ {
		wolPacket.macAddr[i] = macAddr
	}

	return &wolPacket
}