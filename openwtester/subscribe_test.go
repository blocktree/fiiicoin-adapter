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

package openwtester

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/common/file"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openw"
	"github.com/blocktree/openwallet/openwallet"
	"path/filepath"
	"testing"
)

////////////////////////// 测试单个扫描器 //////////////////////////

type subscriberSingle struct {
	manager *openw.WalletManager
}

//BlockScanNotify 新区块扫描完成通知
func (sub *subscriberSingle) BlockScanNotify(header *openwallet.BlockHeader) error {
	log.Notice("header:", header)
	return nil
}

//BlockTxExtractDataNotify 区块提取结果通知
func (sub *subscriberSingle) BlockExtractDataNotify(sourceKey string, data *openwallet.TxExtractData) error {
	log.Notice("account:", sourceKey)

	for i, input := range data.TxInputs {
		log.Std.Notice("data.TxInputs[%d]: %+v", i, input)
	}

	for i, output := range data.TxOutputs {
		log.Std.Notice("data.TxOutputs[%d]: %+v", i, output)
	}

	log.Std.Notice("data.Transaction: %+v", data.Transaction)

	walletID := "WKFkmvsSFz5mC1cAX3edJC2e6hH6ow3X9E"
	accountID := "HX4tUVg5eETb6SvZeGeAFwk4PQ1CWS6dQeyjj3CqfYyK"

	balance, err := sub.manager.GetAssetsAccountBalance(testApp, walletID, accountID)
	if err != nil {
		log.Error("GetAssetsAccountBalance failed, unexpected error:", err)
		return nil
	}
	log.Notice("account balance:", balance)

	return nil
}


func TestSubscribeAddress(t *testing.T) {

	var (
		endRunning = make(chan bool, 1)
		symbol     = "FIII"
		addrs      = map[string]string{
			"fiiimUwLmiZ5gwyVZvam1eeSbweNz2vaVP6GtB": "sender",
			"fiiimP62d42Hyej8SmtYxzYxx6xnom3fHw1FAe": "receiver",
			"fiiimZUdC8SPVHHD6VW4CuBfUSKFgXftg6nsRu": "receiver",
			"fiiimAvd7FuDgVbozaMsThUTpRXLdSoWEtht3r": "receiver",
			"fiiimW3ni6W9usNVp1M91hcf6kXwxrUnfzhsKY": "receiver",
			"fiiimH3DFDRszk6hXNKDhRcZhAyikekLuC4rpL": "receiver",
			"fiiimK9GqjcvpAhy3AYng3M8maQuCkRpVsKk4r": "receiver",
			"fiiimRbFFq3TbdHNj21MKd8oANBqBpTnxT99n9": "receiver",
		}
	)

	//GetSourceKeyByAddress 获取地址对应的数据源标识
	scanAddressFunc := func(address string) (string, bool) {
		key, ok := addrs[address]
		if !ok {
			return "", false
		}
		return key, true
	}

	assetsMgr, err := openw.GetAssetsAdapter(symbol)
	if err != nil {
		log.Error(symbol, "is not support")
		return
	}

	//读取配置
	absFile := filepath.Join(configFilePath, symbol+".ini")

	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return
	}
	assetsMgr.LoadAssetsConfig(c)

	assetsLogger := assetsMgr.GetAssetsLogger()
	if assetsLogger != nil {
		assetsLogger.SetLogFuncCall(true)
	}

	//log.Debug("already got scanner:", assetsMgr)
	scanner := assetsMgr.GetBlockScanner()
	if scanner.SupportBlockchainDAI() {
		dbFilePath := filepath.Join("data", "db")
		dbFileName := "blockchain.db"
		file.MkdirAll(dbFilePath)
		dai, err := openwallet.NewBlockchainLocal(filepath.Join(dbFilePath, dbFileName), false)
		if err != nil {
			log.Error("NewBlockchainLocal err: %v", err)
			return
		}

		scanner.SetBlockchainDAI(dai)
	}
	//scanner.SetRescanBlockHeight(231829)

	if scanner == nil {
		log.Error(symbol, "is not support block scan")
		return
	}

	scanner.SetBlockScanAddressFunc(scanAddressFunc)

	sub := subscriberSingle{manager:testInitWalletManager()}
	scanner.AddObserver(&sub)

	scanner.Run()

	<-endRunning
}


func TestSubscribeScanBlock(t *testing.T) {

	var (
		symbol     = "FIII"
		addrs      = map[string]string{
			"fiiimUwLmiZ5gwyVZvam1eeSbweNz2vaVP6GtB": "sender",
			"fiiimP62d42Hyej8SmtYxzYxx6xnom3fHw1FAe": "sender",
			"fiiimZUdC8SPVHHD6VW4CuBfUSKFgXftg6nsRu": "receiver",
			"fiiimAvd7FuDgVbozaMsThUTpRXLdSoWEtht3r": "receiver",
			"fiiimW3ni6W9usNVp1M91hcf6kXwxrUnfzhsKY": "receiver",
			"fiiimH3DFDRszk6hXNKDhRcZhAyikekLuC4rpL": "receiver",
			"fiiimK9GqjcvpAhy3AYng3M8maQuCkRpVsKk4r": "receiver",
			"fiiimRbFFq3TbdHNj21MKd8oANBqBpTnxT99n9": "receiver",
		}
	)

	//GetSourceKeyByAddress 获取地址对应的数据源标识
	scanAddressFunc := func(address string) (string, bool) {
		key, ok := addrs[address]
		if !ok {
			return "", false
		}
		return key, true
	}

	assetsMgr, err := openw.GetAssetsAdapter(symbol)
	if err != nil {
		log.Error(symbol, "is not support")
		return
	}

	//读取配置
	absFile := filepath.Join(configFilePath, symbol+".ini")

	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return
	}
	assetsMgr.LoadAssetsConfig(c)

	assetsLogger := assetsMgr.GetAssetsLogger()
	if assetsLogger != nil {
		assetsLogger.SetLogFuncCall(true)
	}

	//log.Debug("already got scanner:", assetsMgr)
	scanner := assetsMgr.GetBlockScanner()
	if scanner == nil {
		log.Error(symbol, "is not support block scan")
		return
	}

	scanner.SetBlockScanAddressFunc(scanAddressFunc)

	sub := subscriberSingle{}
	scanner.AddObserver(&sub)

	scanner.ScanBlock(46044)
}
