// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"io/ioutil"

	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/wsapi"
	"github.com/FactomProject/web"
)

func HandleV2Get(ctx *web.Context) {
	HandleV2(ctx, false)
}

func HandleV2Post(ctx *web.Context) {
	HandleV2(ctx, true)
}

func HandleV2(ctx *web.Context, post bool) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		wsapi.HandleV2Error(ctx, nil, wsapi.NewInvalidRequestError())
		return
	}

	j, err := primitives.ParseJSON2Request(string(body))
	if err != nil {
		wsapi.HandleV2Error(ctx, nil, wsapi.NewInvalidRequestError())
		return
	}

	var jsonResp *primitives.JSON2Response
	var jsonError *primitives.JSONError
	if post == true {
		jsonResp, jsonError = HandleV2PostRequest(j)
	} else {
		jsonResp, jsonError = HandleV2GetRequest(j)
	}

	if jsonError != nil {
		wsapi.HandleV2Error(ctx, j, jsonError)
		return
	}

	ctx.Write([]byte(jsonResp.String()))
}

func HandleV2PostRequest(j *primitives.JSON2Request) (*primitives.JSON2Response, *primitives.JSONError) {
	params := j.Params
	var resp interface{}
	var jsonError *primitives.JSONError
	switch j.Method {
	/*case "compose-chain-submit":
		resp, jsonError = HandleV2ComposeChainSubmit(params)
		break
	case "compose-entry-submit":
		resp, jsonError = HandleV2ComposeEntrySubmit(params)
		break
	case "factoid-new-transaction":
		resp, jsonError = HandleV2FactoidNewTransaction(params)
		break
	case "factoid-delete-transaction":
		resp, jsonError = HandleV2FactoidDeleteTransaction(params)
		break
	case "factoid-add-fee":
		resp, jsonError = HandleV2FactoidAddFee(params)
		break
	case "factoid-add-input":
		resp, jsonError = HandleV2FactoidAddInput(params)
		break
	case "factoid-add-output":
		resp, jsonError = HandleV2FactoidAddOutput(params)
		break
	case "factoid-add-ecoutput":
		resp, jsonError = HandleV2FactoidAddECOutput(params)
		break
	case "factoid-sign-transaction":
		resp, jsonError = HandleV2FactoidSignTransaction(params)
		break*/
	case "commit-chain":
		resp, jsonError = HandleV2CommitChain(params)
		break
	case "commit-entry":
		resp, jsonError = HandleV2CommitEntry(params)
		break
		/*case "factoid-submit":
			resp, jsonError = HandleV2FactoidSubmit(params)
			break
		case "factoid-get-processed-transactions":
			resp, jsonError = HandleV2GetProcessedTransactions(params)
			break
		case "factoid-get-processed-transactionsj/":
			resp, jsonError = HandleV2GetProcessedTransactionsj(params)
			break*/
	}
	if jsonError != nil {
		return nil, jsonError
	}

	jsonResp := primitives.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = resp

	return jsonResp, nil
}

func HandleV2GetRequest(j *primitives.JSON2Request) (*primitives.JSON2Response, *primitives.JSONError) {
	params := j.Params
	var resp interface{}
	var jsonError *primitives.JSONError

	switch j.Method {
	case "factoid-balance":
		resp, jsonError = HandleV2FactoidBalance(params)
		break
	case "entry-credit-balance":
		resp, jsonError = HandleV2EntryCreditBalance(params)
		break
	case "factoid-generate-address":
		resp, jsonError = HandleV2FactoidGenerateAddress(params)
		break
	case "factoid-generate-ec-address":
		resp, jsonError = HandleV2FactoidGenerateECAddress(params)
		break
	case "factoid-generate-address-from-private-key":
		resp, jsonError = HandleV2FactoidGenerateAddressFromPrivateKey(params)
		break
	case "factoid-generate-ec-address-from-private-key":
		resp, jsonError = HandleV2FactoidGenerateECAddressFromPrivateKey(params)
		break
	case "factoid-generate-address-from-human-readable-private-key":
		resp, jsonError = HandleV2FactoidGenerateAddressFromHumanReadablePrivateKey(params)
		break
	case "factoid-generate-ec-address-from-human-readable-private-key":
		resp, jsonError = HandleV2FactoidGenerateECAddressFromHumanReadablePrivateKey(params)
		break
	case "factoid-generate-address-from-token-sale":
		resp, jsonError = HandleV2FactoidGenerateAddressFromMnemonic(params)
		break
	case "resolve-address":
		resp, jsonError = HandleV2ResolveAddress(params)
		break
	case "verify-address-type":
		resp, jsonError = HandleV2VerifyAddressType(params)
		break
		/*case "factoid-validate":
			resp, jsonError = HandleV2FactoidValidate(params)
			break
		case "factoid-get-fee":
			resp, jsonError = HandleV2GetFee(params)
			break
		case "factoid-validate":
			resp, jsonError = HandleV2FactoidValidate(params)
			break
		case "factoid-get-fee":
			resp, jsonError = HandleV2GetFee(params)
			break
		case "properties":
			resp, jsonError = HandleV2Properties(params)
			break
		case "factoid-get-addresses":
			resp, jsonError = HandleV2GetAddresses(params)
			break
		case "factoid-get-transactions":
			resp, jsonError = HandleV2GetTransactions(params)
			break
		case "factoid-get-transactionsj":
			resp, jsonError = HandleV2GetTransactionsj(params)
			break*/
	}

	if jsonError != nil {
		return nil, jsonError
	}

	jsonResp := primitives.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = resp

	return jsonResp, nil
}
