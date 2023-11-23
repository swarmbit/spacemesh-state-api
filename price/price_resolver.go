package price

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"strconv"
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
	if err != nil || resp.StatusCode == 404 {
		fmt.Println("Error calling coinpaprika, try XT exchange")
		respXT, err := http.Get("https://www.xt.com/sapi/v4/market/public/ticker/24h?symbol=smh_usdt")
		if err != nil {
			p.priceMap.Delete(priceKey)
			return
		}
		var xtResponce PriceXTResponse
		if err := json.NewDecoder(respXT.Body).Decode(&xtResponce); err != nil {
			p.priceMap.Delete(priceKey)
			fmt.Println("Error:", err)
			return
		}

		if len(xtResponce.Result) > 0 {
			price, err := strconv.ParseFloat(xtResponce.Result[0].Current, 64)
			if err != nil {
				p.priceMap.Delete(priceKey)
				fmt.Println("Error no price on XT response")
				return
			}
			fmt.Println("Price: ", price)
			quotes := make(map[string]*PriceQuote)
			quotes["USD"] = &PriceQuote{
				Price: price,
			}
			p.priceMap.Store(priceKey, &PriceResponse{
				Quotes: quotes,
			})
			return
		} else {
			p.priceMap.Delete(priceKey)
			fmt.Println("Error no price on XT response")
			return
		}

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

type PriceXTResponse struct {
	Result []PriceXTResult `json:"result"`
}

type PriceXTResult struct {
	Current string `json:"c"`
}
