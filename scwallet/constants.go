// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package scwallet

const (
	// Wallet
	//==================
	W_SEEDS            = "wallet.address.seeds"      // Holds the root seeds for address generation
	W_SEED_HEADS       = "wallet.address.seed.heads" // Holds the latest generated seed for each root seed.
	W_RCD_ADDRESS_HASH = "wallet.address.addr"
	W_ADDRESS_PUB_KEY  = "wallet.public.key"
	W_NAME             = "wallet.address.name"
	DB_BUILD_TRANS     = "Transactions_Under_Construction"
	DB_TRANSACTIONS    = "Transactions_For_Addresses"

	PRIVATE_LENGTH = 64 // length of a Private Key
	ADDRESS_LENGTH = 32 // Length of an Address or a Hash or Public Key
)

var CURRENT_SEED = [32]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
