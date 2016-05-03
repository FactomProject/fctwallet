// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"encoding/json"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"

	"github.com/FactomProject/fctwallet/scwallet"
)

func ComposeChainSubmit(name string, data string) (string, error) {
	type chainsubmit struct {
		ChainID     string
		ChainCommit json.RawMessage
		EntryReveal json.RawMessage
	}

	e := new(factom.Entry)
	if err := e.UnmarshalJSON([]byte(data)); err != nil {
		return "", err
	}
	c := factom.NewChain(e)

	sub := new(chainsubmit)
	we, err := wallet.GetDB().FetchWalletEntryByName([]byte(name))
	if err != nil {
		return "", err
	}

	pub := new([scwallet.ADDRESS_LENGTH]byte)
	copy(pub[:], we.(interfaces.IWalletEntry).GetKey(0))
	pri := new([scwallet.PRIVATE_LENGTH]byte)
	copy(pri[:], we.(interfaces.IWalletEntry).GetPrivKey(0))

	sub.ChainID = c.ChainID

	if j, err := factom.ComposeChainCommit(pub, pri, c); err != nil {
		return "", err
	} else {
		sub.ChainCommit = j
	}

	if j, err := factom.ComposeEntryReveal(c.FirstEntry); err != nil {
		return "", err
	} else {
		sub.EntryReveal = j
	}

	return primitives.EncodeJSONString(sub)
}

func ComposeEntrySubmit(name string, data string) (string, error) {
	type entrysubmit struct {
		EntryCommit json.RawMessage
		EntryReveal json.RawMessage
	}

	e := new(factom.Entry)
	if err := e.UnmarshalJSON([]byte(data)); err != nil {
		return "", err
	}

	sub := new(entrysubmit)
	we, err := wallet.GetDB().FetchWalletEntryByName([]byte(name))
	if err != nil {
		return "", err
	}

	pub := new([scwallet.ADDRESS_LENGTH]byte)
	copy(pub[:], we.(interfaces.IWalletEntry).GetKey(0))
	pri := new([scwallet.PRIVATE_LENGTH]byte)
	copy(pri[:], we.(interfaces.IWalletEntry).GetPrivKey(0))

	if j, err := factom.ComposeEntryCommit(pub, pri, e); err != nil {
		return "", err
	} else {
		sub.EntryCommit = j
	}

	if j, err := factom.ComposeEntryReveal(e); err != nil {
		return "", err
	} else {
		sub.EntryReveal = j
	}

	return primitives.EncodeJSONString(sub)
}
