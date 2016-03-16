// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/FactomProject/factomd/common/primitives"
)

func NewInvalidNameError() *primitives.JSONError {
	return primitives.NewJSONError(-32602, "Invalid params", "Name provided is not valid")
}

type RequestParams struct {
}

//Balance

type EntryCreditBalanceResponse struct {
	Balance int64
}

type FactoidBalanceResponse struct {
	Balance int64
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
	Name string `json:"name"`
	PrivateKey string `json:"privateKey"`
}