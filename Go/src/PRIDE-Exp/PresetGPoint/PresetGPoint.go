package main

import (
	"PRIDE-Exp/Constant"
	cryptoRand "crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/bn256"
	"math/big"
	"strconv"
)

func generateKey() {
	var G [Constant.HIGH_G + 1][2]big.Int

	fmt.Println("---- Go code ----")
	fmt.Println("var xBytes, yBytes []byte")

	for i := 0; i <= Constant.HIGH_G; i++ {
		_, g1, err := bn256.RandomG1(cryptoRand.Reader)
		if err != nil {
			panic(err)
		}

		//前32字节是x，后32字节是y
		g1Bytes := g1.Marshal()

		G[i][0].SetBytes(g1Bytes[0:32])
		G[i][1].SetBytes(g1Bytes[32:64])

		fmt.Print("xBytes, _ = hex.DecodeString(\"")
		fmt.Print(hex.EncodeToString(g1Bytes[0:32]))
		fmt.Println("\")")

		fmt.Print("yBytes, _ = hex.DecodeString(\"")
		fmt.Print(hex.EncodeToString(g1Bytes[32:64]))
		fmt.Println("\")")

		fmt.Println("G[" + strconv.Itoa(i) + "].Unmarshal(append(xBytes[:],yBytes[:]...))")
	}

	fmt.Println("---- Solidity code ----")

	for i := 0; i <= Constant.HIGH_G; i++ {
		fmt.Println("G[" + strconv.Itoa(i) + "][0]=" + G[i][0].String() + ";")
		fmt.Println("G[" + strconv.Itoa(i) + "][1]=" + G[i][1].String() + ";")
	}

}

func main() {
	generateKey()
}
