package main

import (
	"log/syslog"
	"net/http"
	"flag"
	"net"
	"fmt"
	"encoding/binary"
	"bytes"
	"time"
	wololo "github.com/FuzzyLogic/Wololo/internal"
)

var globalConfig wololo.WololoConfig
var globalLog *syslog.Writer
var globalVerbose bool
var finishedRequestWait chan bool

// HTTP handler function to handle requests.
// This handler will send the WOL packet to the configured destination.
// The response will contain information on whether the packet was transmitted or an error occurred.
func wolHandler(respWr http.ResponseWriter, req *http.Request) {
	// Get WOL parameters from request or, alternatively, use default values from config
	var macAddr wololo.MACAddress
	paramMACAddr := req.URL.Query().Get("macaddr")
	if paramMACAddr != "" {
		err := wololo.CheckMACAddr(paramMACAddr)
		if err != nil {
			wololo.WriteToLog(globalLog, "Error: " + err.Error())
			fmt.Fprintf(respWr, "Error: " + err.Error() + "\n")
			return
		} 
		convertedMacAddr, err := wololo.ParseMAC(paramMACAddr)
		if err != nil {
			wololo.WriteToLog(globalLog, "Error: " + err.Error())
			fmt.Fprintf(respWr, "Error: " + err.Error() + "\n")
			return
		}
		macAddr = *convertedMacAddr
	} else {
		macAddr = globalConfig.MacAddr
	}

	// Build the packet, consisting of indicator and target device's MAC address
	wolPacket := wololo.BuildWolPacket(macAddr)

	// WOL needs to be broadcasted, create UDPAddress from config or query parameters
	var udpBcastAddr string
	paramUdpBcastAddr := req.URL.Query().Get("udpbcastaddr")
	if paramUdpBcastAddr != "" {
		err := wololo.CheckBcastAddr(paramUdpBcastAddr)
		if err != nil {
			wololo.WriteToLog(globalLog, "Error: " + err.Error())
			fmt.Fprintf(respWr, "Error: " + err.Error() + "\n")
			return
		}
		udpBcastAddr = paramUdpBcastAddr
	} else {
		udpBcastAddr = globalConfig.UdpBcastAddr
	}

	resolvedUdpBcastAddr, err := net.ResolveUDPAddr("udp", udpBcastAddr)
	if err != nil {
		wololo.WriteToLog(globalLog, "Error: " + err.Error())
		fmt.Fprintf(respWr, "Error: " + err.Error() + "\n")
		return
	}

	// Get local IP address from specified interface
	localUDPAddr, err := wololo.InterfaceToIp(globalConfig.Iface)
	if err != nil {
		wololo.WriteToLog(globalLog, "Error: " + err.Error())
		fmt.Fprintf(respWr, "Error: " + err.Error() + "\n")
		return
	}

	// Open UDP connection to send WOL packet and start timer s.t. network flooding
	// through too many requests is mitigated
	<- finishedRequestWait
	con, err := net.DialUDP("udp", localUDPAddr, resolvedUdpBcastAddr)
	go func(finishedRequestWait chan bool) {
		time.Sleep(3 * time.Second)
		finishedRequestWait <- true
	}(finishedRequestWait)
	if err != nil {
		wololo.WriteToLog(globalLog, "Error: " + err.Error())
		fmt.Fprintf(respWr, "Error: " + err.Error() + "\n")
		return
	}
	defer con.Close()

	// Broadcast the WOL packet
	var packetBuf bytes.Buffer
	binary.Write(&packetBuf, binary.BigEndian, wolPacket)
	bytesWritten, err := con.Write(packetBuf.Bytes())
	if err != nil {
		wololo.WriteToLog(globalLog, "Error sending WOL packet")
		fmt.Fprintf(respWr, "Error: " + err.Error() + "\n")
	} else if bytesWritten != 102 {
		// Not an error but something went wrong - inform user
		wololo.WriteToLog(globalLog, "Warning: WOL packet transmission may have been incomplete")
		fmt.Fprintf(respWr, "Warning: WOL packet transmission may have been incomplete\n")
	} else {
		// Notify user that WOL packet was sent
		wololo.WriteToLog(globalLog, "WOL packet sent")
		fmt.Fprintf(respWr, "Device is off...\nWOLOLO\nDevice (%s) is on!", macAddr)
	}
}

func main() {
	// Command line arg parsing
	configPathPtr := flag.String("config", "/etc/wololo/config.json", "Path to Wololo configuration file")
	flag.Parse()

	// Start logging and defer connection close to end
	globalLog = wololo.SetupLog()
	defer func() {
		if err := globalLog.Close(); err != nil {
			panic(err)
		}
	}()

	// Apply seccomp sandbox to application if activated
	err := wololo.Sandbox()
	if err != nil {
		wololo.WriteToLog(globalLog, "Error sandboxing application")
		panic(err)
	}

	// Read configuration into global variable
	globalConfigPtr, err := wololo.ReadConfig(*configPathPtr)
	if err != nil {
		wololo.WriteToLog(globalLog, "Error reading configuration")
		panic(err)
	}
	globalConfig = *globalConfigPtr

	// Start HTTP handler
	finishedRequestWait = make(chan bool, 1)
	finishedRequestWait <- true
	wololo.WriteToLog(globalLog, "Starting server")
	http.HandleFunc("/", wolHandler)
	if err := http.ListenAndServe(globalConfig.ListenAddr+":"+globalConfig.ListenPort, nil); err != nil {
		wololo.WriteToLog(globalLog, "Unable to start server")
		panic(err)
	}
}
