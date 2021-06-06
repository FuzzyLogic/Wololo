package main

import (
    "log"
    "net"
    "bytes"
    "net/http"
)

// Start a UDP server and wait for a message on port 7.
// Check if the received sequence corresponds to the expected
// WOL packet.
func main() {
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
        panic(err)
    }

    // Listen on port 7 for WOL signal
    log.Println("Starting UDP server on port 7")
    servCon, err := net.ListenUDP("udp", servAddr)
    if err != nil {
        panic(err)
    }
    defer servCon.Close()

    // Receive data
    cn := make(chan int)
    wolBuf := make([]byte, 102)
    go func(servCon *net.UDPConn, cn chan int, buf []byte) {
        log.Println("Waiting for data")
        tmpN, _, err := servCon.ReadFromUDP(buf)
        if err != nil {
            panic(err)
        }

        // Output number of read bytes to channel
        cn <- tmpN
    }(servCon, cn, wolBuf)

    // Send an HTTP GET to Wololo service
    log.Println("Triggering Wololo service")
    _, err = http.Get("http://wololo-test_wololo:5000")
    if err != nil {
        panic(err)
    }

    // Check number of bytes received from cn
    n := <- cn
    log.Println("Received data on UDP connection! Checking...")
    if n != 102 {
        log.Println("Fail: Incorrect number of bytes received")
    }

    // Check received sequence
    if bytes.Compare(wolExpected, wolBuf) != 0 {
        log.Println("Fail: Incorrect sequence received")
    } else {
        log.Println("Success!")
    }
}