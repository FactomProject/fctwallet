// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"fmt"
	
	fct "github.com/FactomProject/factoid"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factoid/wallet"
)

func ComposeEntrySubmit(name string, data []byte) ([]byte, error) {
	e := factom.NewEntry()
	e.UnmarshalJSON(data)

	we := factoidState.GetDB().GetRaw([]byte(fct.W_NAME), []byte(name))
	jout := make([]byte, 0)
	switch we.(type) {
	case wallet.IWalletEntry:
		pub := new([fct.ADDRESS_LENGTH]byte)
		copy(pub[:], we.(wallet.IWalletEntry).GetKey(0))
		pri := new([fct.PRIVATE_LENGTH]byte)
		copy(pri[:], we.(wallet.IWalletEntry).GetPrivKey(0))

		j, err := factom.ComposeEntryCommit(pub, pri, e)
		if err != nil {
			return nil, err
		}
		jout = j
	default:
		return nil, fmt.Errorf("Cannot use non Entry Credit Address for Entry Commit")
	}
	
	return jout, nil
	// TODO - add entry reveal as part of json return
}