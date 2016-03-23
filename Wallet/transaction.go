// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
)

// New Transaction:  key --
// We create a new transaction, and track it with the user supplied key.  The
// user can then use this key to make subsequent calls to add inputs, outputs,
// and to sign. Then they can submit the transaction.
//
// When the transaction is submitted, we clear it from our working memory.
// Multiple transactions can be under construction at one time, but they need
// their own keys. Once a transaction is either submitted or deleted, the key
// can be reused.
func FactoidNewTransaction(key string) error {
	// Make sure we have a key
	if len(key) == 0 {
		return fmt.Errorf("Missing transaction key")
	}

	ok := Utility.IsValidKey(key)
	if !ok {
		return fmt.Errorf("Invalid name for transaction")
	}

	// Make sure we don't already have a transaction in process with this key
	t, err := wallet.GetDB().FetchTransaction([]byte(key))
	if err != nil {
		return err
	}
	if t != nil {
		return fmt.Errorf("Duplicate key: '%s'", key)
	}
	// Create a transaction
	t = wallet.CreateTransaction(interfaces.GetTimeMilli())
	// Save it with the key
	return wallet.GetDB().SaveTransaction([]byte(key), t)
}

// Delete Transaction:  key --
// Remove a transaction rather than sign and submit the transaction.  Sometimes
// you just need to throw one a way, and rebuild it.
//
func FactoidDeleteTransaction(key string) error {
	// Make sure we have a key
	if len(key) == 0 {
		return fmt.Errorf("Missing transaction key")
	}
	// Wipe out the key
	return wallet.GetDB().DeleteTransaction([]byte(key))
}

func FactoidAddFee(trans interfaces.ITransaction, key string, address interfaces.IAddress, name string) (uint64, error) {
	ins, err := trans.TotalInputs()
	if err != nil {
		return 0, err
	}
	outs, err := trans.TotalOutputs()
	if err != nil {
		return 0, err
	}
	ecs, err := trans.TotalECs()
	if err != nil {
		return 0, err
	}

	if ins != outs+ecs {
		return 0, fmt.Errorf("Inputs and outputs don't add up")
	}

	ok := Utility.IsValidKey(key)
	if !ok {
		return 0, fmt.Errorf("Invalid name for transaction")
	}

	fee, err := GetFee()
	if err != nil {
		return 0, err
	}

	transfee, err := trans.CalculateFee(uint64(fee))
	if err != nil {
		return 0, err
	}

	adr, err := wallet.GetAddressHash(address)
	if err != nil {
		return 0, err
	}

	for _, input := range trans.GetInputs() {
		if input.GetAddress().IsSameAs(adr) {
			amt, err := factoid.ValidateAmounts(input.GetAmount(), transfee)
			if err != nil {
				return 0, err
			}
			input.SetAmount(amt)

			// Update our map with our new transaction to the same key. Otherwise, all
			// of our work will go away!
			err = wallet.GetDB().SaveTransaction([]byte(key), trans)
			if err != nil {
				return 0, err
			}

			return transfee, nil
		}
	}
	return 0, fmt.Errorf("%s is not an input to the transaction.", key)
}

func FactoidSubFee(trans interfaces.ITransaction, key string, address interfaces.IAddress, name string) (uint64, error) {
	{
		ins, err := trans.TotalInputs()
		if err != nil {
			return 0, err
		}
		outs, err := trans.TotalOutputs()
		if err != nil {
			return 0, err
		}
		ecs, err := trans.TotalECs()
		if err != nil {
			return 0, err
		}

		if ins != outs+ecs {
			return 0, fmt.Errorf("Inputs and outputs don't add up")
		}
	}

	ok := Utility.IsValidKey(key)
	if !ok {
		return 0, fmt.Errorf("Invalid name for transaction")
	}

	fee, err := GetFee()
	if err != nil {
		return 0, err
	}

	transfee, err := trans.CalculateFee(uint64(fee))
	if err != nil {
		return 0, err
	}

	adr, err := wallet.GetAddressHash(address)
	if err != nil {
		return 0, err
	}

	for _, output := range trans.GetOutputs() {
		if output.GetAddress().IsSameAs(adr) {
			output.SetAmount(output.GetAmount() - transfee)
			return transfee, nil
		}
	}
	return 0, fmt.Errorf("%s is not an input to the transaction.", key)
}

func FactoidAddInput(trans interfaces.ITransaction, key string, address interfaces.IAddress, amount uint64) error {
	ok := Utility.IsValidKey(key)
	if !ok {
		return fmt.Errorf("Invalid name for transaction")
	}

	// First look if this is really an update
	for _, input := range trans.GetInputs() {
		if input.GetAddress().IsSameAs(address) {
			input.SetAmount(amount)
			return nil
		}
	}

	// Add our new input
	err := wallet.AddInput(trans, address, amount)
	if err != nil {
		return fmt.Errorf("Failed to add input")
	}

	// Update our map with our new transaction to the same key. Otherwise, all
	// of our work will go away!
	return wallet.GetDB().SaveTransaction([]byte(key), trans)
}

func FactoidAddOutput(trans interfaces.ITransaction, key string, address interfaces.IAddress, amount uint64) error {
	ok := Utility.IsValidKey(key)
	if !ok {
		return fmt.Errorf("Invalid name for transaction")
	}

	// First look if this is really an update
	for _, output := range trans.GetOutputs() {
		if output.GetAddress().IsSameAs(address) {
			output.SetAmount(amount)
			return nil
		}
	}
	// Add our new Output
	err := wallet.AddOutput(trans, address, uint64(amount))
	if err != nil {
		return fmt.Errorf("Failed to add output")
	}

	// Update our map with our new transaction to the same key.  Otherwise, all
	// of our work will go away!
	return wallet.GetDB().SaveTransaction([]byte(key), trans)
}

func FactoidAddECOutput(trans interfaces.ITransaction, key string, address interfaces.IAddress, amount uint64) error {
	ok := Utility.IsValidKey(key)
	if !ok {
		return fmt.Errorf("Invalid name for transaction")
	}
	// First look if this is really an update
	for _, ecoutput := range trans.GetECOutputs() {
		if ecoutput.GetAddress().IsSameAs(address) {
			ecoutput.SetAmount(amount)
			return nil
		}
	}
	// Add our new Entry Credit Output
	err := wallet.AddECOutput(trans, address, uint64(amount))
	if err != nil {
		return fmt.Errorf("Failed to add Entry Credit Output")
	}

	// Update our map with our new transaction to the same key.  Otherwise, all
	// of our work will go away!
	return wallet.GetDB().SaveTransaction([]byte(key), trans)
}

func FactoidSignTransaction(key string) error {
	ok := Utility.IsValidKey(key)
	if !ok {
		return fmt.Errorf("Invalid name for transaction")
	}

	// Get the transaction
	trans, err := GetTransaction(key)
	if err != nil {
		return fmt.Errorf("Failed to get the transaction")
	}

	err = wallet.Validate(1, trans)
	if err != nil {
		return err
	}

	valid, err := wallet.SignInputs(trans)
	if !valid {
		return fmt.Errorf("Do not have all the private keys required to sign this transaction\n" +
			err.Error())
	}
	if err != nil {
		return err
	}

	err = wallet.ValidateSignatures(trans)
	if err != nil {
		fmt.Printf("FactoidSignTransaction - Signature invalid - %v, %v\n", trans, err)
		return err
	}
	// Update our map with our new transaction to the same key.  Otherwise, all
	// of our work will go away!
	return wallet.GetDB().SaveTransaction([]byte(key), trans)
}

func FactoidSubmit(jsonkey string) (string, error) {
	type submitReq struct {
		Transaction string
	}

	in := new(submitReq)
	json.Unmarshal([]byte(jsonkey), in)

	key := in.Transaction
	// Get the transaction
	trans, err := GetTransaction(key)
	if err != nil {
		return "", err
	}

	fmt.Printf("Fetched transaction - %v\n", trans)

	err = wallet.ValidateSignatures(trans)
	if err != nil {
		fmt.Printf("Signature invalid - %v\n", err)
		fmt.Println(err)
		return "", err
	}

	err = isReasonableFee(trans)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Okay, transaction is good, so marshal and send to factomd!
	data, err := trans.MarshalBinary()
	if err != nil {
		fmt.Printf("Error marshalling transaction - %v\n", err)
		fmt.Println(err)
		return "", err
	}

	transdata := string(hex.EncodeToString(data))

	s := struct{ Transaction string }{transdata}

	j, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Printf("Encoded transaction - %v\n", j)

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/v1/factoid-submit/", ipaddressFD, portNumberFD),
		"application/json",
		bytes.NewBuffer(j))

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	resp.Body.Close()

	r := new(Response)
	if err := json.Unmarshal(body, r); err != nil {
		return "", err
	}

	if r.Success {
		wallet.GetDB().DeleteTransaction([]byte(key))
		return "", nil
	} else {
		return "", fmt.Errorf(r.Response)
	}
	return r.Response, nil
}

func isReasonableFee(trans interfaces.ITransaction) error {
	feeRate, getErr := GetFee()
	if getErr != nil {
		return getErr
	}

	reqFee, err := trans.CalculateFee(uint64(feeRate))
	if err != nil {
		return err
	}

	sreqFee := int64(reqFee)

	tin, err := trans.TotalInputs()
	if err != nil {
		return err
	}

	tout, err := trans.TotalOutputs()
	if err != nil {
		return err
	}

	tec, err := trans.TotalECs()
	if err != nil {
		return err
	}

	cfee := int64(tin) - int64(tout) - int64(tec)

	if cfee >= (sreqFee * 10) {
		return fmt.Errorf("Unbalanced transaction (fee too high). Fee should be less than 10x the required fee.")
	}

	if cfee < sreqFee {
		return fmt.Errorf("Insufficient fee - %v vs %v", cfee, sreqFee)
	}

	return nil
}

func GetFee() (int64, error) {
	str := fmt.Sprintf("http://%s:%d/v1/factoid-get-fee/", ipaddressFD, portNumberFD)
	resp, err := http.Get(str)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return 0, err
	}
	resp.Body.Close()

	type x struct {
		Response struct {
			Fee int64
		}
		Success bool
	}
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return 0, err
	}

	return b.Response.Fee, nil
}

func GetProperties() (protocol, factomd, fctwallet string, err error) {
	type prop struct {
		Protocol_Version  string
		Factomd_Version   string
		Fctwallet_Version string
	}

	str := fmt.Sprintf("http://%s:%d/v1/properties/", ipaddressFD, portNumberFD)
	resp, err := http.Get(str)
	if err != nil {
		return "", "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return "", "", "", err
	}
	resp.Body.Close()

	b := new(prop)
	if err := json.Unmarshal(body, b); err != nil {
		return "", "", "", err
	}
	b.Fctwallet_Version = Version
	return b.Protocol_Version, b.Factomd_Version, Version, nil
}

func GetAddresses() ([]interfaces.IWalletEntry, error) {
	return wallet.GetDB().FetchAllWalletEntriesByName()
}

func GetTransactions() ([][]byte, []interfaces.ITransaction, error) {
	// Get the transactions in flight.
	keys, err := wallet.GetDB().FetchAllTransactionKeys()
	if err != nil {
		return nil, nil, err
	}
	values, err := wallet.GetDB().FetchAllTransactions()
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(keys)-1; i++ {
		for j := 0; j < len(keys)-i-1; j++ {
			if bytes.Compare(keys[j], keys[j+1]) > 0 {
				t := keys[j]
				keys[j] = keys[j+1]
				keys[j+1] = t
				t2 := values[j]
				values[j] = values[j+1]
				values[j+1] = t2
			}
		}
	}

	return keys, values, nil
}

func GetWalletNames() ([][]byte, []interfaces.IWalletEntry, error) {
	keys, err := wallet.GetDB().FetchAllAddressNameKeys()
	if err != nil {
		return nil, nil, err
	}
	values, err := wallet.GetDB().FetchAllWalletEntriesByName()
	if err != nil {
		return nil, nil, err
	}

	return keys, values, nil
}

func GenerateFctAddress(name []byte, m int, n int) (hash interfaces.IAddress, err error) {
	return wallet.GenerateFctAddress(name, m, n)
}

func NewSeed(data []byte) {
	wallet.NewSeed(data)
}
