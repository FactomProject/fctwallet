// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"encoding/json"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
)

func ComposeChainSubmit(name string, data []byte) ([]byte, error) {
	type chainsubmit struct {
		ChainID     string
		ChainCommit json.RawMessage
		EntryReveal json.RawMessage
	}

	e := factom.NewEntry()
	if err := e.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	c := factom.NewChain(e)

	sub := new(chainsubmit)
	we, err := wallet.GetDB().FetchWalletEntryByName([]byte(name))
	if err != nil {
		return nil, err
	}

	pub := new([constants.ADDRESS_LENGTH]byte)
	copy(pub[:], we.(interfaces.IWalletEntry).GetKey(0))
	pri := new([constants.PRIVATE_LENGTH]byte)
	copy(pri[:], we.(interfaces.IWalletEntry).GetPrivKey(0))

	sub.ChainID = c.ChainID

	if j, err := factom.ComposeChainCommit(pub, pri, c); err != nil {
		return nil, err
	} else {
		sub.ChainCommit = j
	}

	if j, err := factom.ComposeEntryReveal(c.FirstEntry); err != nil {
		return nil, err
	} else {
		sub.EntryReveal = j
	}

	return json.Marshal(sub)
}

func ComposeEntrySubmit(name string, data []byte) ([]byte, error) {
	type entrysubmit struct {
		EntryCommit json.RawMessage
		EntryReveal json.RawMessage
	}

	e := factom.NewEntry()
	if err := e.UnmarshalJSON(data); err != nil {
		return nil, err
	}

	sub := new(entrysubmit)
	we, err := wallet.GetDB().FetchWalletEntryByName([]byte(name))
	if err != nil {
		return nil, err
	}
	
	pub := new([constants.ADDRESS_LENGTH]byte)
	copy(pub[:], we.(interfaces.IWalletEntry).GetKey(0))
	pri := new([constants.PRIVATE_LENGTH]byte)
	copy(pri[:], we.(interfaces.IWalletEntry).GetPrivKey(0))

	if j, err := factom.ComposeEntryCommit(pub, pri, e); err != nil {
		return nil, err
	} else {
		sub.EntryCommit = j
	}

	if j, err := factom.ComposeEntryReveal(e); err != nil {
		return nil, err
	} else {
		sub.EntryReveal = j
	}

	return json.Marshal(sub)
}
