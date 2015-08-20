package Wallet_test

import (
	//"fmt"
	"github.com/FactomProject/fctwallet/Wallet"
	"testing"
)

func TestValidateKey(t *testing.T) {
	validKeys := []string{"dceb1ce5778444e7777172e1f586488d2382fb1037887cd79a70b0cba4fb3dce", "9881aeb264452a4f7fafa1cc7bc4b93a05c55537c0703453e585f6d83ce77dca"}
	invalidKeys := []string{"", "cat", "deadbeef", "FA3eNd17NgaXZA3rXQVvzSvWHrpXfHWPzLQjJy2PQVQSc4ZutjC1", "FA38F8fY6duMqDLyCNUYWemdFSWgXDSteeNvNCmJ1Eyb86Z3VNZo"}

	for _, v := range validKeys {
		err := Wallet.ValidateKey(v)
		if err != nil {
			t.Errorf("ValidateKey returned error for valid key %v - %v\n", v, err)
		}
	}
	for _, v := range invalidKeys {
		err := Wallet.ValidateKey(v)
		if err == nil {
			t.Errorf("ValidateKey did not return error for invalid key `%v`\n", v)
		}
	}
}
