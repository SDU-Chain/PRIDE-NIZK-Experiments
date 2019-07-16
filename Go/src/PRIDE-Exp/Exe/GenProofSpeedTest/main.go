package main

import (
	"PRIDE-Exp/Config"
	"PRIDE-Exp/Constant"
	"PRIDE-Exp/Util"
	"flag"
	"fmt"
	"geth-timing/crypto/bn256"

	"math/big"
	mathRand "math/rand"
	"time"
)

//提交 Commit 时的时间戳。保持递增。
var Timestamp int64 = 0

//向量 V 的和。数字需要对 G1 的阶取余。
var SigmaV [Constant.HIGH_V + 1]big.Int

//向量 A 的和。数字需要对 G1 的阶取余。
var SigmaA [Constant.HIGH_A + 1]big.Int

var PiV bn256.G1 = Util.NewG1IdenticalElement()
var PiA bn256.G1 = Util.NewG1IdenticalElement()

func updateSigmaV(vectorV [Constant.HIGH_V + 1]big.Int) {
	for i := 0; i <= Constant.HIGH_V; i++ {
		SigmaV[i].Add(&SigmaV[i], &vectorV[i])
		SigmaV[i].Mod(&SigmaV[i], &Constant.G1Order)
	}
}
func updateSigmaA(vectorA [Constant.HIGH_A + 1]big.Int) {
	for i := 0; i <= Constant.HIGH_A; i++ {
		SigmaA[i].Add(&SigmaA[i], &vectorA[i])
		SigmaA[i].Mod(&SigmaA[i], &Constant.G1Order)
	}
}

func genProof() {
	start := time.Now()

	var r int64 = int64(mathRand.Int31())

	vGamma := Util.NewG1IdenticalElement()
	vGamma.ScalarMult(&Config.G[0], big.NewInt(r))
	aGamma := Util.NewG1IdenticalElement()
	aGamma.ScalarMult(&Config.G[0], big.NewInt(r))

	vHash := Util.CalcHash(PiV, vGamma)
	aHash := Util.CalcHash(PiA, aGamma)

	proof := struct {
		vProduct [2]big.Int
		vGamma   [2]big.Int
		vY       [1]big.Int
		vSigma   [Constant.HIGH_V + 1]big.Int
		aProduct [2]big.Int
		aGamma   [2]big.Int
		aY       [1]big.Int
		aSigma   [Constant.HIGH_A + 1]big.Int
	}{
		vProduct: Util.G1ToBigXY(PiV),
		vGamma:   Util.G1ToBigXY(vGamma),
		vY:       [1]big.Int{*new(big.Int).Mod(new(big.Int).Add(big.NewInt(r), new(big.Int).Mul(new(big.Int).Set(&SigmaV[0]), &vHash)), &Constant.G1Order)},
		vSigma:   SigmaV,
		aProduct: Util.G1ToBigXY(PiA),
		aGamma:   Util.G1ToBigXY(aGamma),
		aY:       [1]big.Int{*new(big.Int).Mod(new(big.Int).Add(big.NewInt(r), new(big.Int).Mul(new(big.Int).Set(&SigmaA[0]), &aHash)), &Constant.G1Order)},
		aSigma:   SigmaA,
	}

	end := time.Now()
	duration := end.Sub(start)
	fmt.Println(proof)
	fmt.Println()

	fmt.Println("genproof:", duration.Nanoseconds(), "ns")
}

func commit(v int, a int) {
	//log.Println("[Commit] v=", v, "a=", a)
	Timestamp++

	vectorV := Util.IntToVectorV(v)
	vectorA := Util.IntToVectorA(a)

	updateSigmaV(vectorV)
	updateSigmaA(vectorA)

	//指数上的向量 V
	var tildeV bn256.G1 = Util.VectorVToTildeV(vectorV)
	//指数上的向量 A
	var tildeA bn256.G1 = Util.VectorAToTildeA(vectorA)

	//update Pi_V and Pi_A
	PiV.Add(&PiV, &tildeV)
	PiA.Add(&PiA, &tildeA)

}

func main() {

	loopCount := flag.Int("count", 0, "The count of commitments. (required)")

	flag.Parse()
	if *loopCount == 0 {
		flag.Usage()
		return
	}
	start := time.Now()
	for i := 1; i <= *loopCount; i++ {
		v := mathRand.Intn(Constant.HIGH_V + 1)
		a := mathRand.Intn(Constant.HIGH_A+1) - Constant.HIGH_A/2
		commit(v, a)
	}
	end := time.Now()
	duration := end.Sub(start)

	genProof()
	fmt.Println("commit:", duration.Nanoseconds()/(int64)(*loopCount), "ns")
}
