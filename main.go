package main

import (
	"github.com/EatsLemons/fa_currencies/currency"
	"github.com/EatsLemons/fa_currencies/currency/crypto"
	"github.com/EatsLemons/fa_currencies/currency/storage"
	"github.com/EatsLemons/fa_currencies/rest"
)

func main() {
	ccAPI := crypto.NewCryptoCompareAPIClient()
	mongo := storage.NewMongoDB("localhost:27017", "", "", "fa_test")

	currencyService := currency.NewCurrencyService(mongo, ccAPI, 60)

	go currencyService.Run()

	srv := rest.Rest{
		CurrService: currencyService,
	}

	srv.Run(8080)
}
