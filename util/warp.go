package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/charmbracelet/log"
)

var (
	CommonWarpPorts = []int{
		854, 859, 864, 878, 880, 890, 891, 894, 903,
		908, 928, 934, 939, 942, 943, 945, 946, 955, 968,
		987, 988, 1002, 1010, 1014, 1018, 1070, 1074, 1180, 1387,
		1701, 1843, 2371, 2408, 2506, 3138, 3476, 3581, 3854, 4177,
		4198, 4233, 4500, 5279, 5956, 7103, 7152, 7156, 7281, 7559, 8319, 8742, 8854, 8886,
	}
	ServerIPs = []string{
		"162.159.192.0/24",
		"162.159.193.0/24",
		"162.159.195.0/24",
		"162.159.204.0/24",
		"188.114.96.0/24",
		"188.114.97.0/24",
		"188.114.98.0/24",
		"188.114.99.0/24",
	}
)

type PassTest struct {
	Address    string
	ReturnTime *int64
}

type Speed struct {
	Server  string
	TimeOut int64
}

func generateRandomHex(length int) (*string, error) {
	randomBytes := make([]byte, length)

	// Read random bytes from a cryptographically random source
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Convert random bytes to hexadecimal string
	randomHex := hex.EncodeToString(randomBytes)

	return &randomHex, nil
}

func HandshakePacket() ([]byte, error) {
	// wireGuardHandshake, err := randomPacket(16)
	// 	if err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}
	//warpHandshakePacket, err := hex.DecodeString(fmt.Sprintf("0100000030ec356d08af3939c1b09d3143c2e3773be539e4c7be2e2996e043f1871497be7ed28138b0473350f28647ca3013fe8de10f1ec7e448542c0ef0f0c5b2976455b6bc3f0224d06f14abfbabb7fc8753865f6dad38d7b1c2156c6cea13f57edc39c6627139659075a1c25d49743a86a40517ec45cf8e151bf0796b3f992070839600000000000000000000000000000000%s", wireGuardHandshake))
	start, err := generateRandomHex(33)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	end, err := generateRandomHex(16)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	warpHandshakePacket, err := hex.DecodeString(fmt.Sprintf("%s000000000000000000000000%s", *start, *end))

	if err != nil {
		log.Error("HandshakePacket Init Fail")
		return nil, err
	}
	return warpHandshakePacket, nil
}
