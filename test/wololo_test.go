package main

import (
    "testing"
    "net"
    "bytes"
    "net/http"
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
    cn := make(chan int)
    wolBuf := make([]byte, 102)
    go func(servCon *net.UDPConn, cn chan int, buf []byte) {
        tmpN, _, err := servCon.ReadFromUDP(buf)
        if err != nil {
            t.Errorf(err.Error())
        }

        // Output number of read bytes to channel
        cn <- tmpN
    }(servCon, cn, wolBuf)

    // Send an HTTP GET to server
    _, err = http.Get("http://172.17.0.2:5000")
    if err != nil {
        t.Errorf(err.Error())
    }

    // Check number of bytes received from cn
    n := <- cn
    if n != 102 {
        t.Error("Error: Incorrect number of bytes received")
    }

    // Check received sequence
    if bytes.Compare(wolExpected, wolBuf) != 0 {
        t.Error("Error: Incorrect sequence received")
    }
}
