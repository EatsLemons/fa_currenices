package currency

import (
	"log"
	"time"

	"github.com/EatsLemons/fa_currencies/store"
)

type storage interface {
	Save([]store.Ratio) error
	GetCurrPair(from, to string) (*store.Ratio, error)
}

type cryptoStock interface {
	Prices(cryptoCurr []string, fiatCurr []string) ([]store.Ratio, error)
	CoinsList() ([]string, error)
}

var fiatCurrencies = []string{"USD", "RUB", "EUR", "GBP"}

type CurrencyService struct {
	Storage     storage
	cryptoStock cryptoStock
	reloadTime  int
}

func NewCurrencyService(s storage, c cryptoStock, reloadTime int) *CurrencyService {
	cs := CurrencyService{
		Storage:     s,
		cryptoStock: c,
		reloadTime:  reloadTime,
	}

	return &cs
}

func (cs *CurrencyService) Run() {
	for {
		dataToSave := cs.getCurrencyData()
		cs.saveCurrencyData(dataToSave)

		time.Sleep(time.Second * time.Duration(cs.reloadTime))
	}
}

func (cs *CurrencyService) saveCurrencyData(rates []store.Ratio) {
	err := cs.Storage.Save(rates)
	if err != nil {
		log.Printf("[WARN] fail to save new rates %s", err)
	}
}

func (cs *CurrencyService) getCurrencyData() []store.Ratio {
	crCoins, err := cs.cryptoStock.CoinsList()
	if err != nil {
		log.Println("[WARN] fail while reqeusts coins list %s", err)
		return nil
	}

	rates, priceErr := cs.cryptoStock.Prices(crCoins, fiatCurrencies)
	if priceErr != nil {
		log.Println("[WARN] fail while reqeusts currency rates %s", priceErr)
		return nil
	}

	return rates
}
