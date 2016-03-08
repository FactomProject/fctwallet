// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"io/ioutil"
	
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

	j, err := Wallet.ComposeEntrySubmit(name, data)
	if err != nil {
		fmt.Println(err)
		ctx.WriteHeader(httpBad)
		ctx.Write([]byte(err.Error()))
		return
	}

	ctx.WriteHeader(200)
	ctx.Write(j)
	return
}

func HandleComposeChainSubmit(ctx *web.Context, name string) {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Could not read from http request:", err)
		ctx.WriteHeader(httpBad)
		ctx.Write([]byte(err.Error()))
		return
	}

	j, err := Wallet.ComposeChainSubmit(name, data)
	if err != nil {
		fmt.Println(err)
		ctx.WriteHeader(httpBad)
		ctx.Write([]byte(err.Error()))
		return
	}

	ctx.WriteHeader(200)
	ctx.Write(j)
	return
}
