// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/FactomProject/fctwallet/Wallet/Utility"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

/*********************************************************************************************************/
/********************************Factoid Addresses********************************************************/
/*********************************************************************************************************/

func GenerateAddress(name string) (interfaces.IAddress, error) {
	ok := Utility.IsValidKey(name)
	if !ok {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := wallet.GenerateFctAddress([]byte(name), 1, 1)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func GenerateAddressString(name string) (string, error) {
	addr, err := GenerateAddress(name)
	if err != nil {
		return "", err
	}
	return primitives.ConvertFctAddressToUserStr(addr), nil
}

func GenerateAddressFromPrivateKey(name string, privateKey string) (interfaces.IAddress, error) {
	if Utility.IsValidKey(name) == false {
		return nil, fmt.Errorf("Invalid name or address")
	}
	if len(privateKey) != 64 && len(privateKey) != 128 {
		return nil, fmt.Errorf("Invalid private key length")
	}
	if Utility.IsValidHex(privateKey) == false {
		return nil, fmt.Errorf("Invalid private key format")
	}
	priv, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	addr, err := wallet.GenerateFctAddressFromPrivateKey([]byte(name), priv, 1, 1)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func GenerateAddressStringFromPrivateKey(name string, privateKey string) (string, error) {
	addr, err := GenerateAddressFromPrivateKey(name, privateKey)
	if err != nil {
		return "", err
	}
	return primitives.ConvertFctAddressToUserStr(addr), nil
}

func GenerateAddressFromHumanReadablePrivateKey(name string, privateKey string) (interfaces.IAddress, error) {
	if Utility.IsValidKey(name) == false {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := wallet.GenerateFctAddressFromHumanReadablePrivateKey([]byte(name), privateKey, 1, 1)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func GenerateAddressStringFromHumanReadablePrivateKey(name string, privateKey string) (string, error) {
	addr, err := GenerateAddressFromHumanReadablePrivateKey(name, privateKey)
	if err != nil {
		return "", err
	}
	return primitives.ConvertFctAddressToUserStr(addr), nil
}

func GenerateAddressFromMnemonic(name string, privateKey string) (interfaces.IAddress, error) {
	if Utility.IsValidKey(name) == false {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := wallet.GenerateFctAddressFromMnemonic([]byte(name), privateKey, 1, 1)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func GenerateAddressStringFromMnemonic(name string, privateKey string) (string, error) {
	addr, err := GenerateAddressFromMnemonic(name, privateKey)
	if err != nil {
		return "", err
	}
	return primitives.ConvertFctAddressToUserStr(addr), nil
}

/*********************************************************************************************************/
/*************************************EC Addresses********************************************************/
/*********************************************************************************************************/

func GenerateECAddress(name string) (interfaces.IAddress, error) {
	ok := Utility.IsValidKey(name)
	if !ok {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := wallet.GenerateECAddress([]byte(name))
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func GenerateECAddressString(name string) (string, error) {
	addr, err := GenerateECAddress(name)
	if err != nil {
		return "", err
	}
	return primitives.ConvertECAddressToUserStr(addr), nil
}

func GenerateECAddressFromPrivateKey(name string, privateKey string) (interfaces.IAddress, error) {
	if Utility.IsValidKey(name) == false {
		return nil, fmt.Errorf("Invalid name or address")
	}
	if len(privateKey) != 64 && len(privateKey) != 128 {
		return nil, fmt.Errorf("Invalid private key length")
	}
	if Utility.IsValidHex(privateKey) == false {
		return nil, fmt.Errorf("Invalid private key format")
	}
	priv, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	addr, err := wallet.GenerateECAddressFromPrivateKey([]byte(name), priv)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func GenerateECAddressStringFromPrivateKey(name string, privateKey string) (string, error) {
	addr, err := GenerateECAddressFromPrivateKey(name, privateKey)
	if err != nil {
		return "", err
	}
	return primitives.ConvertECAddressToUserStr(addr), nil
}

func GenerateECAddressFromHumanReadablePrivateKey(name string, privateKey string) (interfaces.IAddress, error) {
	if Utility.IsValidKey(name) == false {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := wallet.GenerateECAddressFromHumanReadablePrivateKey([]byte(name), privateKey)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func GenerateECAddressStringFromHumanReadablePrivateKey(name string, privateKey string) (string, error) {
	addr, err := GenerateECAddressFromHumanReadablePrivateKey(name, privateKey)
	if err != nil {
		return "", err
	}
	return primitives.ConvertECAddressToUserStr(addr), nil
}

/*********************************************************************************************************/
/*********************************************check address type******************************************/
/*********************************************************************************************************/
/*********************************************************************************************************/
func VerifyAddressType(address string) (string, bool) {
	var resp string = "Not a Valid Factoid Address"
	var pass bool = false

	if strings.HasPrefix(address, "FA") {
		if primitives.ValidateFUserStr(address) {
			resp = "Factoid - Public"
			pass = true
		}
	} else if strings.HasPrefix(address, "EC") {
		if primitives.ValidateECUserStr(address) {
			resp = "Entry Credit - Public"
			pass = true
		}
	} else if strings.HasPrefix(address, "Fs") {
		if primitives.ValidateFPrivateUserStr(address) {
			resp = "Factoid - Private"
			pass = true
		}
	} else if strings.HasPrefix(address, "Es") {
		if primitives.ValidateECPrivateUserStr(address) {
			resp = "Entry Credit - Private"
			pass = true
		}
	}

	//  Add Netki resolution here
	//else if (checkNetki) {
	//	if (primitives.ValidateECPrivateUserStr(address)) {
	//		resp = "{\"AddressType\":\"Factoid - Public\", \"TypeCode\":4 ,\"Success\":true}"
	//	}
	//}

	return resp, pass
}
