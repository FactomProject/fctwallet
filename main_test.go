package main_test

import (
	"fmt"
	"github.com/FactomProject/fctwallet/Wallet"
	"testing"
)

func Test(t *testing.T) {
	initWallet()
	err := Wallet.Synchronize()
	t.Errorf("Test - %v", err)
}

func initWallet() {
	fmt.Printf("\ninitWallet\n")
	keys, _ := Wallet.GetWalletNames()
	if len(keys) == 0 {
		for i := 1; i <= 10; i++ {
			name := fmt.Sprintf("%02d-Fountain", i)
			_, err := Wallet.GenerateFctAddress([]byte(name), 1, 1)
			if err != nil {
				fmt.Printf("\nError - %v\n", err)
				return
			}
		}
	}
	fmt.Printf("initWallet done\n")
}
