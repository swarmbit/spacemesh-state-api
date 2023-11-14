package price

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const priceKey = "priceKey"

type PriceResolver struct {
	priceMap *sync.Map
}

func NewPriceResolver() *PriceResolver {
	priceResolver := &PriceResolver{
		priceMap: &sync.Map{},
	}
	priceResolver.fetchPrice()
	priceResolver.periodicPriceFetch()
	return priceResolver
}

func (p *PriceResolver) GetPrice() float64 {
	priceResponse, present := p.priceMap.Load(priceKey)
	if !present {
		return -1
	}
	value := priceResponse.(*PriceResponse).Quotes["USD"]
	if value != nil {
		return value.Price
	}
	return -1
}

func (p *PriceResolver) periodicPriceFetch() {
	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for range ticker.C {
			p.fetchPrice()
		}
	}()
}

func (p *PriceResolver) fetchPrice() {
	resp, err := http.Get("https://api.coinpaprika.com/v1/tickers/smh-spacemesh")
	if err != nil {
		fmt.Println("Error:", err)
		p.priceMap.Delete(priceKey)
		return
	}
	defer resp.Body.Close()

	var response PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		p.priceMap.Delete(priceKey)
		fmt.Println("Error:", err)
		return
	}

	p.priceMap.Store(priceKey, &response)

}

type PriceResponse struct {
	Quotes map[string]*PriceQuote `json:"quotes"`
}

type PriceQuote struct {
	Price float64 `json:"price"`
}
