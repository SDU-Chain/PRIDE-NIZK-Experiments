package main

import (
	"PRIDE-Exp/Constant"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptoRand "crypto/rand"
	"encoding/hex"
	"fmt"
	"geth-timing/crypto/bn256/google"
	"math/big"
	"strconv"
)

func generateG() {

	fmt.Println("---- Generating G ----")
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

		fmt.Println("_, _ = G[" + strconv.Itoa(i) + "].Unmarshal(append(xBytes[:],yBytes[:]...))")
	}

	fmt.Println("---- Solidity code ----")

	for i := 0; i <= Constant.HIGH_G; i++ {
		fmt.Println("G[" + strconv.Itoa(i) + "][0]=" + G[i][0].String() + ";")
		fmt.Println("G[" + strconv.Itoa(i) + "][1]=" + G[i][1].String() + ";")
	}

	fmt.Println("-----------------------")
}

func generateCloudSig() (err error) {
	fmt.Println("---- Generate Cloud Signature ----")

	fmt.Println("---- Go code ----")
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), cryptoRand.Reader)
	if err != nil {
		return err
	}

	publicKey := privateKey.PublicKey

	fmt.Println("privKey_D := UtilShit.BigFromBase10(\"" + fmt.Sprint(privateKey.D) + "\")")
	fmt.Println("privKey := ecdsa.PrivateKey{D: &privKey_D,}")
	fmt.Println("pubKey_X := UtilShit.BigFromBase10(\"" + fmt.Sprint(publicKey.X) + "\")")
	fmt.Println("pubKey_Y := UtilShit.BigFromBase10(\"" + fmt.Sprint(publicKey.Y) + "\")")
	fmt.Println("pubKey := ecdsa.PublicKey{X: &pubKey_X,Y: &pubKey_Y,}")

	fmt.Println("-----------------------")
	return nil
}

func main() {
	var err error
	generateG()
	err = generateCloudSig()
	if err != nil {
		_ = fmt.Errorf("%s\n", err.Error())
	}
}
