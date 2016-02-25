// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hoisie/web"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/fctwallet/Wallet"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/directoryBlock"
	"github.com/FactomProject/factomd/common/factoid/block"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/wsapi"
)

// Older blocks smaller indexes.  All the Factoid Directory blocks
var DirectoryBlocks = make([]*directoryBlock.DirectoryBlock, 0, 100)
var FactoidBlocks = make([]interfaces.IFBlock, 0, 100)
var DBHead []byte
var DBHeadStr string = ""

// Refresh the Directory Block Head.  If it has changed, return true.
// Otherwise return false.
func getDBHead() bool {
	db, err := factom.GetDBlockHead()

	if err != nil {
		panic(err.Error())
	}

	if db != DBHeadStr {
		DBHeadStr = db
		DBHead, err = hex.DecodeString(db)
		if err != nil {
			panic(err.Error())
		}

		return true
	}
	return false
}

func getAll() error {
	dbs := make([]*directoryBlock.DirectoryBlock, 0, 100)
	next := DBHeadStr

	for {
		blk, err := factom.GetRaw(next)
		if err != nil {
			panic(err.Error())
		}
		db := new(directoryBlock.DirectoryBlock)
		err = db.UnmarshalBinary(blk)
		if err != nil {
			panic(err.Error())
		}
		dbs = append(dbs, db)
		if bytes.Equal(db.GetHeader().GetPrevKeyMR().Bytes(), constants.ZERO_HASH[:]) {
			break
		}
		next = hex.EncodeToString(db.GetHeader().GetPrevKeyMR().Bytes())
	}

	for i := len(dbs) - 1; i >= 0; i-- {
		DirectoryBlocks = append(DirectoryBlocks, dbs[i])
		fb := new(block.FBlock)
		for _, dbe := range dbs[i].DBEntries {
			if bytes.Equal(dbe.GetChainID().Bytes(), constants.FACTOID_CHAINID) {
				hashstr := hex.EncodeToString(dbe.GetKeyMR().Bytes())
				fdata, err := factom.GetRaw(hashstr)
				if err != nil {
					panic(err.Error())
				}
				err = fb.UnmarshalBinary(fdata)
				if err != nil {
					panic(err.Error())
				}
				FactoidBlocks = append(FactoidBlocks, fb)
				break
			}
		}
		if fb == nil {
			fmt.Println("Missing Factoid Block")
		}
	}
	return nil
}

func refresh() error {
	if DBHead == nil {
		getDBHead()
		getAll()
	}
	if getDBHead() {

	}
	return nil
}

func FctBalance(adr string) (int64, error) {
	err := refresh()
	if err != nil {
		return 0, err
	}
	return Wallet.FactoidBalance(adr)
}

func ECBalance(adr string) (int64, error) {
	return Wallet.ECBalance(adr)
}

func HandleEntryCreditBalance(ctx *web.Context, adr string) {
	req := primitives.NewJSON2Request(1, adr, "entry-credit-balance")

	jsonResp, jsonError := HandleV2GetRequest(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	str := fmt.Sprintf("%d", jsonResp.Result.(*EntryCreditBalanceResponse).Balance)
	reportResults(ctx, str, true)
}

func HandleFactoidBalance(ctx *web.Context, adr string) {
	req := primitives.NewJSON2Request(1, adr, "factoid-balance")

	jsonResp, jsonError := HandleV2GetRequest(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	str := fmt.Sprintf("%d", jsonResp.Result.(*FactoidCreditBalance).Balance)
	reportResults(ctx, str, true)
}


func HandleV2EntryCreditBalance(params interface{}) (interface{}, *primitives.JSONError) {
	adr, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	v, err := ECBalance(adr)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}
	resp:=new(EntryCreditBalanceResponse)
	resp.Balance = v

	return resp, nil
}

func HandleV2FactoidBalance(params interface{}) (interface{}, *primitives.JSONError) {
	adr, ok := params.(string)
	if ok == false {
		return nil, wsapi.NewInvalidParamsError()
	}

	v, err := FctBalance(adr)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}
	resp:=new(FactoidBalanceResponse)
	resp.Balance = v

	return resp, nil
}
