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
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/log"
	"path/filepath"
	"testing"
)

var (
	tw *WalletManager
)

func init() {

	tw = testNewWalletManager()
}

func testNewWalletManager() *WalletManager {
	wm := NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "conf.ini")
	//log.Debug("absFile:", absFile)
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	//wm.ExplorerClient.Debug = false
	wm.WalletClient.Debug = true
	return wm
}

func TestWalletManager_GetAddressesByTag(t *testing.T) {
	addresses, err := tw.GetAddressesByTag("")
	if err != nil {
		t.Errorf("GetAddressesByTag failed unexpected error: %v\n", err)
		return
	}

	for i, a := range addresses {
		t.Logf("GetAddressesByAccount address[%d] = %s\n", i, a)
	}
}

func TestWalletManager_ExportAddresses(t *testing.T) {
	addresses, err := tw.ExportAddresses()
	if err != nil {
		t.Errorf("ExportAddresses failed unexpected error: %v\n", err)
		return
	}

	for i, a := range addresses {
		t.Logf("ExportAddresses address[%d] = %s\n", i, a)
	}
}

func TestGetBlockChainInfo(t *testing.T) {
	b, err := tw.GetBlockChainInfo()
	if err != nil {
		t.Errorf("GetBlockChainInfo failed unexpected error: %v\n", err)
	} else {
		log.Infof("GetBlockChainInfo info: %+v\n", b)
	}
}

func TestListUnspent(t *testing.T) {
	utxos, err := tw.ListUnspent(1, "fiiimYt7qZekpQKZauBGxv8kGFJGdMyYtzSgdP")
	if err != nil {
		t.Errorf("ListUnspent failed unexpected error: %v\n", err)
		return
	}

	for _, u := range utxos {
		t.Logf("ListUnspent %s: %s = %d\n", u.Address, u.AccountID, u.Amount)
	}
}

func TestEstimateFee(t *testing.T) {
	feeRate, _ := tw.EstimateFeeRate()
	t.Logf("EstimateFee feeRate = %s\n", feeRate.StringFixed(8))
	fees, _ := tw.EstimateFee(10, 2, feeRate)
	t.Logf("EstimateFee fees = %s\n", fees.StringFixed(8))
}

func TestWalletManager_ImportAddress(t *testing.T) {
	addr := "134id8BvKerWe4MGjn2oRKySX4ipw8ZayP"
	err := tw.AddWatchOnlyAddress(addr)
	if err != nil {
		t.Errorf("RestoreWallet failed unexpected error: %v\n", err)
		return
	}
	log.Info("imported success")
}

func TestWalletManager_GetAccountByAddress(t *testing.T) {
	b, err := tw.GetAccountByAddress("fiiimYt7qZekpQKZauBGxv8kGFJGdMyYtzSgdP")
	if err != nil {
		t.Errorf("GetAccountByAddress failed unexpected error: %v\n", err)
		return
	}

	t.Logf("balance = %d\n", b)
}