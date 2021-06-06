package wololo

import (
	"encoding/hex"
	"errors"
	"regexp"
	"encoding/json"
	"io/ioutil"
)

// MAC address byte array
type MACAddress [6]byte

// Configuration information for the application
type WololoConfig struct {
	ListenAddr   	string		`json:"listenAddr"`
	ListenPort   	string		`json:"listenPort"`
	UdpBcastAddr 	string		`json:"udpBcastAddr"`
	Iface        	string		`json:"iface"`
	MacAddrStr   	string		`json:"macAddr"`
	MacAddr			MACAddress	
}

// Convert a MAC address string from the configuration file
// to the MACAddress data structure
func parseMAC(macAddr string) (*MACAddress, error) {
	var retVal MACAddress

	curByteIdx := 0
	curByte := ""
	for i, c := range macAddr {
		if c == ':' {
			// Decode previous byte
			hexBytes, err := hex.DecodeString(curByte)
			if err != nil {
				return nil, err
			}

			retVal[curByteIdx] = hexBytes[0]
			curByte = ""
			curByteIdx += 1
		} else {
			curByte = curByte + string(c)

			// Decode if this is the last element
			if i == len(macAddr)-1 {
				hexBytes, err := hex.DecodeString(curByte)
				if err != nil {
					return nil, err
				}

				retVal[curByteIdx] = hexBytes[0]
			}
		}
	}

	return &retVal, nil
}

// Read the configuration file and extract all relevant information
// to fill a WololoConfig structure
func ReadConfig(path string) (*WololoConfig, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON
	configData := WololoConfig{}
	err = json.Unmarshal(configFile, &configData)
	if err != nil {
		return nil, err
	}

	// Regular expressions for some very basic config data checking
	listenAddrRe := regexp.MustCompile(`[0-9.]+`)
	if ok := listenAddrRe.Match([]byte(configData.ListenAddr)); !ok {
		return nil, errors.New("Invalid listen address")
	}	

	listenPortRe := regexp.MustCompile(`[0-9]+`)
	if ok := listenPortRe.Match([]byte(configData.ListenPort)); !ok {
		return nil, errors.New("Invalid listen port")
	}

	bcastAddrRe := regexp.MustCompile(`[0-9.]+:[0-9]+`)
	if ok := bcastAddrRe.Match([]byte(configData.UdpBcastAddr)); !ok {
		return nil, errors.New("Invalid UDP broadcast address")
	}

	ifaceRe := regexp.MustCompile(`[a-zA-Z0-9]+`)
	if ok := ifaceRe.Match([]byte(configData.Iface)); !ok {
		return nil, errors.New("Invalid interface")
	}

	macRe := regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)
	if ok := macRe.Match([]byte(configData.MacAddrStr)); !ok {
		return nil, errors.New("Invalid MAC address")
	}

    // The MAC address must be set
    nullAddr := MACAddress{0, 0, 0, 0, 0, 0}
	convertedMacAddr, err := parseMAC(configData.MacAddrStr)
	if err != nil {
		return nil, err
	}
    if nullAddr == *convertedMacAddr {
        return nil, errors.New("MAC address cannot be all zeros")
    }
	configData.MacAddr = *convertedMacAddr

	return &configData, nil
}
