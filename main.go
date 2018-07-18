package main

import (
	"log"
	"os"

	"github.com/EatsLemons/fa_currencies/currency"
	"github.com/EatsLemons/fa_currencies/currency/crypto"
	"github.com/EatsLemons/fa_currencies/currency/storage"
	"github.com/EatsLemons/fa_currencies/rest"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Port int `long:"port" env:"FA_CURRENCIES_LIST_PORT" default:"8080" description:"port"`

	CurrCacheTurnOff bool `long:"cache-reload" env:"CACHE_RELOAD" description:"Make it true, if you don't need to update currency cache by this instance"`
	CurrCacheReload  int  `long:"cache-reload-time" env:"CACHE_RELOAD_TIME" default:"300" description:"time between cache reloads"`

	MongoDBHost     string `long:"mongodb-host" env:"MONGODB_HOST" default:"localhost:27017" description:"MongoDB host"`
	MongoDBLogin    string `long:"mongodb-login" env:"MONGODB_LOGIN" default:"" description:"MongoDB login"`
	MongoDBPassword string `long:"mongodb-password" env:"MONGODB_PASSWORD" default:"" description:"MongoDB password"`
	MongoDBBaseName string `long:"mongodb-base-name" env:"MONGODB_BASE_NAME" default:"fa_test" description:"MongoDB base name"`

	FiatCurrencies []string `long:"fiat-currencies" env:"FIAT_CURRENCY" default:"USD" default:"EUR" default:"GBR" default:"RUB" description:"Set the fiat currencies"`
}

func main() {
	p := flags.NewParser(&opts, flags.Default)
	if _, e := p.ParseArgs(os.Args[1:]); e != nil {
		log.Println(e.Error())
		os.Exit(1)
	}

	log.Println("Started with:")
	log.Printf("%+v", opts)

	ccAPI := crypto.NewCryptoCompareAPIClient()
	mongo := storage.NewMongoDB(opts.MongoDBHost, opts.MongoDBLogin, opts.MongoDBPassword, opts.MongoDBBaseName)

	currencyService := currency.NewCurrencyService(mongo, ccAPI, opts.CurrCacheReload, opts.FiatCurrencies)

	if !opts.CurrCacheTurnOff {
		go currencyService.Run()
	}

	srv := rest.Rest{
		CurrService: currencyService,
	}

	srv.Run(8080)
}
