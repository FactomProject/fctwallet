// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/wsapi"
	"github.com/FactomProject/fctwallet/Wallet"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
	"github.com/FactomProject/web"
)

func HandleV2FactoidGenerateAddress(params interface{}) (interface{}, *primitives.JSONError) {
	name, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	if Utility.IsValidKey(name) == false {
		return nil, NewInvalidNameError()
	}

	adrstr, err := Wallet.GenerateAddressString(name)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(GenerateAddressResponse)
	resp.Address = adrstr

	return resp, nil
}

func HandleFactoidGenerateAddress(ctx *web.Context, name string) {
	if Utility.IsValidKey(name) == false {
		reportResults(ctx, "Name provided is not valid", false)
		return
	}

	adrstr, err := Wallet.GenerateAddressString(name)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}

	reportResults(ctx, adrstr, true)
}

func HandleFactoidGenerateECAddress(ctx *web.Context, name string) {
	if Utility.IsValidKey(name) == false {
		reportResults(ctx, "Name provided is not valid", false)
		return
	}

	adrstr, err := Wallet.GenerateECAddressString(name)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}

	reportResults(ctx, adrstr, true)
}

/*********************************************************************************************************/
/******************************************From Private Key***********************************************/
/*********************************************************************************************************/

func HandleFactoidGenerateAddressFromPrivateKey(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	privateKey := ctx.Params["privateKey"]
	if Utility.IsValidKey(name) == false {
		reportResults(ctx, "Name provided is not valid", false)
		return
	}
	if len(privateKey) != 64 && len(privateKey) != 128 {
		reportResults(ctx, "Invalid private key length", false)
		return
	}
	if Utility.IsValidHex(privateKey) == false {
		reportResults(ctx, "Invalid private key format", false)
		return
	}

	adrstr, err := Wallet.GenerateAddressStringFromPrivateKey(name, privateKey)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}

	reportResults(ctx, adrstr, true)
}

func HandleFactoidGenerateECAddressFromPrivateKey(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	privateKey := ctx.Params["privateKey"]
	if Utility.IsValidKey(name) == false {
		reportResults(ctx, "Name provided is not valid", false)
		return
	}
	if len(privateKey) != 64 && len(privateKey) != 128 {
		reportResults(ctx, "Invalid private key length", false)
		return
	}
	if Utility.IsValidHex(privateKey) == false {
		reportResults(ctx, "Invalid private key format", false)
		return
	}

	adrstr, err := Wallet.GenerateECAddressStringFromPrivateKey(name, privateKey)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}

	reportResults(ctx, adrstr, true)
}

/*********************************************************************************************************/
/********************************From human readable private key******************************************/
/*********************************************************************************************************/

func HandleFactoidGenerateAddressFromHumanReadablePrivateKey(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	privateKey := ctx.Params["privateKey"]
	if Utility.IsValidKey(name) == false {
		reportResults(ctx, "Name provided is not valid", false)
		return
	}

	adrstr, err := Wallet.GenerateAddressStringFromHumanReadablePrivateKey(name, privateKey)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}

	reportResults(ctx, adrstr, true)
}

func HandleFactoidGenerateECAddressFromHumanReadablePrivateKey(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	privateKey := ctx.Params["privateKey"]
	if Utility.IsValidKey(name) == false {
		reportResults(ctx, "Name provided is not valid", false)
		return
	}

	adrstr, err := Wallet.GenerateECAddressStringFromHumanReadablePrivateKey(name, privateKey)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}

	reportResults(ctx, adrstr, true)
}

/*********************************************************************************************************/
/*********************************************From mnemonic***********************************************/
/*********************************************************************************************************/

func HandleFactoidGenerateAddressFromMnemonic(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	mnemonic := ctx.Params["mnemonic"]
	if Utility.IsValidKey(name) == false {
		reportResults(ctx, "Name provided is not valid", false)
		return
	}

	adrstr, err := Wallet.GenerateAddressStringFromMnemonic(name, mnemonic)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}

	reportResults(ctx, adrstr, true)
}

/*********************************************************************************************************/
/*********************************************Check address type******************************************/
/*********************************************************************************************************/
func HandleVerifyAddressType(ctx *web.Context, params string) {
	address := ctx.Params["address"]

	answer, err := HandleV2VerifyAddressType(address)
	if err != nil {
		reportResults(ctx, err.Error(), false)
	}

	reportResults(ctx, answer.(*VerifyAddressTypeResponse).Type, answer.(*VerifyAddressTypeResponse).Valid)
}

func HandleV2VerifyAddressType(params interface{}) (interface{}, *primitives.JSONError) {
	address, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	resp, pass := Wallet.VerifyAddressType(address)

	answer := new(VerifyAddressTypeResponse)
	answer.Type = resp
	answer.Valid = pass
	return answer, nil
}
