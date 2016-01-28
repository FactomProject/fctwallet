// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"encoding/hex"
	"fmt"
	"github.com/FactomProject/factoid"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
	"strings"
)

/*********************************************************************************************************/
/********************************Factoid Addresses********************************************************/
/*********************************************************************************************************/

func GenerateAddress(name string) (factoid.IAddress, error) {
	ok := Utility.IsValidKey(name)
	if !ok {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := factoidState.GetWallet().GenerateFctAddress([]byte(name), 1, 1)
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
	return factoid.ConvertFctAddressToUserStr(addr), nil
}

func GenerateAddressFromPrivateKey(name string, privateKey string) (factoid.IAddress, error) {
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
	addr, err := factoidState.GetWallet().GenerateFctAddressFromPrivateKey([]byte(name), priv, 1, 1)
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
	return factoid.ConvertFctAddressToUserStr(addr), nil
}

func GenerateAddressFromHumanReadablePrivateKey(name string, privateKey string) (factoid.IAddress, error) {
	if Utility.IsValidKey(name) == false {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := factoidState.GetWallet().GenerateFctAddressFromHumanReadablePrivateKey([]byte(name), privateKey, 1, 1)
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
	return factoid.ConvertFctAddressToUserStr(addr), nil
}

func GenerateAddressFromMnemonic(name string, privateKey string) (factoid.IAddress, error) {
	if Utility.IsValidKey(name) == false {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := factoidState.GetWallet().GenerateFctAddressFromMnemonic([]byte(name), privateKey, 1, 1)
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
	return factoid.ConvertFctAddressToUserStr(addr), nil
}

/*********************************************************************************************************/
/*************************************EC Addresses********************************************************/
/*********************************************************************************************************/

func GenerateECAddress(name string) (factoid.IAddress, error) {
	ok := Utility.IsValidKey(name)
	if !ok {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := factoidState.GetWallet().GenerateECAddress([]byte(name))
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
	return factoid.ConvertECAddressToUserStr(addr), nil
}

func GenerateECAddressFromPrivateKey(name string, privateKey string) (factoid.IAddress, error) {
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
	addr, err := factoidState.GetWallet().GenerateECAddressFromPrivateKey([]byte(name), priv)
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
	return factoid.ConvertECAddressToUserStr(addr), nil
}

func GenerateECAddressFromHumanReadablePrivateKey(name string, privateKey string) (factoid.IAddress, error) {
	if Utility.IsValidKey(name) == false {
		return nil, fmt.Errorf("Invalid name or address")
	}
	addr, err := factoidState.GetWallet().GenerateECAddressFromHumanReadablePrivateKey([]byte(name), privateKey)
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
	return factoid.ConvertECAddressToUserStr(addr), nil
}


/*********************************************************************************************************/
/*********************************************check address type******************************************/
/*********************************************************************************************************/
/*********************************************************************************************************/
func VerifyAddressType(address string)(string, bool) {
  	var resp string = "Not a Valid Factoid Address"
	var pass bool = false
  
  	if (strings.HasPrefix(address,"FA")) {
		if (factoid.ValidateFUserStr(address)) {
			resp = "Factoid - Public"
			pass = true
		}
	} else if (strings.HasPrefix(address,"EC")) {
		if (factoid.ValidateECUserStr(address)) {
			resp = "Entry Credit - Public"
			pass = true
		}
	} else if (strings.HasPrefix(address,"Fs")) {
		if (factoid.ValidateFPrivateUserStr(address)) {
			resp = "Factoid - Private"
			pass = true
		}
	} else if (strings.HasPrefix(address,"Es")) {
		if (factoid.ValidateECPrivateUserStr(address)) {
			resp = "Entry Credit - Public"
			pass = true
		}
	} 



	//  Add Netki resolution here
	//else if (checkNetki) {
	//	if (factoid.ValidateECPrivateUserStr(address)) {
	//		resp = "{\"AddressType\":\"Factoid - Public\", \"TypeCode\":4 ,\"Success\":true}"
	//	}
	//} 


	return resp,pass
}

