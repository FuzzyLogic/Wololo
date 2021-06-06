package main

import (
	"log/syslog"
	"net/http"
	"flag"
	"net"
	"fmt"
	"encoding/binary"
	"bytes"
	wololo "github.com/FuzzyLogic/Wololo/internal"
)

var globalConfig wololo.WololoConfig
var globalLog *syslog.Writer
var globalVerbose bool

// HTTP handler function to handle requests.
// This handler will send the WOL packet to the configured destination.
// The response will contain information on whether the packet was transmitted or an error occurred.
func wolHandler(respWr http.ResponseWriter, req *http.Request) {
	wolPacket := wololo.BuildWolPacket(globalConfig)

	// WOL needs to be broadcasted, create UDPAddress from config
	bcastUDPAddr, err := net.ResolveUDPAddr("udp", globalConfig.UdpBcastAddr)
	if err != nil {
		wololo.WriteToLog(globalLog, "Unable to obtain UDP broadcast address object")
		fmt.Fprintf(respWr, "Unable to obtain UDP broadcast address object")
		fmt.Fprintf(respWr, "Error: "+err.Error())
		return
	}

	// Get local IP address from specified interfaces
	localUDPAddr, err := wololo.InterfaceToIp(globalConfig.Iface)
	if err != nil {
		wololo.WriteToLog(globalLog, "Unable to get local IP address from interface "+globalConfig.Iface)
		fmt.Fprintf(respWr, "Unable to get local IP address from interface\n"+globalConfig.Iface)
		fmt.Fprintf(respWr, "Error: "+err.Error()+"\n")
		return
	}

	// Open UDP connection to send WOL packet
	con, err := net.DialUDP("udp", localUDPAddr, bcastUDPAddr)
	if err != nil {
		wololo.WriteToLog(globalLog, "Unable to create UDP connection")
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
		wololo.WriteToLog(globalLog, "Error sending WOL packet")
		fmt.Fprintf(respWr, "Error sending WOL packet\n")
		fmt.Fprintf(respWr, "Error: "+err.Error()+"\n")
	} else if bytesWritten != 102 {
		// Not an error but something went wrong - inform user
		wololo.WriteToLog(globalLog, "Warning: WOL packet transmission may have been incomplete")
		fmt.Fprintf(respWr, "Warning: WOL packet transmission may have been incomplete\n")
	} else {
		// Notify user that WOL packet was sent
		wololo.WriteToLog(globalLog, "WOL packet sent")
		fmt.Fprintf(respWr, "Device is off...\nWOLOLO\nDevice is on!")
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
	wololo.WriteToLog(globalLog, "Starting server")
	http.HandleFunc("/", wolHandler)
	if err := http.ListenAndServe(globalConfig.ListenAddr+":"+globalConfig.ListenPort, nil); err != nil {
		wololo.WriteToLog(globalLog, "Unable to start server")
		panic(err)
	}
}
