package UtilShit

import "math/big"

func BigFromBase10(s string) big.Int {
	n, _ := new(big.Int).SetString(s, 10)
	return *n
}
