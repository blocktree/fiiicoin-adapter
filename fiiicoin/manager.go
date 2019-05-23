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
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

const (
	maxAddresNum = 10000
)

type WalletManager struct {
	openwallet.AssetsAdapterBase

	WalletClient    *Client                       // 节点客户端
	Config          *WalletConfig                 //钱包管理配置
	Blockscanner    *FIIIBlockScanner              //区块扫描器
	Decoder         *AddressDecoder                //地址编码器
	TxDecoder       openwallet.TransactionDecoder //交易单编码器
	Log             *log.OWLogger                 //日志工具
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol)
	//区块扫描器
	wm.Blockscanner = NewFIIIBlockScanner(&wm)
	wm.Decoder = NewAddressDecoder(&wm)
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	return &wm
}

func (wm *WalletManager) GetAddressesByTag(tag string) ([]string, error) {

	var (
		addresses = make([]string, 0)
	)

	request := []interface{}{
		tag,
	}

	result, err := wm.WalletClient.Call("GetAddressesByTag", request)
	if err != nil {
		return nil, err
	}

	array := result.Array()
	for _, a := range array {
		addresses = append(addresses, a.String())
	}

	return addresses, nil

}

//AddWatchOnlyAddress 导入地址核心钱包
func (wm *WalletManager) AddWatchOnlyAddress(publickey string) error {

	request := []interface{}{
		publickey,
	}

	_, err := wm.WalletClient.Call("AddWatchOnlyAddress", request)

	if err != nil {
		return err
	}

	return nil

}

//ExportAddresses 导入地址核心钱包
func (wm *WalletManager) ExportAddresses() ([]string, error) {

	var (
		addresses = make([]string, 0)
	)

	result, err := wm.WalletClient.Call("ExportAddresses", nil)
	if err != nil {
		return nil, err
	}

	array := gjson.Parse(result.String()).Array()
	for _, a := range array {
		addresses = append(addresses, a.Get("Id").String())
	}

	return addresses, nil

}

// GetAccountByAddress
func (wm *WalletManager) GetAccountByAddress(address string) (uint64, error) {

	request := []interface{}{
		address,
	}

	result, err := wm.WalletClient.Call("GetAccountByAddress", request)
	if err != nil {
		return 0, err
	}

	balance := result.Get("Balance").Uint()

	return balance, nil
}

//GetBlockChainInfo 获取钱包区块链信息
func (wm *WalletManager) GetBlockChainInfo() (*BlockchainInfo, error) {

	result, err := wm.WalletClient.Call("GetBlockChainInfo", nil)
	if err != nil {
		return nil, err
	}

	blockchain := NewBlockchainInfo(result)

	return blockchain, nil

}

//ListUnspent 获取未花记录
func (wm *WalletManager) ListUnspent(min uint64, addresses ...string) ([]*Unspent, error) {

	//:分页限制

	var (
		limit       = 100
		searchAddrs = make([]string, 0)
		max         = len(addresses)
		step        = max / limit
		utxo        = make([]*Unspent, 0)
		pice        []*Unspent
		err         error
	)

	for i := 0; i <= step; i++ {
		begin := i * limit
		end := (i + 1) * limit
		if end > max {
			end = max
		}

		searchAddrs = addresses[begin:end]

		pice, err = wm.getListUnspentByCore(min, searchAddrs...)
		if err != nil {
			return nil, err
		}
		utxo = append(utxo, pice...)
	}
	return utxo, nil
}

//getTransactionByCore 获取交易单
func (wm *WalletManager) getListUnspentByCore(min uint64, addresses ...string) ([]*Unspent, error) {

	var (
		utxos = make([]*Unspent, 0)
	)

	request := []interface{}{
		min,
		9999999,
	}

	if len(addresses) > 0 {
		request = append(request, addresses)
	}

	result, err := wm.WalletClient.Call("ListUnspent", request)
	if err != nil {
		return nil, err
	}

	array := result.Array()
	for _, a := range array {
		utxos = append(utxos, NewUnspent(&a))
	}

	return utxos, nil
}

//EstimateFee 预估手续费
func (wm *WalletManager) EstimateFee(inputs, outputs int64, feeRate decimal.Decimal) (decimal.Decimal, error) {

	var base int64 = 68

	//TransactionHash数据量计算
	//单位：Byte
	//基础数据: 68 （整个个tx只要加一次）
	//output: 101 （每个output都要加一次，如一个就是101，两个就是202）
	//input: 262 （每个input都要加一次，如一个就是262，两个就是524）
	//
	//如果一个tx由一个input与两个output组成，那么它的数据量计算方法是：
	//68 + 262 * 1 + 101 * 2
	//计算公式如下：148 * 输入数额 + 34 * 输出数额 + 10
	trx_bytes := decimal.New(inputs*262+outputs*101+base, 0)
	trx_fee := trx_bytes.Div(decimal.New(1000, 0)).Mul(feeRate)
	trx_fee = trx_fee.Round(wm.Decimal())
	return trx_fee, nil
}

//EstimateFeeRate 预估的没KB手续费率
func (wm *WalletManager) EstimateFeeRate() (decimal.Decimal, error) {

	feeRate := decimal.Zero

	estimatesmartfee, err := wm.WalletClient.Call("EstimateSmartFee", nil)
	if err != nil {
		return decimal.Zero, err
	}

	feeRate, _ = decimal.NewFromString(estimatesmartfee.String())
	feeRate = feeRate.Shift(-wm.Decimal())

	return feeRate, nil
}

func (wm *WalletManager) CreateRawTransaction(senders []Unspent, receivers map[string]uint64, changeAddress string, feeRate uint64) (*gjson.Result, error) {

	input := make([]interface{}, 0)
	for _, in := range senders {
		input = append(input, map[string]interface{}{
			"Txid": in.TxID,
			"Vout": in.Vout,
		})
	}

	output := make([]interface{}, 0)
	for addr, amount := range receivers {
		output = append(output, map[string]interface{}{
			"Address": addr,
			"Amount": amount,
		})
	}

	request := []interface{}{
		input,
		output,
		changeAddress,
		0,
		feeRate,
	}

	trx, err := wm.WalletClient.Call("CreateRawTransaction", request)
	if err != nil {
		return nil, err
	}
	return trx, nil
}


//BroadcastTransaction 广播交易
func (wm *WalletManager) BroadcastTransaction(msg interface{}) error {

	request := []interface{}{
		msg,
	}

	_, err := wm.WalletClient.Call("BroadcastTransaction", request)
	if err != nil {
		return err
	}

	return nil

}