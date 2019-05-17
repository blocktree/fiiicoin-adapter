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
	"github.com/blocktree/fiiicoin-adapter/fiiicoin_addrdec"
	"testing"
)

func TestAddressDecoder_AddressEncode(t *testing.T) {
	fiiicoin_addrdec.Default.IsTestNet = false

	p2pk, _ := hex.DecodeString("bd0c059c996ce6b4e1e42e71e91059626e4e135a")
	p2pkAddr, _ := fiiicoin_addrdec.Default.AddressEncode(p2pk)
	t.Logf("p2pkAddr: %s", p2pkAddr)

}

func TestAddressDecoder_AddressDecode(t *testing.T) {

	fiiicoin_addrdec.Default.IsTestNet = false

	p2pkAddr := "fiiimU5jzazxf7B9naGSQauE5XwPCZBKajiQe2"
	p2pkHash, _ := fiiicoin_addrdec.Default.AddressDecode(p2pkAddr)
	t.Logf("p2pkHash: %s", hex.EncodeToString(p2pkHash))

}