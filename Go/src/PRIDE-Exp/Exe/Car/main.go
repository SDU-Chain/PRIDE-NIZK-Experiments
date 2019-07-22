package main

import (
	"PRIDE-Exp/Config"
	"PRIDE-Exp/Constant"
	"PRIDE-Exp/Util"
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"geth-timing/accounts/abi"
	"geth-timing/crypto/bn256/google"
	geth "geth-timing/mobile"
	ethereumRpc "geth-timing/rpc"
	"log"
	"math/big"
	mathRand "math/rand"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
	"time"
)

func main() {
	var err error

	//读取设置
	cloudProviderHost := flag.String("cloud", "localhost:12345", "The cloud server in the form of \"hostname:port\"")
	ethereumHost := flag.String("ethereum", "http://localhost:8545", "The ethereum client's rpc url.")
	contractAddress := flag.String("contract", "", "The address of smart contract PRIDE-NIZK. A hex string begin with 0x. (required)")
	contractAccountIndex := flag.Int("account", 0, "The index of eth.accounts[]. 0 is the first account. (default 0)")
	loopCount := flag.Int("count", 0, "The count of commitments. (required)")

	flag.Parse()

	if *cloudProviderHost == "" ||
		*ethereumHost == "" ||
		*contractAddress == "" ||
		*contractAccountIndex < 0 ||
		*loopCount <= 0 {
		flag.Usage()
		return
	}

	CloudProviderHost = *cloudProviderHost
	EthereumHost = *ethereumHost
	ContractAddress = *contractAddress
	ContractAccountIndex = *contractAccountIndex

	//初始化随机数种子
	mathRand.Seed(time.Now().Unix())

	initializeRandomCarID()

	err = connectToCloud()
	if err != nil {
		log.Panic(err)
	}

	err = connectToEthereum()
	if err != nil {
		log.Panic(err)
	}

	EthereumAccounts, err = ethereumEthAccounts()
	if err != nil {
		log.Panic(err)
	}
	EthereumAccount = EthereumAccounts[ContractAccountIndex]
	log.Println("Use account", EthereumAccount)

	GasLimit, err = ethereumLatestGasLimit()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Gas limit is", GasLimit)

	err = rpcNewSession(CarID)
	if err != nil {
		log.Panic(err)
	}
	{
		transactionHash, err := ethereumNewSession()
		if err != nil {
			log.Panic(err)
		}

		//等 Transaction 写入区块
		succeed := false
		succeed, err = ethereumGetTransactionReceiptLoop(transactionHash)
		if err != nil {
			log.Panic(err)
		}
		if !succeed {
			log.Panic("Transaction confirmed. Execution failed.")
		} else {
			log.Println("Transaction confirmed. Execution succeed.")
		}
	}
	for i := 1; i <= *loopCount; i++ {
		v := mathRand.Intn(Constant.HIGH_V + 1)
		a := mathRand.Intn(Constant.HIGH_A+1) - Constant.HIGH_A/2
		err = commit(v, a)
		if err != nil {
			log.Panic(err)
		}
	}

	signature, err := rpcSign(CarID, PiV, PiA)
	if err != nil {
		log.Panic(err)
	}

	//log.Println("Cloud signature: " + signature)
	{
		transactionHash, err := ethereumProof(signature)
		if err != nil {
			log.Panic(err)
		}
		//等 Transaction 写入区块
		succeed, err := ethereumGetTransactionReceiptLoop(transactionHash)
		if err != nil {
			log.Panic(err)
		}
		if !succeed {
			log.Panic("Transaction confirmed. Execution failed.")
		} else {
			log.Println("Transaction confirmed. Execution succeed.")
		}
	}

}

//CloudProvider 的 JsonRPC (over TCP) 服务器。以地址:端口的形式给出。
//
//例如：
//  127.0.0.1:12345
//  [::1]:12345
var CloudProviderHost string

// 以太坊客户端的 JsonRPC (over HTTP) 服务器。以 URL 的形式给出。
//
//例如：
//  http://127.0.0.1:12345
//  http://[::1]:12345
var EthereumHost string

//已经建立的 RPC 连接。若连接已断开，则需重新连接（暂不实现）。
var RpcClient *rpc.Client

//已经与以太坊客户端建立的 RPC 连接。若连接已断开，则需重新连接（暂不实现）。
var EthereumClient *ethereumRpc.Client

//PRIDE-NIZK 智能合约所在地址。以 0x 开头的十六进制字符串。
var ContractAddress string

//从以太坊客户端的 eth.accounts 中选择第几个作为使用的账户。0 是第一个。
var ContractAccountIndex int

//被选择的账户地址
var EthereumAccount string

//从以太坊客户端中获得的账户地址
var EthereumAccounts []string

func connectToCloud() (err error) {
	RpcClient, err = jsonrpc.Dial("tcp", CloudProviderHost)
	if err == nil {
		log.Println("Connected to server " + CloudProviderHost)
	}
	return err
}

func connectToEthereum() (err error) {
	EthereumClient, err = ethereumRpc.Dial(EthereumHost)
	if err == nil {
		log.Println("Ethereum RPC server was set to " + EthereumHost)
	}
	return err
}

func initializeRandomCarID() {
	CarID = mathRand.Uint64()
}

func ethereumNewSession() (hashHex string, err error) {
	ret, err := ethereumNewSessionCall(EthereumAccount, ContractAddress)
	if err != nil {
		return "", err
	}

	//fmt.Println("RET:", ret.String())
	if ret.Cmp(big.NewInt(1)) == 0 {
		log.Println("Start new session on ethereum...")
	} else {
		return "", errors.New("start new session on ethereum failed")
	}

	return ethereumNewSessionTransaction(EthereumAccount, ContractAddress)
}

func commit(v int, a int) error {
	log.Println("[Commit] v=", v, "a=", a)
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

	return rpcCommit(tildeV, tildeA, Timestamp, CarID)
}

var CarID uint64 = 0

//提交 Commit 时的时间戳。保持递增。
var Timestamp int64 = 0

//向量 V 的和。数字需要对 G1 的阶取余。
var SigmaV [Constant.HIGH_V + 1]big.Int

//向量 A 的和。数字需要对 G1 的阶取余。
var SigmaA [Constant.HIGH_A + 1]big.Int

var PiV bn256.G1 = Util.NewG1IdenticalElement()
var PiA bn256.G1 = Util.NewG1IdenticalElement()

var GasLimit int64

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

/*
由于 JSON-RPC 服务器返回 error 时，客户端这边也算 error，所以暂时不用这个函数。
func rpcCallTwice(serviceMethod string, args interface{}, reply interface{}) error {
	retryTime := 0
	maxRetryTime := 1
	for {
		err := RpcClient.Call(serviceMethod, args, reply)
		if err != nil {
			log.Println(err)
			err2 := connectToCloud()
			if err2 != nil {
				return err2
			} else {
				if retryTime < maxRetryTime {
					retryTime++
				} else {
					return err
				}
			}
		} else {
			return nil
		}
	}
}*/

func rpcCallOnce(serviceMethod string, args interface{}, reply interface{}) error {
	err := RpcClient.Call(serviceMethod, args, reply)
	if err == nil {
		log.Println("[RPC] \"" + serviceMethod + "\" sent.")
	} else {
		log.Println("[RPC] \"" + serviceMethod + "\" sent but failed.")
	}
	return err
}

func rpcCommit(tildeV, tildeA bn256.G1, timestamp int64, carID uint64) error {
	arg := Util.RpcCommitArgument{
		CarID:     carID,
		Timestamp: timestamp,
		TildeV:    Util.G1ToStringXY(tildeV),
		TildeA:    Util.G1ToStringXY(tildeA),
	}
	var res Util.RpcNullResponse
	return rpcCallOnce("Cloud.Commit", arg, &res)

}

func rpcNewSession(carID uint64) error {
	arg := Util.RpcCarIDArgument{CarID: carID}
	var res Util.RpcNullResponse
	return rpcCallOnce("Cloud.NewSession", arg, &res)
}

func rpcSign(carID uint64, piV bn256.G1, piA bn256.G1) (signature string, err error) {
	arg := Util.RpcSignArgument{
		CarID: carID,
		PiV:   Util.G1ToStringXY(piV),
		PiA:   Util.G1ToStringXY(piA),
	}

	var res Util.RpcSignResponse
	err = rpcCallOnce("Cloud.Sign", arg, &res)
	if err == nil {
		return res.Signature, nil
	} else {
		return "", err
	}
}

func ethereumProof(cloudSignature string) (hashHex string, err error) {
	var r int64 = int64(mathRand.Int31())

	vGamma := Util.NewG1IdenticalElement()
	vGamma.ScalarMult(&Config.G[0], big.NewInt(r))
	aGamma := Util.NewG1IdenticalElement()
	aGamma.ScalarMult(&Config.G[0], big.NewInt(r))

	vHash := Util.CalcHash(PiV, vGamma)
	aHash := Util.CalcHash(PiA, aGamma)

	proof := Proof{
		vProduct: Util.G1ToBigXY(PiV),
		vGamma:   Util.G1ToBigXY(vGamma),
		vY:       [1]big.Int{*new(big.Int).Mod(new(big.Int).Add(big.NewInt(r), new(big.Int).Mul(new(big.Int).Set(&SigmaV[0]), &vHash)), &Constant.G1Order)},
		vSigma:   SigmaV,
		aProduct: Util.G1ToBigXY(PiA),
		aGamma:   Util.G1ToBigXY(aGamma),
		aY:       [1]big.Int{*new(big.Int).Mod(new(big.Int).Add(big.NewInt(r), new(big.Int).Mul(new(big.Int).Set(&SigmaA[0]), &aHash)), &Constant.G1Order)},
		aSigma:   SigmaA,
	}

	////test
	//fmt.Println("PROOF:")
	//fmt.Println("vProduct:(", proof.vProduct[0].String(), ",", proof.vProduct[1].String(), ")")
	//fmt.Println("vGamma:(", proof.vGamma[0].String(), ",", proof.vGamma[1].String(), ")")
	//fmt.Println("vY:", proof.vY[0].String())
	//fmt.Println("vSigma[0]:", proof.vSigma[0].String())
	//fmt.Println("vSigma[1]:", proof.vSigma[1].String())
	//fmt.Println("aProduct:(", proof.aProduct[0].String(), ",", proof.aProduct[1].String(), ")")
	//fmt.Println("aGamma:(", proof.aGamma[0].String(), ",", proof.aGamma[1].String(), ")")
	//fmt.Println("aY:", proof.aY[0].String())
	//fmt.Println("aSigma[0]:", proof.aSigma[0].String())
	//fmt.Println("aSigma[1]:", proof.aSigma[1].String())

	ret, err := ethereumProofCall(EthereumAccount, ContractAddress, proof)
	if err != nil {
		return "", err
	}
	//fmt.Println()
	//fmt.Println("RET:", ret.String())
	if ret.Cmp(big.NewInt(1)) == 0 {
		log.Println("Proof verified.")
	} else {
		return "", errors.New("ethereumProof denied")
	}

	return ethereumProofTransaction(EthereumAccount, ContractAddress, proof)
}

type Proof struct {
	vProduct [2]big.Int
	vGamma   [2]big.Int
	vY       [1]big.Int
	vSigma   [Constant.HIGH_V + 1]big.Int
	aProduct [2]big.Int
	aGamma   [2]big.Int
	aY       [1]big.Int
	aSigma   [Constant.HIGH_A + 1]big.Int
}

func getNewSessionArgumentData() (ret []byte, err error) {
	newSessionAbi, err := abi.JSON(strings.NewReader(Config.SmartContractABI))
	if err != nil {
		return nil, err
	}
	ret, err = newSessionAbi.Pack("newSession")
	if err != nil {
		return nil, err
	}
	//fmt.Println("NewSession DATA:" + hex.EncodeToString(ret))
	return ret, nil
}

func getProofArgumentData(proof Proof) (ret []byte, err error) {

	proofAbi, err := abi.JSON(strings.NewReader(Config.SmartContractABI))
	if err != nil {
		return nil, err
	}

	//参考 GitHub\go-ethereum\accounts\abi\abi_test.go 的 TestInputFixedArrayAndVariableInputLength

	input1 := [2]*big.Int{&proof.vProduct[0], &proof.vProduct[1]}
	input2 := [2]*big.Int{&proof.vGamma[0], &proof.vGamma[1]}
	input3 := [1]*big.Int{&proof.vY[0]}
	input4 := [Constant.HIGH_V + 1]*big.Int{}
	for i := 0; i <= Constant.HIGH_V; i++ {
		input4[i] = &proof.vSigma[i]
	}

	input5 := [2]*big.Int{&proof.aProduct[0], &proof.aProduct[1]}
	input6 := [2]*big.Int{&proof.aGamma[0], &proof.aGamma[1]}
	input7 := [1]*big.Int{&proof.aY[0]}
	input8 := [Constant.HIGH_A + 1]*big.Int{}
	for i := 0; i <= Constant.HIGH_A; i++ {
		input8[i] = &proof.aSigma[i]
	}

	ret, err = proofAbi.Pack(
		"closeSession",
		input1, input2, input3, input4,
		input5, input6, input7, input8,
	)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Proof DATA:" + hex.EncodeToString(ret))
	return ret, nil

}

type EthereumParam struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Data  string `json:"data"`
	Value string `json:"value"`
	Gas   string `json:"gas"`
}

func getNewSessionParam(accountAddressHex, contractAddressHex string) (ret EthereumParam, err error) {
	data, err := getNewSessionArgumentData()
	if err != nil {
		return EthereumParam{}, err
	}
	return getEthereumParam(accountAddressHex, contractAddressHex, data)
}

func getProofParam(accountAddressHex, contractAddressHex string, proof Proof) (ret EthereumParam, err error) {
	data, err := getProofArgumentData(proof)
	if err != nil {
		return EthereumParam{}, err
	}
	return getEthereumParam(accountAddressHex, contractAddressHex, data)
}

func getEthereumParam(accountAddressHex, contractAddressHex string, data []byte) (ret EthereumParam, err error) {
	accountAddress, err := geth.NewAddressFromHex(accountAddressHex)
	if err != nil {
		return EthereumParam{}, err
	}
	contractAddress, err := geth.NewAddressFromHex(contractAddressHex)
	if err != nil {
		return EthereumParam{}, err
	}

	message := geth.NewCallMsg()
	message.SetFrom(accountAddress)
	message.SetTo(contractAddress)
	message.SetData(data)
	message.SetValue(geth.NewBigInt(0))
	//Gas 的值不能高于 block gas limit，也不能低于执行合约需要的 gas。注意 gasLimit 是会动态变化的。
	message.SetGas(GasLimit / 2)

	return EthereumParam{
		From:  message.GetFrom().GetHex(),
		To:    message.GetTo().GetHex(),
		Data:  "0x" + hex.EncodeToString(message.GetData()),
		Value: "0x" + message.GetValue().GetString(16),
		Gas:   "0x" + strconv.FormatInt(message.GetGas(), 16),
	}, nil
}

func ethereumNewSessionTransaction(accountAddressHex, contractAddressHex string) (hashHex string, err error) {
	param, err := getNewSessionParam(accountAddressHex, contractAddressHex)
	if err != nil {
		return "", err
	}
	return ethereumSendTransaction(param)
}

func ethereumProofTransaction(accountAddressHex, contractAddressHex string, proof Proof) (hashHex string, err error) {
	param, err := getProofParam(accountAddressHex, contractAddressHex, proof)
	if err != nil {
		return "", err
	}

	return ethereumSendTransaction(param)
}

func ethereumSendTransaction(param EthereumParam) (hashHex string, err error) {
	err = ethereumRpcCallOnce(&hashHex, "eth_sendTransaction", param)
	log.Println("Transaction hash: " + hashHex)
	return hashHex, err
}
func ethereumCall(param EthereumParam) (ret big.Int, err error) {

	var retHex string
	err = ethereumRpcCallOnce(&retHex, "eth_call", param, "latest")
	if err != nil {
		return ret, err
	}
	//fmt.Println("Call result: ", retHex)
	//base 必须设置成 0 而不是 16，否则 0x 可能被忽略
	ret.SetString(retHex, 0)
	//fmt.Println("RECV INT:", ret.String())

	return ret, nil
}

func ethereumProofCall(accountAddressHex, contractAddressHex string, proof Proof) (ret big.Int, err error) {
	param, err := getProofParam(accountAddressHex, contractAddressHex, proof)
	if err != nil {
		return *new(big.Int), err
	}
	return ethereumCall(param)
}

func ethereumNewSessionCall(accountAddressHex, contractAddressHex string) (ret big.Int, err error) {
	param, err := getNewSessionParam(accountAddressHex, contractAddressHex)
	if err != nil {
		return *new(big.Int), err
	}
	return ethereumCall(param)
}

func ethereumRpcCallOnce(result interface{}, method string, args ...interface{}) error {
	//超时时间暂时设置为 5s
	duration, _ := time.ParseDuration("5s")
	ctx, _ := context.WithTimeout(context.Background(), duration)

	err := EthereumClient.CallContext(ctx, &result, method, args...)
	if err == nil {
		log.Println("<RPC> \"" + method + "\" sent.")
	} else {
		log.Println("<RPC> \"" + method + "\" sent but failed.")
	}
	return err
}

func ethereumEthAccounts() (ret []string, err error) {
	err = ethereumRpcCallOnce(&ret, "eth_accounts")
	return ret, err
}

func ethereumLatestGasLimit() (ret int64, err error) {
	var output map[string]interface{}
	err = ethereumRpcCallOnce(&output, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		return 0, err
	}
	//fmt.Println(output["gasLimit"])
	var t big.Int
	//必须是 0，不能是 16
	t.SetString(output["gasLimit"].(string), 0)
	return t.Int64(), nil
}
func ethereumGetTransactionReceiptLoop(hashHex string) (status bool, err error) {
	var output map[string]interface{}
	for output == nil {
		time.Sleep(time.Duration(1) * time.Second)
		err = ethereumRpcCallOnce(&output, "eth_getTransactionReceipt", hashHex)
		if err != nil {
			return false, err
		}
	}
	var t big.Int
	//必须是 0，不能是 16
	t.SetString(output["status"].(string), 0)

	return t.Cmp(big.NewInt(1)) == 0, nil
}
