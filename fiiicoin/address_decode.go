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
	"encoding/hex"
	"fmt"
	"github.com/blocktree/fiiicoin-adapter/fiiicoin_addrdec"
)

type AddressDecoder struct {
	wm *WalletManager //钱包管理者
}

//NewAddressDecoder 地址解析器
func NewAddressDecoder(wm *WalletManager) *AddressDecoder {
	decoder := AddressDecoder{}
	decoder.wm = wm
	return &decoder
}

//PrivateKeyToWIF 私钥转WIF
func (decoder *AddressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {

	cfg := fiiicoin_addrdec.FIII_mainnetPrivateWIFCompressed
	if decoder.wm.Config.IsTestNet {
		cfg = fiiicoin_addrdec.FIII_testnetPrivateWIFCompressed
	}

	wif, _ := fiiicoin_addrdec.Default.AddressEncode(priv, cfg)

	return wif, nil

}

//PublicKeyToAddress 公钥转地址
func (decoder *AddressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {
	fiiicoin_addrdec.Default.IsTestNet = decoder.wm.Config.IsTestNet
	//暂不清楚意义，
	prefix := []byte{0x30,0x2A,0x30,0x05,0x06,0x03,0x2B,0x65,0x70,0x03,0x21,0x00}
	pub = append(prefix, pub...)
	address, err := fiiicoin_addrdec.Default.AddressEncode(pub)
	if err != nil {
		return "", err
	}

	err = decoder.wm.AddWatchOnlyAddress(hex.EncodeToString(pub))
	if err != nil {
		return "", err
	}

	return address, nil

}

//RedeemScriptToAddress 多重签名赎回脚本转地址
func (decoder *AddressDecoder) RedeemScriptToAddress(pubs [][]byte, required uint64, isTestnet bool) (string, error) {

	return "", fmt.Errorf("RedeemScriptToAddress is not supported")

}

//WIFToPrivateKey WIF转私钥
func (decoder *AddressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {

	cfg := fiiicoin_addrdec.FIII_mainnetPrivateWIFCompressed
	if decoder.wm.Config.IsTestNet {
		cfg = fiiicoin_addrdec.FIII_testnetPrivateWIFCompressed
	}

	priv, err := fiiicoin_addrdec.Default.AddressDecode(wif, cfg)
	if err != nil {
		return nil, err
	}

	return priv, err

}

//ScriptPubKeyToBech32Address scriptPubKey转Bech32地址
func (decoder *AddressDecoder) ScriptPubKeyToBech32Address(scriptPubKey []byte) (string, error) {

	return "", fmt.Errorf("ScriptPubKeyToBech32Address is not supported")

}