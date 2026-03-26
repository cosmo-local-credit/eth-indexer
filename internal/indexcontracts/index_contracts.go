package indexcontracts

import "strings"

const (
	TokenContractAddress = "0xe2CEf4000d6003958c891D251328850f84654eb9"
	PoolContractAddress  = "0x01eD8Fe01a2Ca44Cb26D00b1309d7D777471D00C"
)

func IsPoolContractAddress(address string) bool {
	return strings.EqualFold(address, PoolContractAddress)
}

func IsTrackedIndexContractAddress(address string) bool {
	return IsPoolContractAddress(address) || strings.EqualFold(address, TokenContractAddress)
}
