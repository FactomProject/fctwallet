// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"github.com/FactomProject/FactomCode/util"
	"github.com/FactomProject/factoid/state/stateinit"
)

var (
	cfg             = util.ReadConfig().Wallet
	applicationName = "Factom/fctwallet"

	ipaddressFD  = cfg.FactomdAddress
	portNumberFD = cfg.FactomdPort

	databasefile = "factoid_wallet_bolt.db"
)

var factoidState = stateinit.NewFactoidState(cfg.BoltDBPath + databasefile)

const Version = "0.1.8.0"
