// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/util"
	//"github.com/FactomProject/factoid/state/stateinit"
)

const (
	httpOK  = 200
	httpBad = 400
)

var (
	cfg              = util.ReadConfig("")
	IpAddress        = cfg.Wallet.Address
	PortNumber       = cfg.Wallet.Port
	applicationName  = "Factom/fctwallet"
	dataStorePath    = cfg.Wallet.DataFile
	refreshInSeconds = cfg.Wallet.RefreshInSeconds
	ipaddressFD      = cfg.Wallet.Address
	portNumberFD     = cfg.Wsapi.PortNumber

	//databasefile = "factoid_wallet_bolt.db"
)

func init() {
	factom.SetServer(fmt.Sprintf("%v:%v", ipaddressFD, portNumberFD))
}

/*
var factoidState = stateinit.NewFactoidState(cfg.BoltDBPath + databasefile)
*/
