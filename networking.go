package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// Byte sequence for WOL, global variable for handler and main
type WolPacket struct {
	header  [6]byte
	macAddr [16]MACAddress
}

// Extrace the IPv4 address from a specified interface
func interfaceToIp(iface string) (*net.UDPAddr, error) {
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
func buildWolPacket(config WololoConfig) *WolPacket {
	// Build the WOL packet
	var wolPacket WolPacket
	wolPacket.header = [6]byte{0, 0, 0, 0, 0, 0}
	for i := 0; i < 16; i++ {
		wolPacket.macAddr[i] = globalConfig.macAddr
	}

	return &wolPacket
}

// HTTP handler function to handle requests.
// This handler will send the WOL packet to the configured destination.
// The response will contain information on whether the packet was transmitted or an error occurred.
func wolHandler(respWr http.ResponseWriter, req *http.Request) {
	wolPacket := buildWolPacket(globalConfig)

	// WOL needs to be broadcasted, create UDPAddress from config
	bcastUDPAddr, err := net.ResolveUDPAddr("udp", globalConfig.udpBcastAddr)
	if err != nil {
		writeToLog(globalLog, "Unable to obtain UDP broadcast address object")
		fmt.Fprintf(respWr, "Unable to obtain UDP broadcast address object")
		fmt.Fprintf(respWr, "Error: "+err.Error())
		return
	}

	// Get local IP address from specified interfaces
	localUDPAddr, err := interfaceToIp(globalConfig.iface)
	if err != nil {
		writeToLog(globalLog, "Unable to get local IP address from interface "+globalConfig.iface)
		fmt.Fprintf(respWr, "Unable to get local IP address from interface\n"+globalConfig.iface)
		fmt.Fprintf(respWr, "Error: "+err.Error()+"\n")
		return
	}

	// Open UDP connection to send WOL packet
	con, err := net.DialUDP("udp", localUDPAddr, bcastUDPAddr)
	if err != nil {
		writeToLog(globalLog, "Unable to create UDP connection")
		fmt.Fprintf(respWr, "Unable to create UDP connection\n")
		fmt.Fprintf(respWr, "Error: "+err.Error()+"\n")
		return
	}
	defer con.Close()

	// Broadcast the WOL packet
	var packetBuf bytes.Buffer
	binary.Write(&packetBuf, binary.BigEndian, wolPacket)
	bytesWritten, err := con.Write(packetBuf.Bytes())
	if err != nil {
		writeToLog(globalLog, "Error sending WOL packet")
		fmt.Fprintf(respWr, "Error sending WOL packet\n")
		fmt.Fprintf(respWr, "Error: "+err.Error()+"\n")
	} else if bytesWritten != 102 {
		// Not an error but something went wrong - inform user
		writeToLog(globalLog, "Warning: WOL packet transmission may have been incomplete")
		fmt.Fprintf(respWr, "Warning: WOL packet transmission may have been incomplete\n")
	} else {
		// Notify user that WOL packet was sent
		writeToLog(globalLog, "WOL packet sent")
		fmt.Fprintf(respWr, "Device is off...\nWOLOLO\nDevice is on!")
	}
}
