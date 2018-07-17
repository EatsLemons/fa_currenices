package crypto

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/EatsLemons/fa_currencies/store"
)

// CryptoCompareAPI represents client from https://min-api.cryptocompare.com/
type CryptoCompareAPI struct {
	address    string
	httpClient *http.Client
}

func NewCryptoCompareAPIClient() *CryptoCompareAPI {
	client := CryptoCompareAPI{
		address: "https://min-api.cryptocompare.com",
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	return &client
}

func (cc *CryptoCompareAPI) CoinsList() ([]string, error) {
	requestURI := "/data/all/coinlist"

	response := coinsListResponse{}
	err := cc.makeGetRequest(requestURI, &response, true)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(response.Data))

	for curr := range response.Data {
		result = append(result, curr)
	}

	return result, nil
}

func (cc *CryptoCompareAPI) Prices(cryptoCurrencies, fiatCurrencies []string) ([]store.Ratio, error) {
	// Making request URIs to satisfy provider condition:
	// fsyms Comma separated cryptocurrency symbols list [Max character length: 300]
	fsyms := cc.makeCurrencyRQStrings(cryptoCurrencies, 300)
	// tsyms Comma separated cryptocurrency symbols list to convert into [Max character length: 100]
	tsyms := cc.makeCurrencyRQStrings(fiatCurrencies, 100)
	result := make([]store.Ratio, 0, len(cryptoCurrencies))

	for _, fsym := range fsyms {
		for _, tsym := range tsyms {
			requestURI := "/data/pricemulti?fsyms=" + fsym + "&tsyms=" + tsym
			response := make(map[string]map[string]float64, 0)
			log.Println("req")
			err := cc.makeGetRequest(requestURI, &response, true)
			if err != nil {
				return nil, err
			}

			for curr, ratios := range response {
				ratio := store.Ratio{
					From: curr,
					To:   make(map[string]float64),
				}

				for fiatCurr, rate := range ratios {
					ratio.To[fiatCurr] = rate
				}

				result = append(result, ratio)
			}
		}
	}

	return result, nil
}

func (cc *CryptoCompareAPI) makeGetRequest(url string, result interface{}, checkRateLimit bool) error {
	if checkRateLimit {
		ok, restErr := cc.needToRestForRateLimit()
		if !ok || restErr != nil {
			time.Sleep(time.Second * 1)
		}
	}

	r, err := cc.httpClient.Get(cc.address + url)
	if err != nil {
		log.Printf("[WARN] request to cryptocompare has failed %s for query %s", err, url)
		return err
	}

	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(result)
}

func (cc *CryptoCompareAPI) needToRestForRateLimit() (bool, error) {
	requestURI := "/stats/rate/limit"

	response := statsRateLimitResponse{}
	err := cc.makeGetRequest(requestURI, &response, false)
	if err != nil {
		return true, err
	}

	if response.Second.CallsMade.Price >= 50 {
		return true, nil
	}

	return false, nil
}

func (cc *CryptoCompareAPI) makeCurrencyRQStrings(currencies []string, sLength int) []string {
	result := make([]string, 0)
	tmpRq := currencies[0]
	for i := 1; i < len(currencies); i++ {
		if utf8.RuneCountInString(tmpRq)+utf8.RuneCountInString(currencies[i])+1 < sLength {
			tmpRq = tmpRq + "," + currencies[i]

			if i == len(currencies)-1 {
				result = append(result, tmpRq)
			}

			continue
		}

		result = append(result, tmpRq)
		tmpRq = currencies[i]
	}

	return result
}

type priceReponse map[string]float64

type statsRateLimitResponse struct {
	Second struct {
		CallsMade struct {
			Price int
		}
	}
}

type coinsListResponse struct {
	Data map[string]interface{} `json:"Data"`
}
