// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"github.com/FactomProject/web"
	"strconv"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/wsapi"
	"github.com/FactomProject/fctwallet/Wallet"
)

func NewInvalidNameError() *primitives.JSONError {
	return primitives.NewJSONError(-32602, "Invalid params", "Name provided is not valid")
}

type NameRequest struct {
	Name string `json:"name"`
}

type AddressRequest struct {
	Address string `json:"address"`
}

type KeyRequest struct {
	Key string `json:"key"`
}

//Balance

type EntryCreditBalanceResponse struct {
	Balance int64 `json:"balance"`
}

type FactoidBalanceResponse struct {
	Balance int64 `json:"balance"`
}

type ResolveAddressResponse struct {
	FactoidAddress     string `json:"factoidaddress"`
	EntryCreditAddress string `json:"entrycreditaddress"`
}

//Commit

type CommitRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type CommitResponse struct {
	Success string `json:"success"`
}

//Compose

type ComposeRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type ComposeResponse struct {
	Message string `json:"message"`
}

//Generateaddress

type GenerateAddressResponse struct {
	Address string `json:"address"`
}

type VerifyAddressTypeResponse struct {
	Type  string `json:"type"`
	Valid bool   `json:"valid"`
}

//Transaction

type FactoidFeeResponse struct {
	Message  string `json:"message"`
	FeeDelta int64  `json:"feedelta"`
}

type GetFeeResponse struct {
	Fee int64 `json:"fee"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type PropertiesResponse struct {
	ProtocolVersion  string `json:"protocolversion"`
	FactomdVersion   string `json:"factomdversion"`
	FCTWalletVersion string `json:"fctwalletversion"`
	Message          string `json:"message"`
}

type AddressesResponse struct {
	EntryCreditAddresses []Address `json:"entrycreditaddresses"`
	FactoidAddresses     []Address `json:"factoidaddresses"`
}

type Address struct {
	Address        string  `json:"address"`
	Key            string  `json:"key"`
	Balance        float64 `json:"balance"`
	BalanceDecimal int64   `json:"balancedecimal"`
}

/*Requests*/

type GenerateAddressFromPrivateKeyRequest struct {
	Name       string `json:"name"`
	PrivateKey string `json:"privateKey,omitempty"`
	Mnemonic   string `json:"mnemonic,omitempty"`
}

type RequestParams struct {
	Key         string                  `json:"key,omitempty"`
	Name        string                  `json:"name,omitempty"`
	Amount      int64                   `json:"amount,omitempty"`
	Address     interfaces.IAddress     `json:"address"`
	Transaction interfaces.ITransaction `json:"transaction"`
}

// &key=<key>&name=<name or address>&amount=<amount>
// If no amount is specified, a zero is returned.
func V1toV2Params(ctx *web.Context) *RequestParams {
	req := new(RequestParams)
	req.Key = ctx.Params["key"]
	req.Name = ctx.Params["name"]
	StrAmount := ctx.Params["amount"]

	if len(StrAmount) == 0 {
		StrAmount = "0"
	}

	amount, err := strconv.ParseInt(StrAmount, 10, 64)
	if err != nil {
		return nil
	}
	req.Amount = amount

	return req
}

// &key=<key>&name=<name or address>&amount=<amount>
// If no amount is specified, a zero is returned.
func GetV2Params(params interface{}) (*RequestParams, *primitives.JSONError) {
	req := new(RequestParams)
	err := wsapi.MapToObject(params, req)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}

	if len(req.Key) == 0 || len(req.Name) == 0 {
		return nil, wsapi.NewCustomInvalidParamsError(fmt.Sprintln("Missing Parameters: key='", req.Key, "' name='", req.Name, "' amount='", req.Amount, "'"))
	}

	_, valid := ValidateKey(req.Key)
	if !valid {
		return nil, wsapi.NewCustomInvalidParamsError("Invalid key")
	}

	// Get the transaction
	trans, err := Wallet.GetTransaction(req.Key)
	if err != nil {
		return nil, wsapi.NewCustomInternalError("Failure to locate the transaction")
	}
	req.Transaction = trans

	// Get the input/output/ec address.  Which could be a name.  First look and see if it is
	// a name.  If it isn't, then look and see if it is an address.  Someone could
	// do a weird Address as a name and fool the code, but that seems unlikely.
	// Could check for that some how, but there are many ways around such checks.

	if len(req.Name) <= constants.ADDRESS_LENGTH {
		we, err := Wallet.GetWalletEntry([]byte(req.Name))
		if err != nil {
			return nil, wsapi.NewCustomInternalError("Failure to locate the transaction")
		}
		if we != nil {
			address, err := we.GetAddress()
			if err != nil || address == nil {
				return nil, wsapi.NewCustomInternalError("Should not get an error geting a address from a Wallet Entry")
			}
			req.Address = address
			return req, nil
		}
	}
	if (primitives.ValidateFUserStr(req.Name) || primitives.ValidateECUserStr(req.Name)) == false {
		return nil, wsapi.NewCustomInvalidParamsError(fmt.Sprintf("The address specified isn't defined or is invalid: %s", req.Name))
	}
	baddr := primitives.ConvertUserStrToAddress(req.Name)

	address := factoid.NewAddress(baddr)
	req.Address = address

	return req, nil
}
