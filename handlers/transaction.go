// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	//"time"
	"github.com/FactomProject/web"
	"strings"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/wsapi"
	"github.com/FactomProject/fctwallet/Wallet"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
)

// Returns either an unbounded list of transactions, or the list of
// transactions that involve a given address.
//
func HandleGetProcessedTransactions(ctx *web.Context, parms string) {
	cmd := ctx.Params["cmd"]
	adr := ctx.Params["address"]

	if cmd == "all" {
		list, err := Utility.DumpTransactions(nil)
		if err != nil {
			reportResults(ctx, err.Error(), false)
			return
		}
		reportResults(ctx, string(list), true)
	} else {

		adr, err := Wallet.LookupAddress("FA", adr)
		if err != nil {
			adr, err = Wallet.LookupAddress("EC", adr)
			if err != nil {
				reportResults(ctx, fmt.Sprintf("Could not understand address %s", adr), false)
				return
			}
		}
		badr, err := hex.DecodeString(adr)

		var adrs [][]byte
		adrs = append(adrs, badr)

		list, err := Utility.DumpTransactions(adrs)
		if err != nil {
			reportResults(ctx, err.Error(), false)
			return
		}
		reportResults(ctx, string(list), true)
	}
}

// Returns either an unbounded list of transactions, or the list of
// transactions that involve a given address.
//
// Return in JSON
//
func HandleGetProcessedTransactionsj(ctx *web.Context, parms string) {
	cmd := ctx.Params["cmd"]
	adr := ctx.Params["address"]
	sstart := ctx.Params["start"]
	send := ctx.Params["end"]

	if len(sstart) == 0 {
		sstart = "0"
	}
	if len(send) == 0 {
		send = "0"
	}

	start, err1 := strconv.ParseInt(sstart, 10, 32)
	end, err2 := strconv.ParseInt(send, 10, 32)
	if err1 != nil || err2 != nil {
		start = 0
		end = 1000000000
	}

	if cmd == "all" {
		list, err := Utility.DumpTransactionsJSON(nil, int(start), int(end))
		if err != nil {
			reportResults(ctx, err.Error(), false)
			return
		}
		reportResults(ctx, string(list), true)
	} else {

		adr, err := Wallet.LookupAddress("FA", adr)
		if err != nil {
			adr, err = Wallet.LookupAddress("EC", adr)
			if err != nil {
				reportResults(ctx, fmt.Sprintf("Could not understand address %s", adr), false)
				return
			}
		}
		badr, err := hex.DecodeString(adr)

		var adrs [][]byte
		adrs = append(adrs, badr)

		list, err := Utility.DumpTransactionsJSON(adrs, int(start), int(end))
		if err != nil {
			reportResults(ctx, err.Error(), false)
			return
		}
		reportResults(ctx, string(list), true)
	}
}

// Setup:  seed --
// Setup creates the 10 fountain Factoid Addresses, then sets address
// generation to be unique for this wallet.  You CAN call setup multiple
// times, but once the Fountain addresses are created, Setup only changes
// the seed.
//
// Setup must be called once before you do anything else with the wallet.
//

/*
func HandleFactoidSetup(ctx *web.Context, seed string) {
	// Make sure we have a seed.
	if len(seed) == 0 {
		msg := "You must supply some random seed. For example (don't use this!)\n" +
			"factom-cli setup 'woe!#in31!%234ng)%^&$%oeg%^&*^jp45694a;gmr@#t4 q34y'\n" +
			"would make a nice seed.  The more random the better.\n\n" +
			"Note that if you create an address before you call Setup, you must\n" +
			"use those address(s) as you access the fountians."

		reportResults(ctx, msg, false)
	}
	setFountian := false
	keys, _ := Wallet.GetWalletNames()
	if len(keys) == 0 {
		setFountian = true
		for i := 1; i <= 10; i++ {
			name := fmt.Sprintf("%02d-Fountain", i)
			_, err := Wallet.GenerateFctAddress([]byte(name), 1, 1)
			if err != nil {
				reportResults(ctx, err.Error(), false)
				return
			}
		}
	}

	seedprime := fct.Sha([]byte(fmt.Sprintf("%s%v", seed, time.Now().UnixNano()))).Bytes()
	Wallet.NewSeed(seedprime)

	if setFountian {
		reportResults(ctx, "New seed set, fountain addresses defined", true)
	} else {
		reportResults(ctx, "New seed set, no fountain addresses defined", true)
	}
}
*/

// New Transaction:  key --
// We create a new transaction, and track it with the user supplied key.  The
// user can then use this key to make subsequent calls to add inputs, outputs,
// and to sign. Then they can submit the transaction.
//
// When the transaction is submitted, we clear it from our working memory.
// Multiple transactions can be under construction at one time, but they need
// their own keys. Once a transaction is either submitted or deleted, the key
// can be reused.
func HandleFactoidNewTransaction(ctx *web.Context, key string) {
	req := primitives.NewJSON2Request("factoid-new-transaction", 1, key)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*MessageResponse).Message, true)
}

func HandleV2FactoidNewTransaction(params interface{}) (interface{}, *primitives.JSONError) {
	key, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	msg, valid := ValidateKey(key)
	if !valid {
		return nil, wsapi.NewCustomInvalidParamsError(msg)
	}

	err := Wallet.FactoidNewTransaction(key)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(MessageResponse)
	resp.Message = "Success building a transaction"
	return resp, nil
}

// Delete Transaction:  key --
// Remove a transaction rather than sign and submit the transaction.  Sometimes
// you just need to throw one a way, and rebuild it.
//
func HandleFactoidDeleteTransaction(ctx *web.Context, key string) {
	req := primitives.NewJSON2Request("factoid-delete-transaction", 1, key)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*MessageResponse).Message, true)
}

func HandleV2FactoidDeleteTransaction(params interface{}) (interface{}, *primitives.JSONError) {
	key, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	// Make sure we have a key
	if len(key) == 0 {
		return nil, wsapi.NewInvalidParamsError()
	}
	err := Wallet.FactoidDeleteTransaction(key)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}
	resp := new(MessageResponse)
	resp.Message = "Success deleting transaction"
	return resp, nil
}

func HandleProperties(ctx *web.Context) {
	req := primitives.NewJSON2Request("properties", 1, nil)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*PropertiesResponse).Message, true)
}

func HandleV2Properties(params interface{}) (interface{}, *primitives.JSONError) {
	p, f, w, err := Wallet.GetProperties()
	if err != nil {
		return nil, wsapi.NewCustomInternalError("Failed to retrieve properties")
	}

	ret := fmt.Sprintf("Protocol Version:   %s\n", p)
	ret = ret + fmt.Sprintf("factomd Version:    %s\n", f)
	ret = ret + fmt.Sprintf("fctwallet Version:  %s\n", w)

	resp := new(PropertiesResponse)
	resp.Message = ret
	resp.ProtocolVersion = p
	resp.FactomdVersion = f
	resp.FCTWalletVersion = w

	return resp, nil
}

func HandleFactoidAddFee(ctx *web.Context, params string) {
	par := V1toV2Params(ctx)
	if par == nil {
		fmt.Println("Not OK")
		return
	}
	req := primitives.NewJSON2Request("factoid-add-fee", 1, par)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*FactoidFeeResponse).Message, true)
}

func HandleV2FactoidAddFee(params interface{}) (interface{}, *primitives.JSONError) {
	pars, jerror := GetV2Params(params)
	if jerror != nil {
		return nil, jerror
	}

	ins, err := pars.Transaction.TotalInputs()
	if err != nil {
		fmt.Println(err.Error())
		return nil, wsapi.NewCustomInternalError(err.Error())
	}
	outs, err := pars.Transaction.TotalOutputs()
	if err != nil {
		fmt.Println(err.Error())
		return nil, wsapi.NewCustomInternalError(err.Error())
	}
	ecs, err := pars.Transaction.TotalECs()
	if err != nil {
		fmt.Println(err.Error())
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	if ins != outs+ecs {
		msg := fmt.Sprintf(
			"Addfee requires that all the inputs balance the outputs.\n"+
				"The total inputs of your transaction are              %s\n"+
				"The total outputs + ecoutputs of your transaction are %s",
			primitives.ConvertDecimalToPaddedString(ins), primitives.ConvertDecimalToPaddedString(outs+ecs))

		fmt.Println(msg)
		return nil, wsapi.NewCustomInternalError("Inputs do not balance the outputs")
	}

	transfee, err := Wallet.FactoidAddFee(pars.Transaction, pars.Key, pars.Address, pars.Name)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(FactoidFeeResponse)
	resp.Message = fmt.Sprintf("Added %s to %s", primitives.ConvertDecimalToPaddedString(uint64(transfee)), pars.Name)
	resp.FeeDelta = int64(transfee)
	return resp, nil
}

func HandleFactoidSubFee(ctx *web.Context, params string) {
	par := V1toV2Params(ctx)
	if par == nil {
		fmt.Println("Not OK")
		return
	}
	req := primitives.NewJSON2Request("factoid-sub-fee", 1, par)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*FactoidFeeResponse).Message, true)
}

func HandleV2FactoidSubFee(params interface{}) (interface{}, *primitives.JSONError) {
	pars, jerror := GetV2Params(params)
	if jerror != nil {
		return nil, jerror
	}

	ins, err := pars.Transaction.TotalInputs()
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}
	outs, err := pars.Transaction.TotalOutputs()
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}
	ecs, err := pars.Transaction.TotalECs()
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	if ins != outs+ecs {
		msg := fmt.Sprintf(
			"Subfee requires that all the inputs balance the outputs.\n"+
				"The total inputs of your transaction are              %s\n"+
				"The total outputs + ecoutputs of your transaction are %s",
			primitives.ConvertDecimalToString(ins), primitives.ConvertDecimalToString(outs+ecs))

		fmt.Println(msg)
		return nil, wsapi.NewCustomInternalError("Inputs do not balance the outputs")
	}

	transfee, err := Wallet.FactoidSubFee(pars.Transaction, pars.Key, pars.Address, pars.Name)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(FactoidFeeResponse)
	resp.Message = fmt.Sprintf("Subtracted %s from %s", primitives.ConvertDecimalToPaddedString(uint64(transfee)), pars.Name)
	resp.FeeDelta = -int64(transfee)
	return resp, nil
}

func HandleFactoidAddInput(ctx *web.Context, parms string) {
	par := V1toV2Params(ctx)
	if par == nil {
		fmt.Println("Not OK")
		return
	}
	req := primitives.NewJSON2Request("factoid-add-input", 1, par)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*MessageResponse).Message, true)
}

func HandleV2FactoidAddInput(params interface{}) (interface{}, *primitives.JSONError) {
	pars, jerror := GetV2Params(params)
	if jerror != nil {
		return nil, jerror
	}

	err := Wallet.FactoidAddInput(pars.Transaction, pars.Key, pars.Address, uint64(pars.Amount))
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(MessageResponse)
	resp.Message = "Success adding Input"
	return resp, nil
}

func HandleFactoidAddOutput(ctx *web.Context, parms string) {
	par := V1toV2Params(ctx)
	if par == nil {
		fmt.Println("Not OK")
		return
	}
	req := primitives.NewJSON2Request("factoid-add-output", 1, par)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*MessageResponse).Message, true)
}

func HandleV2FactoidAddOutput(params interface{}) (interface{}, *primitives.JSONError) {
	pars, jerror := GetV2Params(params)
	if jerror != nil {
		return nil, jerror
	}

	err := Wallet.FactoidAddOutput(pars.Transaction, pars.Key, pars.Address, uint64(pars.Amount))
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(MessageResponse)
	resp.Message = "Success adding Output"
	return resp, nil
}

func HandleFactoidAddECOutput(ctx *web.Context, parms string) {
	par := V1toV2Params(ctx)
	if par == nil {
		fmt.Println("Not OK")
		return
	}
	req := primitives.NewJSON2Request("factoid-add-ecoutput", 1, par)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*MessageResponse).Message, true)
}

func HandleV2FactoidAddECOutput(params interface{}) (interface{}, *primitives.JSONError) {
	pars, jerror := GetV2Params(params)
	if jerror != nil {
		return nil, jerror
	}

	err := Wallet.FactoidAddECOutput(pars.Transaction, pars.Key, pars.Address, uint64(pars.Amount))
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(MessageResponse)
	resp.Message = "Success adding Entry Credit Output"
	return resp, nil
}

func HandleFactoidSignTransaction(ctx *web.Context, key string) {
	req := primitives.NewJSON2Request("factoid-sign-transaction", 1, key)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*MessageResponse).Message, true)
}

func HandleV2FactoidSignTransaction(params interface{}) (interface{}, *primitives.JSONError) {
	key, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	err := Wallet.FactoidSignTransaction(key)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(MessageResponse)
	resp.Message = "Success signing transaction"
	return resp, nil
}

func HandleFactoidSubmit(ctx *web.Context, jsonkey string) {
	req := primitives.NewJSON2Request("factoid-submit", 1, jsonkey)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, jsonResp.Result.(*MessageResponse).Message, true)
}

func HandleV2FactoidSubmit(params interface{}) (interface{}, *primitives.JSONError) {
	jsonkey, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	_, err := Wallet.FactoidSubmit(jsonkey)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(MessageResponse)
	resp.Message = "Success Submitting transaction"
	return resp, nil
}

func HandleGetFee(ctx *web.Context, k string) {
	key := ctx.Params["key"]
	req := primitives.NewJSON2Request("factoid-get-fee", 1, key)

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	reportResults(ctx, fmt.Sprintf("%s", primitives.ConvertDecimalToString(uint64(jsonResp.Result.(*GetFeeResponse).Fee))), true)
}

func HandleV2GetFee(params interface{}) (interface{}, *primitives.JSONError) {
	key, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	var trans interfaces.ITransaction
	var err error

	fmt.Println("getfee", key)

	if len(key) > 0 {
		trans, err = getTransaction(key)
		if err != nil {
			return nil, wsapi.NewCustomInternalError("Failure to locate the transaction")
		}
	}

	fee, err := Wallet.GetFee()
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	if trans != nil {
		ufee, _ := trans.CalculateFee(uint64(fee))
		fee = int64(ufee)
	}

	resp := new(GetFeeResponse)
	resp.Fee = int64(fee)
	return resp, nil
}

func GetAddresses() []byte {
	addResp, err := GetV2Addresses()
	if err != nil {
		panic(err)
	}

	var maxlen int

	for i := range addResp.EntryCreditAddressed {
		if len(addResp.EntryCreditAddressed[i].Key) > maxlen {
			maxlen = len(addResp.EntryCreditAddressed[i].Key)
		}
	}
	for i := range addResp.FactoidAddressed {
		if len(addResp.FactoidAddressed[i].Key) > maxlen {
			maxlen = len(addResp.FactoidAddressed[i].Key)
		}
	}

	var out bytes.Buffer
	if len(addResp.FactoidAddressed) > 0 {
		out.WriteString("\n  Factoid Addresses\n\n")
	}
	fstr := fmt.Sprintf("%s%vs    %s38s %s14s\n", "%", maxlen+4, "%", "%")
	for _, fAdd := range addResp.FactoidAddressed {
		bal := primitives.ConvertDecimalToString(uint64(fAdd.BalanceDecimal))
		str := fmt.Sprintf(fstr, fAdd.Key, fAdd.Address, bal)
		out.WriteString(str)
	}
	if len(addResp.EntryCreditAddressed) > 0 {
		out.WriteString("\n  Entry Credit Addresses\n\n")
	}
	for _, ecAdd := range addResp.EntryCreditAddressed {
		bal := primitives.ConvertDecimalToString(uint64(ecAdd.BalanceDecimal))
		str := fmt.Sprintf(fstr, ecAdd.Key, ecAdd.Address, bal)
		out.WriteString(str)
	}

	return out.Bytes()
}

func GetV2Addresses() (*AddressesResponse, error) {
	values, err := Wallet.GetAddresses()
	if err != nil {
		return nil, err
	}

	eca := make([]Address, 0, len(values))
	fa := make([]Address, 0, len(values))

	for _, we := range values {
		var add Address
		if we.GetType() == "ec" {
			address, err := we.GetAddress()
			if err != nil {
				continue
			}
			add.Address = primitives.ConvertECAddressToUserStr(address)
			add.Key = string(we.GetName())
			add.BalanceDecimal, err = ECBalance(add.Address)
			if err != nil {
				return nil, err
			}
			add.Balance = primitives.ConvertDecimalToFloat(uint64(add.BalanceDecimal))
			eca = append(eca, add)
		} else {
			address, err := we.GetAddress()
			if err != nil {
				continue
			}
			add.Address = primitives.ConvertFctAddressToUserStr(address)
			add.Key = string(we.GetName())
			add.BalanceDecimal, err = FctBalance(add.Address)
			if err != nil {
				return nil, err
			}
			add.Balance = primitives.ConvertDecimalToFloat(uint64(add.BalanceDecimal))
			fa = append(fa, add)
		}
	}

	resp := new(AddressesResponse)
	resp.EntryCreditAddressed = eca
	resp.FactoidAddressed = fa

	return resp, nil
}

// Specifying a fee overrides either not being connected, or the current fee.
// Params:
//   key (limit printout to this key)
//   fee (specify the transation fee)
func GetTransactions(ctx *web.Context) ([]byte, error) {
	connected := true

	var _ = connected

	exch, err := GetFee(ctx) // The Fee will be zero if we have no connection.
	if err != nil {
		connected = false
	}

	keys, transactions, err := Wallet.GetTransactions()
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	for i, trans := range transactions {

		fee, _ := trans.CalculateFee(uint64(exch))
		cprt := ""
		cin, err := trans.TotalInputs()
		if err != nil {
			cprt = cprt + err.Error()
		}
		cout, err := trans.TotalOutputs()
		if err != nil {
			cprt = cprt + err.Error()
		}
		cecout, err := trans.TotalECs()
		if err != nil {
			cprt = cprt + err.Error()
		}

		if len(cprt) == 0 {
			v := int64(cin) - int64(cout) - int64(cecout)
			sign := ""
			if v < 0 {
				sign = "-"
				v = -v
			}
			cprt = fmt.Sprintf(" Currently will pay: %s%s",
				sign,
				primitives.ConvertDecimalToString(uint64(v)))
			if sign == "-" || fee > uint64(v) {
				cprt = cprt + "\n\nWARNING: Currently your transaction fee may be too low"
			}
		}

		out.WriteString(fmt.Sprintf("%s:  Fee Due: %s  %s\n\n%s\n",
			strings.TrimSpace(strings.TrimRight(string(keys[i]), "\u0000")),
			primitives.ConvertDecimalToString(fee),
			cprt,
			transactions[i].String()))
	}

	output := out.Bytes()
	// now look for the addresses, and replace them with our names. (the transactions
	// in flight also have a Factom address... We leave those alone.

	names, vs, err := Wallet.GetWalletNames()
	if err != nil {
		return nil, err
	}

	for i, name := range names {
		we, ok := vs[i].(interfaces.IWalletEntry)
		if !ok {
			return nil, fmt.Errorf("Database is corrupt")
		}

		address, err := we.GetAddress()
		if err != nil {
			continue
		} // We shouldn't get any of these, but ignore them if we do.
		adrstr := []byte(hex.EncodeToString(address.Bytes()))

		output = bytes.Replace(output, adrstr, name, -1)
	}

	return output, nil
}

// Specifying a fee overrides either not being connected, or the current fee.
// Params:
//   key (limit printout to this key)
//   fee (specify the transation fee)
func GetTransactionsj(ctx *web.Context) (string, error) {
	connected := true

	var _ = connected

	keys, transactions, _ := Wallet.GetTransactions()
	type pair struct {
		Key     string
		TransID string
	}
	var trans []*pair
	for i, t := range transactions {
		p := new(pair)
		p.Key = strings.TrimRight(string(keys[i]), "\u0000")
		p.TransID = t.GetSigHash().String()
		trans = append(trans, p)
	}

	return primitives.EncodeJSONString(trans)

}

func HandleGetAddresses(ctx *web.Context) {
	b := new(Response)
	b.Response = string(GetAddresses())
	b.Success = true
	j, err := json.Marshal(b)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}
	ctx.ContentType("json")
	ctx.Write(j)
}

func HandleV2GetAddresses(params interface{}) (interface{}, *primitives.JSONError) {
	resp, err := GetV2Addresses()
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}
	return resp, nil
}

func HandleGetTransactionsj(ctx *web.Context) {
	b := new(Response)
	txt, err := GetTransactionsj(ctx)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}
	b.Response = txt
	b.Success = true
	j, err := json.Marshal(b)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}
	ctx.Write(j)
}

func HandleGetTransactions(ctx *web.Context) {
	b := new(Response)
	txt, err := GetTransactions(ctx)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}
	b.Response = string(txt)
	b.Success = true
	j, err := json.Marshal(b)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}
	ctx.ContentType("json")
	ctx.Write(j)
}

func HandleFactoidValidate(ctx *web.Context) {
}

func HandleV2FactoidValidate(params interface{}) (interface{}, *primitives.JSONError) {
	return "OK", nil
}
