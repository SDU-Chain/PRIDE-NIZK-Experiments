package Util

import (
	"PRIDE-Exp/Config"
	"PRIDE-Exp/Constant"
	"PRIDE-Exp/UtilShit"
	"bytes"
	"crypto/sha256"
	"fmt"
	"geth-timing/crypto/bn256/google"
	"math/big"
	mathRand "math/rand"
)

func BigFromBase10(s string) big.Int {
	return UtilShit.BigFromBase10(s)
}

func BigToBytes32(big big.Int) (ret [32]byte) {
	xBytesRaw := big.Bytes()

	//方法1
	copy(ret[32-len(xBytesRaw):], xBytesRaw)
	return ret

	//方法2
	//iMax := 31 - len(xBytesRaw) + 1
	//if iMax <= 0 {
	//	iMax = 0
	//}
	//for i := 31; i >= iMax; i-- {
	//	ret[i] = xBytesRaw[len(xBytesRaw)-(31-i)-1]
	//}
	//return ret
}

func BigXYToG1(bigX big.Int, bigY big.Int) (ret bn256.G1, err error) {
	xBytes := BigToBytes32(bigX)
	yBytes := BigToBytes32(bigY)

	_, err = ret.Unmarshal(append(xBytes[:], yBytes[:]...))
	return ret, err
}

func G1ToBigXY(p bn256.G1) (ret [2]big.Int) {
	pBytes := p.Marshal()
	ret[0].SetBytes(pBytes[0:32])
	ret[1].SetBytes(pBytes[32:64])
	return ret
}

func StringXYToG1(str [2]string) (ret bn256.G1, err error) {
	bigX := BigFromBase10(str[0])
	bigY := BigFromBase10(str[1])
	return BigXYToG1(bigX, bigY)
}

func G1ToStringXY(p bn256.G1) (ret [2]string) {
	bigs := G1ToBigXY(p)
	ret[0] = bigs[0].String()
	ret[1] = bigs[1].String()
	return ret
}

//返回 bn256.G1 的单位元
func NewG1IdenticalElement() bn256.G1 {
	//注意，为了规避 bn256 库的一个 bug（大概是 bug，已经提交 PR，等开发者回复），
	//需要让这个点自己加自己一次。反正是零元（单位元），不影响结果。
	ret := bn256.G1{}
	ret.Add(&ret, &ret)
	return ret
}

func G1Equals(a, b bn256.G1) bool {
	//通过比较封送处理的结果(X,Y)来判断两点是否相等
	return bytes.Equal(a.Marshal(), b.Marshal())
}

//按PRIDE-NIZK的步骤计算Hash。须与智能合约保持一致。
func CalcHash(piVorA bn256.G1, gammaV bn256.G1) big.Int {
	//计算g
	sum := big.NewInt(0)

	for i := 0; i <= Constant.HIGH_G; i++ {
		sum.Add(sum, new(big.Int).SetBytes(Config.G[i].Marshal()[0:32]))
	}
	sum.Add(sum, new(big.Int).SetBytes(piVorA.Marshal()[0:32]))

	sumBytes := sum.Bytes()
	newBytes := make([]byte, 32)

	iMax := 31 - len(sumBytes) + 1
	if iMax <= 0 {
		iMax = 0
	}
	for i := 31; i >= iMax; i-- {
		newBytes[i] = sumBytes[len(sumBytes)-(31-i)-1]
	}
	shaBytes := sha256.Sum256(newBytes[0:32])
	return *new(big.Int).SetBytes(shaBytes[:])
}

//返回适用于 PRIDE-NIZK 中向量首项的随机数。
//当前的随机数范围为 0 到 2^31-1，只是为了方便。
func GetRandIntR() int64 {
	return int64(mathRand.Int31())
}

//将速度 v 转换为向量 V。特别地，第 0 项以随机数填充。
func IntToVectorV(v int) (ret [Constant.HIGH_V + 1]big.Int) {
	if v < 0 || v > Constant.HIGH_V {
		panic("0<=v<=" + fmt.Sprint(Constant.HIGH_V))
	}
	r := GetRandIntR()

	ret[0] = *big.NewInt(r)

	for i := 1; i <= v; i++ {
		ret[i] = *big.NewInt(1)
	}

	return ret
}

//将速度 a 转换为向量 A。特别地，第 0 项以随机数填充。
func IntToVectorA(a int) (ret [Constant.HIGH_A + 1]big.Int) {

	halfLen := Constant.HIGH_A / 2

	if a < -halfLen || a > halfLen {
		panic("-" + fmt.Sprint(halfLen) + "<=a<=" + fmt.Sprint(halfLen))
	}

	r := GetRandIntR()

	ret[0] = *big.NewInt(r)

	if a > 0 {
		for i := halfLen + 1; i <= halfLen+a; i++ {
			ret[i] = *big.NewInt(1)
		}
	} else if a < 0 {
		for i := halfLen; i >= halfLen+a+1; i-- {
			ret[i] = *big.NewInt(1)
		}
	}

	return ret
}

//将向量 V 转换为指数上的向量 V。
func VectorVToTildeV(vectorV [Constant.HIGH_V + 1]big.Int) (ret bn256.G1) {
	for i := 0; i <= Constant.HIGH_V; i++ {
		var t bn256.G1
		t.ScalarMult(&Config.G[i], &vectorV[i])
		ret.Add(&ret, &t)
	}
	return ret
}

//将向量 A 转换为指数上的向量 A。
func VectorAToTildeA(vectorA [Constant.HIGH_A + 1]big.Int) (ret bn256.G1) {
	for i := 0; i <= Constant.HIGH_A; i++ {
		var t bn256.G1
		t.ScalarMult(&Config.G[i], &vectorA[i])
		ret.Add(&ret, &t)
	}
	return ret
}
