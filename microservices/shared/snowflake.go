package shared

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	TIMESTAMP_BITS = 40
	WORKER_ID_BITS = 8
	RANDOM_BITS  = 16
)

func getWorkerIP() (ip net.IP, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ip, fmt.Errorf("Error getting network interfaces: %v", err)
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				return v.IP, nil
			case *net.IPAddr:
				return v.IP, nil
			}
		}
	}

	return ip, fmt.Errorf("No valid ip addresses")
}

var src = rand.NewSource(time.Now().UnixNano())
var Epoch = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

type Snowflake int64

func NewSnowflake() (Snowflake, error) {
	ip, err := getWorkerIP()
	if err != nil {
		return 0, fmt.Errorf("Error getting worker ip: %v", err)
	}

	timestamp := int64(time.Since(Epoch).Milliseconds()) << (WORKER_ID_BITS + RANDOM_BITS)
	workerId := int64(ip[len(ip) - 1]) << RANDOM_BITS
	random := src.Int63() & (1 << RANDOM_BITS - 1)

	return Snowflake(timestamp | workerId | random), err
}

func (snowflake Snowflake) IssuedAt() time.Time {
	ms := time.Duration(snowflake >> (WORKER_ID_BITS + RANDOM_BITS))
	return Epoch.Add(ms * time.Millisecond)
}