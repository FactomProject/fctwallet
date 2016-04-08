// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"

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

//Balance

type EntryCreditBalanceResponse struct {
	Balance int64
}

type FactoidBalanceResponse struct {
	Balance int64
}

type ResolveAddressResponse struct {
	FactoidAddress     string
	EntryCreditAddress string
}

//Commit

type CommitRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type CommitResponse struct {
	Success string
}

//Generateaddress

type GenerateAddressResponse struct {
	Address string
}

type VerifyAddressTypeResponse struct {
	Type  string
	Valid bool
}

type GenerateAddressFromPrivateKeyRequest struct {
	Name       string `json:"name"`
	PrivateKey string `json:"privateKey,omitempty"`
	Mnemonic   string `json:"mnemonic,omitempty"`
}

type RequestParams struct {
	Key         string `json:"key,omitempty"`
	Name        string `json:"name,omitempty"`
	Amount      int64  `json:"amount,omitempty"`
	Address     interfaces.IAddress
	Transaction interfaces.ITransaction
}

// &key=<key>&name=<name or address>&amount=<amount>
// If no amount is specified, a zero is returned.
func GetV2Params(params interface{}) (*RequestParams, *primitives.JSONError) {
	req, ok := params.(*RequestParams)
	if ok == false {
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
