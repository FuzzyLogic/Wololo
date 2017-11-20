package main

import (
    "testing"
    "net"
    "bytes"
)

// Start a UDP server and wait for a message on port 7.
// Check if the received sequence corresponds to the expected
// WOL packet.
func TestWololo(t *testing.T) {
    // The expected WOL byte sequence for the MAC address 00:11:22:33:44:55
    wolExpected := []byte{
        0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
    }

    // Resolve UDP address
    servAddr, err := net.ResolveUDPAddr("udp",":7")
    if err != nil {
        t.Errorf(err.Error())
    }

    // Listen on port 7 for WOL signal
    servCon, err := net.ListenUDP("udp", servAddr)
    if err != nil {
        t.Errorf(err.Error())
    }
    defer servCon.Close()

    // Receive data
    wolBuf := make([]byte, 102)
    n, _, err := servCon.ReadFromUDP(wolBuf)
    if err != nil {
        t.Errorf(err.Error())
    }

    // Check number of bytes
    if n != 102 {
        t.Error("Error: Incorrect number of bytes received")
    }

    // Check received sequence
    if bytes.Compare(wolExpected, wolBuf) != 0 {
        t.Error("Error: Incorrect sequence received")
    }
}