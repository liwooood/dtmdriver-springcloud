package driver

import (
	"fmt"
	"strconv"
	"strings"
)

func parseIpPort(endpoint string) (string, uint64, error) {
	ipPort := strings.Split(endpoint, ":")
	if len(ipPort) != 2 {
		return "", 0, fmt.Errorf("invalid endpoint: %s", endpoint)
	}
	port, err := strconv.ParseUint(ipPort[1], 10, 16)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port in host %s, port: %s", endpoint, ipPort[1])
	}
	return ipPort[0], port, nil
}
