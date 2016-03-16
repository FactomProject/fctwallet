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
	answer, err := HandleV2FactoidGenerateAddress(name)
	if err != nil {
		reportResults(ctx, err.Error(), false)
	}
	reportResults(ctx, answer.(*GenerateAddressResponse).Address, true)
}

func HandleV2FactoidGenerateECAddress(params interface{}) (interface{}, *primitives.JSONError)  {
	name, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	if Utility.IsValidKey(name) == false {
		return nil, NewInvalidNameError()
	}

	adrstr, err := Wallet.GenerateECAddressString(name)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(GenerateAddressResponse)
	resp.Address = adrstr

	return resp, nil
}

func HandleFactoidGenerateECAddress(ctx *web.Context, name string) {
	answer, err := HandleV2FactoidGenerateECAddress(name)
	if err != nil {
		reportResults(ctx, err.Error(), false)
	}
	reportResults(ctx, answer.(*GenerateAddressResponse).Address, true)
}

/*********************************************************************************************************/
/******************************************From Private Key***********************************************/
/*********************************************************************************************************/

func HandleFactoidGenerateAddressFromPrivateKey(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	privateKey := ctx.Params["privateKey"]

	req:=new(GenerateAddressFromPrivateKeyRequest)
	req.Name = name
	req.PrivateKey = privateKey

	answer, err := HandleV2FactoidGenerateAddressFromPrivateKey(req)
	if err != nil {
		reportResults(ctx, err.Error(), false)
	}

	reportResults(ctx, answer.(*GenerateAddressResponse).Address, true)
}

func HandleV2FactoidGenerateAddressFromPrivateKey(params interface{}) (interface{}, *primitives.JSONError)  {
	priv, ok := params.(*GenerateAddressFromPrivateKeyRequest)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	name := priv.Name
	privateKey := priv.PrivateKey
	if Utility.IsValidKey(name) == false {
		return nil, NewInvalidNameError()
	}
	if len(privateKey) != 64 && len(privateKey) != 128 {
		return nil, wsapi.NewInvalidParamsError()
	}
	if Utility.IsValidHex(privateKey) == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	adrstr, err := Wallet.GenerateAddressStringFromPrivateKey(name, privateKey)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(GenerateAddressResponse)
	resp.Address = adrstr

	return resp, nil
}

func HandleFactoidGenerateECAddressFromPrivateKey(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	privateKey := ctx.Params["privateKey"]

	req:=new(GenerateAddressFromPrivateKeyRequest)
	req.Name = name
	req.PrivateKey = privateKey

	answer, err := HandleV2FactoidGenerateECAddressFromPrivateKey(req)
	if err != nil {
		reportResults(ctx, err.Error(), false)
	}

	reportResults(ctx, answer.(*GenerateAddressResponse).Address, true)
}

func HandleV2FactoidGenerateECAddressFromPrivateKey(params interface{}) (interface{}, *primitives.JSONError)  {
	priv, ok := params.(*GenerateAddressFromPrivateKeyRequest)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	name := priv.Name
	privateKey := priv.PrivateKey
	if Utility.IsValidKey(name) == false {
		return nil, NewInvalidNameError()
	}
	if len(privateKey) != 64 && len(privateKey) != 128 {
		return nil, wsapi.NewInvalidParamsError()
	}
	if Utility.IsValidHex(privateKey) == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	adrstr, err := Wallet.GenerateECAddressStringFromPrivateKey(name, privateKey)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(GenerateAddressResponse)
	resp.Address = adrstr

	return resp, nil
}

/*********************************************************************************************************/
/********************************From human readable private key******************************************/
/*********************************************************************************************************/

func HandleFactoidGenerateAddressFromHumanReadablePrivateKey(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	privateKey := ctx.Params["privateKey"]

	req:=new(GenerateAddressFromPrivateKeyRequest)
	req.Name = name
	req.PrivateKey = privateKey

	answer, err := HandleV2FactoidGenerateAddressFromHumanReadablePrivateKey(req)
	if err != nil {
		reportResults(ctx, err.Error(), false)
	}

	reportResults(ctx, answer.(*GenerateAddressResponse).Address, true)
}

func HandleV2FactoidGenerateAddressFromHumanReadablePrivateKey(params interface{}) (interface{}, *primitives.JSONError)  {
	priv, ok := params.(*GenerateAddressFromPrivateKeyRequest)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	name := priv.Name
	privateKey := priv.PrivateKey
	if Utility.IsValidKey(name) == false {
		return nil, NewInvalidNameError()
	}

	adrstr, err := Wallet.GenerateAddressStringFromHumanReadablePrivateKey(name, privateKey)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(GenerateAddressResponse)
	resp.Address = adrstr

	return resp, nil
}

func HandleFactoidGenerateECAddressFromHumanReadablePrivateKey(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	privateKey := ctx.Params["privateKey"]

	req:=new(GenerateAddressFromPrivateKeyRequest)
	req.Name = name
	req.PrivateKey = privateKey

	answer, err := HandleV2FactoidGenerateECAddressFromHumanReadablePrivateKey(req)
	if err != nil {
		reportResults(ctx, err.Error(), false)
	}

	reportResults(ctx, answer.(*GenerateAddressResponse).Address, true)
}

func HandleV2FactoidGenerateECAddressFromHumanReadablePrivateKey(params interface{}) (interface{}, *primitives.JSONError)  {
	priv, ok := params.(*GenerateAddressFromPrivateKeyRequest)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	name := priv.Name
	privateKey := priv.PrivateKey
	if Utility.IsValidKey(name) == false {
		return nil, NewInvalidNameError()
	}

	adrstr, err := Wallet.GenerateECAddressStringFromHumanReadablePrivateKey(name, privateKey)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(GenerateAddressResponse)
	resp.Address = adrstr

	return resp, nil
}

/*********************************************************************************************************/
/*********************************************From mnemonic***********************************************/
/*********************************************************************************************************/

func HandleFactoidGenerateAddressFromMnemonic(ctx *web.Context, params string) {
	name := ctx.Params["name"]
	mnemonic := ctx.Params["mnemonic"]

	req:=new(GenerateAddressFromPrivateKeyRequest)
	req.Name = name
	req.Mnemonic = mnemonic

	answer, err := HandleV2FactoidGenerateAddressFromMnemonic(req)
	if err != nil {
		reportResults(ctx, err.Error(), false)
	}

	reportResults(ctx, answer.(*GenerateAddressResponse).Address, true)
}

func HandleV2FactoidGenerateAddressFromMnemonic(params interface{}) (interface{}, *primitives.JSONError)  {
	priv, ok := params.(*GenerateAddressFromPrivateKeyRequest)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	name := priv.Name
	mnemonic := priv.Mnemonic
	if Utility.IsValidKey(name) == false {
		return nil, NewInvalidNameError()
	}

	adrstr, err := Wallet.GenerateAddressStringFromMnemonic(name, mnemonic)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(GenerateAddressResponse)
	resp.Address = adrstr

	return resp, nil
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
