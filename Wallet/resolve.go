// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"fmt"
	netki "github.com/netkicorp/go-partner-client"
)

func NetkiResolve(adr string) (string, string, error) {
	fct, ferr := netki.WalletNameLookup(adr, "fct")
	ec, eerr := netki.WalletNameLookup(adr, "fec")
	if ferr != nil && eerr != nil {
		return fct, ec, fmt.Errorf("%s\n%s", ferr, eerr)
	}
	return fct, ec, nil
}
