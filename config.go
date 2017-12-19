package wololo

import (
	"bufio"
	"encoding/hex"
	"errors"
	"os"
	"regexp"
)

// MAC address byte array
type MACAddress [6]byte

// Configuration information for the application
type WololoConfig struct {
	ListenIP     string
	ListenPort   string
	UdpBcastAddr string
	Iface        string
	MacAddr      MACAddress
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
    // Default settings here
    // Note that only the MAC address MUST be defined in the configuration file
    // An appropriate check is conducted at the end
	config := WololoConfig{
		ListenIP:     "127.0.0.1",
		ListenPort:   "5000",
		UdpBcastAddr: "255.255.255.255",
		Iface:        "eth0",
		MacAddr:      MACAddress{0, 0, 0, 0, 0, 0},
	}

	// Try to open the configuration file
	conf, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer conf.Close()

	// Read in the individual lines
	var lines []string
	scanner := bufio.NewScanner(conf)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Parse as the regexp for the individual types of lines
	filterListenExp, _ := regexp.Compile(`Listen=(?P<ip>[0-9.]+):(?P<port>[0-9]+)`)
	filterBcastAddrExp, _ := regexp.Compile(`Broadcast=(?P<ip>[0-9.]+:[0-9]+)`)
	filterIfaceExp, _ := regexp.Compile(`Interface=(?P<if>[a-zA-Z0-9]+)`)
	filterMACExp, _ := regexp.Compile(`MAC=(?P<mac>([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2})`)

	for _, line := range lines {
		// Try to match the extraced line
		matchListen := filterListenExp.FindStringSubmatch(line)
		matchBcast := filterBcastAddrExp.FindStringSubmatch(line)
		matchIface := filterIfaceExp.FindStringSubmatch(line)
		matchMAC := filterMACExp.FindStringSubmatch(line)

		// Check if something was matched
		if len(matchListen) != 0 {
			config.ListenIP = matchListen[1]
			config.ListenPort = matchListen[2]
		} else if len(matchMAC) != 0 {
			macAddr, err := parseMAC(matchMAC[1])
			if err != nil {
				return nil, errors.New("Error in configuration file")
			}
			config.MacAddr = *macAddr
		} else if len(matchBcast) > 0 {
			config.UdpBcastAddr = matchBcast[1]
		} else if len(matchIface) > 0 {
			config.Iface = matchIface[1]
		} else {
			return nil, errors.New("Error in configuration file")
		}
	}

    // The MAC address must be set
    nullAddr := MACAddress{0, 0, 0, 0, 0, 0}
    if nullAddr == config.MacAddr {
        return nil, errors.New("No valid MAC address in configuration file")
    }

	return &config, nil
}
