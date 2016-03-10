package scwallet

import (

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/database/databaseOverlay"
	"github.com/FactomProject/factomd/common/factoid"
)

type SCDatabaseOverlay struct {
	databaseOverlay.Overlay
}

var _ interfaces.ISCDatabaseOverlay = (*SCDatabaseOverlay)(nil)

func NewSCOverlay(db interfaces.IDatabase) interfaces.ISCDatabaseOverlay {
	answer := new(SCDatabaseOverlay)
	answer.DB = db
	return answer
}

//Wallet Entries

func (sc *SCDatabaseOverlay) FetchWalletEntryByName(addr []byte) (interfaces.IWalletEntry, error) {
	we, err := sc.DB.Get([]byte(constants.W_NAME), addr, new(WalletEntry))
	if err != nil {
		return nil, err
	}
	if we == nil {
		return nil, nil
	}
	return we.(interfaces.IWalletEntry), nil
}

func (sc *SCDatabaseOverlay) FetchWalletEntryByPublicKey(addr []byte) (interfaces.IWalletEntry, error) {
	we, err := sc.DB.Get([]byte(constants.W_ADDRESS_PUB_KEY), addr, new(WalletEntry))
	if err != nil {
		return nil, err
	}
	if we == nil {
		return nil, nil
	}
	return we.(interfaces.IWalletEntry), nil
}

func (sc *SCDatabaseOverlay) FetchAllWalletEntriesByName() ([]interfaces.IWalletEntry, error) {
	values, err := sc.DB.GetAll([]byte(constants.W_NAME), new(WalletEntry))
	if err != nil {
		return nil, err
	}
	answerWE := []interfaces.IWalletEntry{}
	for _, v := range values {
		we, ok := v.(interfaces.IWalletEntry)
		if !ok {
			panic("Get Addresses finds the database corrupt. Shouldn't happen")
		}
		answerWE = append(answerWE, we)
	}
	return answerWE, nil
}

func (sc *SCDatabaseOverlay) FetchAllWalletEntriesByPublicKey() ([]interfaces.IWalletEntry, error) {
	values, err := sc.DB.GetAll([]byte(constants.W_ADDRESS_PUB_KEY), new(WalletEntry))
	if err != nil {
		return nil, err
	}
	answerWE := []interfaces.IWalletEntry{}
	for _, v := range values {
		we, ok := v.(interfaces.IWalletEntry)
		if !ok {
			panic("Get Addresses finds the database corrupt. Shouldn't happen")
		}
		answerWE = append(answerWE, we)
	}
	return answerWE, nil
}

func (sc *SCDatabaseOverlay) FetchAllAddressNameKeys() ([][]byte, error) {
	return sc.DB.ListAllKeys([]byte(constants.W_NAME))
}

func (sc *SCDatabaseOverlay) FetchAllAddressPublicKeys() ([][]byte, error) {
	return sc.DB.ListAllKeys([]byte(constants.W_ADDRESS_PUB_KEY))
}

//Transactions

func (sc *SCDatabaseOverlay) FetchTransaction(key []byte) (interfaces.ITransaction, error) {
	we, err := sc.DB.Get([]byte(constants.DB_BUILD_TRANS), key, new(factoid.Transaction))
	if err != nil {
		return nil, err
	}
	if we == nil {
		return nil, nil
	}
	return we.(*factoid.Transaction), nil
}

func (sc *SCDatabaseOverlay) SaveTransaction(key []byte, tx interfaces.ITransaction) (error) {
	return sc.DB.Put([]byte(constants.DB_BUILD_TRANS), key, tx)
}

func (sc *SCDatabaseOverlay) DeleteTransaction(key []byte) (error) {
	return sc.DB.Delete([]byte(constants.DB_BUILD_TRANS), key)
}

func (sc *SCDatabaseOverlay) FetchAllTransactionKeys() ([][]byte, error) {
	return sc.DB.ListAllKeys([]byte(constants.DB_BUILD_TRANS))
}

func (sc *SCDatabaseOverlay) FetchAllTransactions() ([]interfaces.ITransaction, error) {
	values, err := sc.DB.GetAll([]byte(constants.DB_BUILD_TRANS), new(factoid.Transaction))
	if err != nil {
		return nil, err
	}
	answerWE := []interfaces.ITransaction{}
	for _, v := range values {
		we, ok := v.(interfaces.ITransaction)
		if !ok {
			panic("Get Addresses finds the database corrupt. Shouldn't happen")
		}
		answerWE = append(answerWE, we)
	}
	return answerWE, nil
}