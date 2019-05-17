/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package fiiicoin

import (
	"fmt"
	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
)

//BlockchainInfo 本地节点区块链信息
type BlockchainInfo struct {
	IsRunning               bool   // P2P网络是否运行
	Connections             uint64 // 当前连接的节点数
	LocalLastBlockHeight    uint64 // 本地节点最新的区块高度
	LocalLastBlockTime      uint64 // 本地节点最新的区块时间
	TempBlockCount          uint64 //本地节点中临时区块的数量
	TempBlockHeights        string //本地节点中临时区块的区块高度集合，高度按升序排序，以逗号连接，比如 "1,2,3"
	RemoteLatestBlockHeight uint64 // 网络上最新的区块高度
	TimeOffset              int64  // 当前节点与网络上最新区块的时间差

}

func NewBlockchainInfo(json *gjson.Result) *BlockchainInfo {
	b := &BlockchainInfo{}
	//解析json
	b.IsRunning = gjson.Get(json.Raw, "isRunning").Bool()
	b.Connections = gjson.Get(json.Raw, "connections").Uint()
	b.LocalLastBlockHeight = gjson.Get(json.Raw, "localLastBlockHeight").Uint()
	b.LocalLastBlockTime = gjson.Get(json.Raw, "localLastBlockTime").Uint()
	b.TempBlockCount = gjson.Get(json.Raw, "tempBlockCount").Uint()
	b.TempBlockHeights = gjson.Get(json.Raw, "TempBlockHeights").String()
	b.RemoteLatestBlockHeight = gjson.Get(json.Raw, "remoteLatestBlockHeight").Uint()
	b.TimeOffset = gjson.Get(json.Raw, "timeOffset").Int()
	return b
}

//Unspent 未花记录
type Unspent struct {

	/*
			{
		        "txid" : "d54994ece1d11b19785c7248868696250ab195605b469632b7bd68130e880c9a",
		        "vout" : 1,
		        "address" : "mgnucj8nYqdrPFh2JfZSB1NmUThUGnmsqe",
		        "account" : "test label",
		        "scriptPubKey" : "76a9140dfc8bafc8419853b34d5e072ad37d1a5159f58488ac",
		        "amount" : 0.00010000,
		        "confirmations" : 6210,
		        "spendable" : true,
		        "solvable" : true
		    }
	*/
	Key           string `storm:"id"`
	TxID          string `json:"txid"`
	Vout          uint64 `json:"vout"`
	Address       string `json:"address"`
	AccountID     string `json:"account" storm:"index"`
	ScriptPubKey  string `json:"scriptPubKey"`
	Amount        string `json:"amount"`
	Confirmations uint64 `json:"confirmations"`
	Spendable     bool   `json:"spendable"`
	Solvable      bool   `json:"solvable"`
	HDAddress     openwallet.Address
}

func NewUnspent(json *gjson.Result) *Unspent {
	obj := &Unspent{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Vout = gjson.Get(json.Raw, "vout").Uint()
	obj.Address = gjson.Get(json.Raw, "address").String()
	obj.AccountID = gjson.Get(json.Raw, "account").String()
	obj.ScriptPubKey = gjson.Get(json.Raw, "scriptPubKey").String()
	obj.Amount = gjson.Get(json.Raw, "amount").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	//obj.Spendable = gjson.Get(json.Raw, "spendable").Bool()
	obj.Spendable = true
	obj.Solvable = gjson.Get(json.Raw, "solvable").Bool()

	return obj
}

type UnspentSort struct {
	Values     []*Unspent
	Comparator func(a, b *Unspent) int
}

func (s UnspentSort) Len() int {
	return len(s.Values)
}
func (s UnspentSort) Swap(i, j int) {
	s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
}
func (s UnspentSort) Less(i, j int) bool {
	return s.Comparator(s.Values[i], s.Values[j]) < 0
}

type Block struct {

	/*

		"hash": "000000000000000127454a8c91e74cf93ad76752cceb7eb3bcff0c398ba84b1f",
		"confirmations": 2,
		"strippedsize": 191875,
		"size": 199561,
		"weight": 775186,
		"height": 1354760,
		"version": 536870912,
		"versionHex": "20000000",
		"merkleroot": "48239e76f8b37d9c8824fef93d42ac3d7c433029c1e9fa23b6416dd0356f3e57",
		"tx": ["c1e12febeb58aefb0b01c04360262138f4ee0faeb207276e79ea3866608ed84f"]
		"time": 1532143012,
		"mediantime": 1532140298,
		"nonce": 3410287696,
		"bits": "19499855",
		"difficulty": 58358570.79038175,
		"chainwork": "00000000000000000000000000000000000000000000006f68c43926cd6c2d1f",
		"previousblockhash": "00000000000000292d142fcc1ddbd9dafd4518310009f152bdca2a66cc589f97",
		"nextblockhash": "0000000000004a50ef5733ab333f718e6ef5c1995e2cfd5a7caa0875f118cd30"

	*/

	Hash              string
	Confirmations     uint64
	Merkleroot        string
	tx                []string
	Previousblockhash string
	Height            uint64 `storm:"id"`
	Version           uint64
	Time              uint64
	Fork              bool
	txDetails         []*Transaction
	isVerbose         bool
}

func NewBlock(json *gjson.Result) *Block {
	obj := &Block{}
	//解析json
	obj.Height = gjson.Get(json.Raw, "height").Uint()
	obj.Hash = gjson.Get(json.Raw, "hash").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Merkleroot = gjson.Get(json.Raw, "merkleroot").String()
	obj.Previousblockhash = gjson.Get(json.Raw, "previousblockhash").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.Time = gjson.Get(json.Raw, "time").Uint()

	txs := make([]string, 0)
	txDetails := make([]*Transaction, 0)
	for _, tx := range gjson.Get(json.Raw, "tx").Array() {
		if tx.IsObject() {
			obj.isVerbose = true
			txObj := newTxByCore(&tx)
			txObj.BlockHeight = obj.Height
			txObj.BlockHash = obj.Hash
			txObj.Blocktime = int64(obj.Time)
			txDetails = append(txDetails, txObj)
		} else {
			obj.isVerbose = false
			txs = append(txs, tx.String())
		}

	}

	obj.tx = txs
	obj.txDetails = txDetails

	return obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader(symbol string) *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	obj.Confirmations = b.Confirmations
	obj.Merkleroot = b.Merkleroot
	obj.Previousblockhash = b.Previousblockhash
	obj.Height = b.Height
	obj.Version = b.Version
	obj.Time = b.Time
	obj.Symbol = symbol

	return &obj
}

//UnscanRecords 扫描失败的区块及交易
type UnscanRecord struct {
	ID          string `storm:"id"` // primary key
	BlockHeight uint64
	TxID        string
	Reason      string
}

func NewUnscanRecord(height uint64, txID, reason string) *UnscanRecord {
	obj := UnscanRecord{}
	obj.BlockHeight = height
	obj.TxID = txID
	obj.Reason = reason
	obj.ID = common.Bytes2Hex(crypto.SHA256([]byte(fmt.Sprintf("%d_%s", height, txID))))
	return &obj
}

type Transaction struct {
	TxID          string
	Size          uint64
	Version       uint64
	LockTime      int64
	Hex           string
	BlockHash     string
	BlockHeight   uint64
	Confirmations uint64
	Blocktime     int64
	IsCoinBase    bool
	Fees          uint64
	Decimals      int32
	Timestamp     int64
	ExpiredTime   int64
	IsDiscarded   bool

	Vins  []*Vin
	Vouts []*Vout
}

type Vin struct {
	Coinbase string
	TxID     string
	Vout     uint64
	N        uint64
	Addr     string
	Amount   uint64
	Size     uint64
}

type Vout struct {
	Vout        uint64
	Addr        string
	Amount      uint64
	LockScript  string
	Spent       bool
	IsDiscarded bool
}

func newTxByCore(json *gjson.Result) *Transaction {

	/*
		{
			"txid": "6595e0d9f21800849360837b85a7933aeec344a89f5c54cf5db97b79c803c462",
			"hash": "f758cb5181d51f8bee1512b4a862faad5b51c7c85a1a11cd6092ffc1c6649bc5",
			"version": 2,
			"size": 249,
			"vsize": 168,
			"locktime": 1414190,
			"vin": [],
			"vout": [],
			"hex": "02000000000101cc8a3077023c08040e677647ad0e528564764f456b01d8519828df165ab3c4550100000017160014aa59f94152351c79b57b14a53e538a923e332468feffffff02a716167c6f00000017a914a0fe07f130a36d9c7581ccd2886895c049b0cc8287ece29c00000000001976a9148c0bceb59d452b3e077f73a420b8bfe09e0550a788ac0247304402205e667171c1798cde426282bb8bff45901866ad6bf0d209e856c1765eda65ba4802203aaa319ea3de00eccef0006e6ee2089aed4b91ada7953f420a47c9c258d424ca0121033cfda2f93d13b01d46ecc406b03ebaba3e1bd526d2148a0a5d579d52f8c7cf022e941500",
			"blockhash": "0000000040730ea7935cce346ce68bf4c07c10b137ba31960bf8a47c4f7da4ec",
			"confirmations": 20076,
			"time": 1537841342,
			"blocktime": 1537841342
		}
	*/

	obj := Transaction{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.LockTime = gjson.Get(json.Raw, "locktime").Int()
	obj.BlockHash = gjson.Get(json.Raw, "blockhash").String()
	//obj.BlockHeight = gjson.Get(json.Raw, "blockheight").Uint()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Blocktime = gjson.Get(json.Raw, "blocktime").Int()
	obj.Size = gjson.Get(json.Raw, "size").Uint()
	//obj.Fees = gjson.Get(json.Raw, "fees").String()
	obj.Decimals = Decimals
	obj.Vins = make([]*Vin, 0)
	if vins := gjson.Get(json.Raw, "vin"); vins.IsArray() {
		for i, vin := range vins.Array() {
			input := newTxVinByCore(&vin)
			input.N = uint64(i)
			obj.Vins = append(obj.Vins, input)
		}
	}

	obj.Vouts = make([]*Vout, 0)
	if vouts := gjson.Get(json.Raw, "vout"); vouts.IsArray() {
		for _, vout := range vouts.Array() {
			output := newTxVoutByCore(&vout)
			obj.Vouts = append(obj.Vouts, output)
		}
	}

	return &obj
}

func newTxVinByCore(json *gjson.Result) *Vin {


	obj := Vin{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Vout = gjson.Get(json.Raw, "vout").Uint()
	obj.Addr = gjson.Get(json.Raw, "AccountId").String()
	obj.Amount = gjson.Get(json.Raw, "Amount").Uint()

	return &obj
}

func newTxVoutByCore(json *gjson.Result) *Vout {

	obj := Vout{}
	//解析json
	obj.Amount = gjson.Get(json.Raw, "value").Uint()
	obj.Vout = gjson.Get(json.Raw, "vout").Uint()
	obj.LockScript = gjson.Get(json.Raw, "scriptPubKey.hex").String()


	return &obj
}
