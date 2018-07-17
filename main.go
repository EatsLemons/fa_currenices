package main

import (
	"log"

	"github.com/EatsLemons/fa_currencies/currency/crypto"
	"github.com/EatsLemons/fa_currencies/currency/storage"
)

func main() {
	ccAPI := crypto.NewCryptoCompareAPIClient()

	currcies, _ := ccAPI.CoinsList()
	result, _ := ccAPI.Prices(currcies, []string{"USD", "RUB"})

	log.Println(len(result))

	mongo := storage.NewMongoDB("localhost:27017", "", "", "fa_test")
	res := mongo.Update(result)
	log.Println(res)
}
