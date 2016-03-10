// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
)

func CommitChain(name string, data []byte) error {
	type walletcommit struct {
		Message string
	}

	type commit struct {
		CommitChainMsg string
	}

	in := new(walletcommit)
	json.Unmarshal(data, in)
	msg, err := hex.DecodeString(in.Message)
	if err != nil {
		return fmt.Errorf("Could not decode message:", err)
	}

	var we interfaces.IWalletEntry

	if Utility.IsValidAddress(name) && strings.HasPrefix(name, "EC") {
		addr := primitives.ConvertUserStrToAddress(name)
		we, err = wallet.GetDB().FetchWalletEntryByPublicKey(addr)
		if err != nil {
			return err
		}
	} else if Utility.IsValidHexAddress(name) {
		addr, err := hex.DecodeString(name)
		if err == nil {
			we, err = wallet.GetDB().FetchWalletEntryByPublicKey(addr)
			if err != nil {
				return err
			}
		}
	} else {
		we, err = wallet.GetDB().FetchWalletEntryByName([]byte(name))
		if err != nil {
			return err
		}
	}

	if we == nil {
		return fmt.Errorf("Unknown address")
	}

	signed := wallet.SignCommit(we, msg)

	com := new(commit)
	com.CommitChainMsg = hex.EncodeToString(signed)
	j, err := json.Marshal(com)
	if err != nil {
		return fmt.Errorf("Could not create json post:", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/v1/commit-chain", ipaddressFD, portNumberFD),
		"application/json",
		bytes.NewBuffer(j))
	if err != nil {
		return fmt.Errorf("Could not post to server:", err)
	}
	resp.Body.Close()

	return nil
}

func CommitEntry(name string, data []byte) error {
	type walletcommit struct {
		Message string
	}

	type commit struct {
		CommitEntryMsg string
	}

	in := new(walletcommit)
	json.Unmarshal(data, in)
	msg, err := hex.DecodeString(in.Message)
	if err != nil {
		return fmt.Errorf("Could not decode message:", err)
	}

	we, err := wallet.GetDB().FetchWalletEntryByName([]byte(name))
	if err != nil {
		return err
	}
	signed := wallet.SignCommit(we.(interfaces.IWalletEntry), msg)

	com := new(commit)
	com.CommitEntryMsg = hex.EncodeToString(signed)
	j, err := json.Marshal(com)
	if err != nil {
		return fmt.Errorf("Could not create json post:", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/v1/commit-entry/", ipaddressFD, portNumberFD),
		"application/json",
		bytes.NewBuffer(j))
	if err != nil {
		return fmt.Errorf("Could not post to server:", err)
	}
	resp.Body.Close()
	return nil
}
