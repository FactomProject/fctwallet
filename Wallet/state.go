// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package Wallet

import (
	"github.com/FactomProject/FactomCode/util"
	"github.com/FactomProject/factoid/state"
	"github.com/FactomProject/factoid/state/stateinit"
)

var (
	cfg             = util.ReadConfig().Wallet
	applicationName = "Factom/fctwallet"

	ipaddressFD  = "localhost:"
	portNumberFD = "8088"

	databasefile = "factoid_wallet_bolt.db"
)

var factoidState state.IFactoidState

func init() {
	factoidState = stateinit.NewFactoidState(cfg.BoltDBPath + databasefile)
	Init()
}
