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

func HandleCommitChain(ctx *web.Context, name string) {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Could not read from http request:", err)
		ctx.WriteHeader(httpBad)
		ctx.Write([]byte(err.Error()))
		return
	}

	req := new(CommitRequest)
	req.Name = name
	req.Data = string(data)

	_, jsonError := HandleV2CommitChain(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}
}

func HandleV2CommitChain(params interface{}) (interface{}, *primitives.JSONError) {
	req := new(CommitRequest)
	err := wsapi.MapToObject(params, req)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}

	err = Wallet.CommitChain(req.Name, []byte(req.Data))
	if err != nil {
		fmt.Println(err)
		return nil, wsapi.NewCustomInternalError(err.Error())
	}
	resp := new(CommitResponse)
	resp.Success = "Success"
	return resp, nil
}

func HandleCommitEntry(ctx *web.Context, name string) {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Could not read from http request:", err)
		ctx.WriteHeader(httpBad)
		ctx.Write([]byte(err.Error()))
		return
	}

	req := new(CommitRequest)
	req.Name = name
	req.Data = string(data)

	_, jsonError := HandleV2CommitEntry(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}
}

func HandleV2CommitEntry(params interface{}) (interface{}, *primitives.JSONError) {
	req := new(CommitRequest)
	err := wsapi.MapToObject(params, req)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}

	err = Wallet.CommitEntry(req.Name, []byte(req.Data))
	if err != nil {
		fmt.Println(err)
		return nil, wsapi.NewCustomInternalError(err.Error())
	}
	resp := new(CommitResponse)
	resp.Success = "Success"
	return resp, nil
}
