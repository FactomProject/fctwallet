// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"io/ioutil"

	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/wsapi"
	"github.com/FactomProject/fctwallet/Wallet"
	"github.com/FactomProject/web"
)

func HandleComposeEntrySubmit(ctx *web.Context, name string) {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Could not read from http request:", err)
		ctx.WriteHeader(httpBad)
		ctx.Write([]byte(err.Error()))
		return
	}
	dat := new(ComposeRequest)
	dat.Name = name
	dat.Data = string(data)

	req := primitives.NewJSON2Request("compose-entry-submit", 1, dat)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		fmt.Println(jsonError.Error())
		ctx.WriteHeader(httpBad)
		ctx.Write([]byte(jsonError.Error()))
		return
	}

	ctx.WriteHeader(200)
	ctx.Write([]byte(jsonResp.Result.(*ComposeResponse).Message))
	return
}

func HandleV2ComposeEntrySubmit(params interface{}) (interface{}, *primitives.JSONError) {
	data, ok := params.(*ComposeRequest)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	j, err := Wallet.ComposeEntrySubmit(data.Name, data.Data)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(ComposeResponse)
	resp.Message = j
	return resp, nil
}

func HandleComposeChainSubmit(ctx *web.Context, name string) {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Could not read from http request:", err)
		ctx.WriteHeader(httpBad)
		ctx.Write([]byte(err.Error()))
		return
	}
	dat := new(ComposeRequest)
	dat.Name = name
	dat.Data = string(data)

	req := primitives.NewJSON2Request("compose-chain-submit", 1, dat)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		fmt.Println(jsonError.Error())
		ctx.WriteHeader(httpBad)
		ctx.Write([]byte(jsonError.Error()))
		return
	}

	ctx.WriteHeader(200)
	ctx.Write([]byte(jsonResp.Result.(*ComposeResponse).Message))
	return
}

func HandleV2ComposeChainSubmit(params interface{}) (interface{}, *primitives.JSONError) {
	data, ok := params.(*ComposeRequest)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	j, err := Wallet.ComposeChainSubmit(data.Name, data.Data)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(ComposeResponse)
	resp.Message = j
	return resp, nil
}
