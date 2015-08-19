package main_test

import (
	"github.com/FactomProject/fctwallet/Wallet"
	"testing"
)

func Test(t *testing.T) {
	err := Wallet.Synchronize()
	t.Errorf("Test - %v", err)
}
