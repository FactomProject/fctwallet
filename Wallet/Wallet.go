package Wallet

import (
	"fmt"
	"regexp"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
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

func GetTransaction(key string) (interfaces.ITransaction, error) {
	ok := Utility.IsValidKey(key)
	if !ok {
		return nil, fmt.Errorf("Invalid name or address")
	}

	// Now get the transaction.  If we don't have a transaction by the given
	// keys there is nothing we can do.  Now we *could* create the transaaction
	// and tie it to the key.  Something to think about.
	return wallet.GetDB().FetchTransaction([]byte(key))
}

func GetWalletEntry(key []byte) (interfaces.IWalletEntry, error) {
	return wallet.GetDB().FetchWalletEntryByName(key)
}
