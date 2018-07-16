package main

import (
	"log"

	"github.com/EatsLemons/fa_currencies/currency/crypto"
)

func main() {
	ccAPI := crypto.NewCryptoCompareAPIClient()

	currcies, _ := ccAPI.CoinsList()
	result, err := ccAPI.Prices(currcies, []string{"USD", "RUB"})
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(len(result))
}
