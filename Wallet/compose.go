// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"encoding/json"
	"fmt"
	
	fct "github.com/FactomProject/factoid"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factoid/wallet"
)

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
	we := factoidState.GetDB().GetRaw([]byte(fct.W_NAME), []byte(name))
	switch we.(type) {
	case wallet.IWalletEntry:
		pub := new([fct.ADDRESS_LENGTH]byte)
		copy(pub[:], we.(wallet.IWalletEntry).GetKey(0))
		pri := new([fct.PRIVATE_LENGTH]byte)
		copy(pri[:], we.(wallet.IWalletEntry).GetPrivKey(0))

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
		
	default:
		return nil, fmt.Errorf("Cannot use non Entry Credit Address for Entry Commit")
	}
	
	return json.Marshal(sub)
}