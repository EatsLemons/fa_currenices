package rest

type CurrenciesResponse struct {
	Result *Rates    `json:"result,omitempty"`
	Errors []ErrorRs `json:"errors,omitempty"`
}

type Rates struct {
	From string  `json:"from,omitempty"`
	To   string  `json:"to,omitempty"`
	Rate float64 `json:"rate,omitempty"`
}

type ErrorRs struct {
	Message string `json:"message,omitempty"`
}
