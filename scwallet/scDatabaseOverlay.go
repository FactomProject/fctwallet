package scwallet

import (

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/database/databaseOverlay"
)

type SCDatabaseOverlay struct {
	databaseOverlay.Overlay
}

func (sc *SCDatabaseOverlay) FetchWalletEntry(addr string) (interfaces.IWalletEntry, error) {
	we, err := db.DB.Get([]byte(constants.W_NAME), []byte(adr), new(WalletEntry))
	if err != nil {
		return nil, err
	}
	if we == nil {
		return nil, nil
	}
	return we.(interfaces.IWalletEntry), nil
}