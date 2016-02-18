package Wallet

import (
	"fmt"
	"regexp"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
	"github.com/FactomProject/fctwallet/scwallet"
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

func ValidateKey(key string) error {
	if Utility.IsValidKey(key) {
		return nil
	}
	return fmt.Errorf("Invalid key")
}

func GetTransaction(key string) (trans interfaces.ITransaction, err error) {
	ok := Utility.IsValidKey(key)
	if !ok {
		return nil, fmt.Errorf("Invalid name or address")
	}

	// Now get the transaction.  If we don't have a transaction by the given
	// keys there is nothing we can do.  Now we *could* create the transaaction
	// and tie it to the key.  Something to think about.
	ib, err := wallet.GetDB().Get([]byte(constants.DB_BUILD_TRANS), []byte(key), new(factoid.Transaction))
	if err != nil {
		return nil, err
	}

	trans, ok = ib.(interfaces.ITransaction)
	if ib == nil || !ok {
		return nil, fmt.Errorf("Unknown Transaction: %s", key)
	}
	return
}

func GetWalletEntry(key []byte) (interfaces.BinaryMarshallable, error) {
	return wallet.GetDB().Get([]byte(constants.W_NAME), key, new(scwallet.WalletEntry))
}
