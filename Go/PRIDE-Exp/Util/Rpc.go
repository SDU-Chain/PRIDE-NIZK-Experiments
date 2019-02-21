package Util

type RpcCommitArgument struct {
	CarID uint64
	//两个 string 分别对应 X 和 Y 的数值
	TildeV [2]string
	//两个 string 分别对应 X 和 Y 的数值
	TildeA    [2]string
	Timestamp int64
}

type RpcNullArgument struct {
}

type RpcNullResponse struct {
}

type RpcTimeResponse struct {
	Timestamp int64
}

type RpcCarIDArgument struct {
	CarID uint64
}

type RpcSignArgument struct {
	CarID uint64
	PiV   [2]string
	PiA   [2]string
}

type RpcSignResponse struct {
	Signature string
}
