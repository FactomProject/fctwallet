// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/FactomProject/factomd/util"
	//"github.com/FactomProject/factoid/state/stateinit"
)

const (
	httpOK  = 200
	httpBad = 400
)

var (
	cfg              = util.ReadConfig("").Wallet
	IpAddress        = cfg.Address
	PortNumber       = cfg.Port
	applicationName  = "Factom/fctwallet"
	dataStorePath    = cfg.DataFile
	refreshInSeconds = cfg.RefreshInSeconds
	ipaddressFD      = cfg.FactomdAddress
	portNumberFD     = cfg.FactomdPort

	//databasefile = "factoid_wallet_bolt.db"
)

/*
var factoidState = stateinit.NewFactoidState(cfg.BoltDBPath + databasefile)
*/
