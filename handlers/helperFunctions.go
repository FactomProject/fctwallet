// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	//"time"
	"github.com/FactomProject/web"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/fctwallet/Wallet"
)

/******************************************
 * Helper Functions
 ******************************************/

var badChar, _ = regexp.Compile("[^A-Za-z0-9_-]")
var badHexChar, _ = regexp.Compile("[^A-Fa-f0-9]")

type Response struct {
	Response string
	Success  bool
}

func ValidateKey(key string) (msg string, valid bool) {
	if len(key) > constants.ADDRESS_LENGTH {
		return "Key is too long.  Keys must be less than 32 characters", false
	}
	if badChar.FindStringIndex(key) != nil {
		str := fmt.Sprintf("The key or name '%s' contains invalid characters.\n"+
			"Keys and names are restricted to alphanumeric characters,\n"+
			"minuses (dashes), and underscores", key)
		return str, false
	}
	return "", true
}

// True is success! False is failure.  The Response is what the CLI
// should report.
func reportResults(ctx *web.Context, response string, success bool) {
	b := Response{
		Response: response,
		Success:  success,
	}
	if p, err := json.Marshal(b); err != nil {

		ctx.WriteHeader(httpBad)
		return
	} else {
		ctx.ContentType("json")
		ctx.Write(p)
	}
}

func getTransaction(key string) (interfaces.ITransaction, error) {
	return Wallet.GetTransaction(key)
}

// &key=<key>&name=<name or address>&amount=<amount>
// If no amount is specified, a zero is returned.
func getParams_(ctx *web.Context, params string, ec bool) (
	trans interfaces.ITransaction,
	key string,
	name string,
	address interfaces.IAddress,
	amount int64,
	ok bool) {

	key = ctx.Params["key"]
	name = ctx.Params["name"]
	StrAmount := ctx.Params["amount"]

	if len(StrAmount) == 0 {
		StrAmount = "0"
	}

	if len(key) == 0 || len(name) == 0 {
		str := fmt.Sprintln("Missing Parameters: key='", key, "' name='", name, "' amount='", StrAmount, "'")
		reportResults(ctx, str, false)
		ok = false
		return
	}

	msg, valid := ValidateKey(key)
	if !valid {
		reportResults(ctx, msg, false)
		ok = false
		return
	}

	amount, err := strconv.ParseInt(StrAmount, 10, 64)
	if err != nil {
		str := fmt.Sprintln("Error parsing amount.\n", err)
		reportResults(ctx, str, false)
		ok = false
		return
	}

	// Get the transaction
	trans, err = getTransaction(key)
	if err != nil {
		reportResults(ctx, "Failure to locate the transaction", false)
		ok = false
		return
	}

	// Get the input/output/ec address.  Which could be a name.  First look and see if it is
	// a name.  If it isn't, then look and see if it is an address.  Someone could
	// do a weird Address as a name and fool the code, but that seems unlikely.
	// Could check for that some how, but there are many ways around such checks.

	if len(name) <= constants.ADDRESS_LENGTH {
		we, err := Wallet.GetWalletEntry([]byte(name))
		if err != nil {
			reportResults(ctx, "Failure to locate the transaction", false)
			ok = false
			return
		}
		if we != nil {
			address, err = we.GetAddress()
			if we.GetType() == "ec" {
				if !ec {
					reportResults(ctx, "Was Expecting a Factoid Address", false)
					ok = false
					return
				}
			} else {
				if ec {
					reportResults(ctx, "Was Expecting an Entry Credit Address", false)
					ok = false
					return
				}
			}
			if err != nil || address == nil {
				reportResults(ctx, "Should not get an error geting a address from a Wallet Entry", false)
				ok = false
				return
			}
			ok = true
			return
		}
	}
	if (!ec && !primitives.ValidateFUserStr(name)) || (ec && !primitives.ValidateECUserStr(name)) {
		reportResults(ctx, fmt.Sprintf("The address specified isn't defined or is invalid: %s", name), false)
		ctx.WriteHeader(httpBad)
		ok = false
		return
	}
	baddr := primitives.ConvertUserStrToAddress(name)

	address = factoid.NewAddress(baddr)

	ok = true
	return
}

/*************************************************************************
 * Handler Functions
 *************************************************************************/
