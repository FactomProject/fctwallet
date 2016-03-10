// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
)

func LookupAddress(adrType string, adr string) (string, error) {
	if Utility.IsValidAddress(adr) && strings.HasPrefix(adr, adrType) {
		baddr := primitives.ConvertUserStrToAddress(adr)
		adr = hex.EncodeToString(baddr)
	} else if Utility.IsValidHexAddress(adr) {
		// the address is good enough.
	} else if Utility.IsValidNickname(adr) {
		we, err := wallet.GetDB().FetchWalletEntryByName([]byte(adr))
		if err != nil {
			return "", err
		}

		if we != nil {
			if we.GetType() == "ec" {
				if strings.ToLower(adrType) == "fa" {
					return "", fmt.Errorf("%s is an entry credit address, not a factoid address.", adr)
				}
			} else if we.GetType() == "fct" {
				if strings.ToLower(adrType) == "ec" {
					return "", fmt.Errorf("%s is a factoid address, not an entry credit address.", adr)
				}
			}

			addr, _ := we.GetAddress()
			adr = hex.EncodeToString(addr.Bytes())
		} else {
			return "", fmt.Errorf("Name %s is undefined.", adr)
		}
	} else {
		return "", fmt.Errorf("Invalid Name.  Check that you have entered the name correctly.")
	}

	return adr, nil
}

func FactoidBalance(adr string) (int64, error) {
	adr, err := LookupAddress("FA", adr)
	if err != nil {
		return 0, err
	}

	str := fmt.Sprintf("http://%s:%d/v1/factoid-balance/%s", ipaddressFD, portNumberFD, adr)
	resp, err := http.Get(str)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	b := new(Response)
	if err := json.Unmarshal(body, b); err != nil {
		return 0, err
	}

	if !b.Success {
		return 0, fmt.Errorf("%s", b.Response)
	}

	v, err := strconv.ParseInt(b.Response, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil

}

func ECBalance(adr string) (int64, error) {
	adr, err := LookupAddress("EC", adr)
	if err != nil {
		return 0, err
	}

	str := fmt.Sprintf("http://%s:%d/v1/entry-credit-balance/%s", ipaddressFD, portNumberFD, adr)
	resp, err := http.Get(str)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	b := new(Response)
	if err := json.Unmarshal(body, b); err != nil {
		return 0, err
	}

	if !b.Success {
		return 0, fmt.Errorf("%s", b.Response)
	}

	v, err := strconv.ParseInt(b.Response, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}
