// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"github.com/FactomProject/factomd/util"
	"github.com/FactomProject/fctwallet/scwallet"
)

var (

	applicationName = "Factom/fctwallet"

	databasefile = "factoid_wallet_bolt.db"

	cfg              = util.ReadConfig("")
	ipaddressFD      = cfg.Wallet.Address
	portNumberFD     = cfg.Wsapi.PortNumber
)

var wallet = scwallet.NewSCWallet(cfg.App.BoltDBPath, databasefile)

const Version = "0.1.6.0"
