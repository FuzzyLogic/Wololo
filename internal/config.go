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
func ParseMAC(macAddr string) (*MACAddress, error) {
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

func checkListenAddr(listenAddr string) error {
	listenAddrRe := regexp.MustCompile(`[0-9.]+`)
	if ok := listenAddrRe.Match([]byte(listenAddr)); !ok {
		return errors.New("Invalid listen address \"" + listenAddr + "\"")
	}

	return nil
}

func checkListenPort(listenPort string) error {
	listenPortRe := regexp.MustCompile(`[0-9]+`)
	if ok := listenPortRe.Match([]byte(listenPort)); !ok {
		return errors.New("Invalid listen port \"" + listenPort + "\"")
	}

	return nil
}

func CheckBcastAddr(bcastAddr string) error {
	bcastAddrRe := regexp.MustCompile(`[0-9.]+:[0-9]+`)
	if ok := bcastAddrRe.Match([]byte(bcastAddr)); !ok {
		return errors.New("Invalid UDP broadcast address \"" + bcastAddr + "\"")
	}

	return nil
}

func checkIface(iface string) error {
	ifaceRe := regexp.MustCompile(`[a-zA-Z0-9-]+`)
	if ok := ifaceRe.Match([]byte(iface)); !ok {
		return errors.New("Invalid interface \"" + iface + "\"")
	}

	return nil
}

func CheckMACAddr(macAddr string) error {
	macRe := regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)
	if ok := macRe.Match([]byte(macAddr)); !ok {
		return errors.New("Invalid MAC address \"" + macAddr + "\"")
	}

	return nil
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
	err = checkListenAddr(configData.ListenAddr)
	if err != nil {
		return nil, err
	}	

	err = checkListenPort(configData.ListenPort)
	if err != nil {
		return nil, err
	}	

	err = CheckBcastAddr(configData.UdpBcastAddr)
	if err != nil {
		return nil, err
	}	

	err = checkIface(configData.Iface)
	if err != nil {
		return nil, err
	}

	err = CheckMACAddr(configData.MacAddrStr)
	if err != nil {
		return nil, err
	}	

    // Convert the MAC address
	convertedMacAddr, err := ParseMAC(configData.MacAddrStr)
	if err != nil {
		return nil, err
	}
	configData.MacAddr = *convertedMacAddr

	return &configData, nil
}
