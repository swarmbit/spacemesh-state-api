package price

import (
	"encoding/json"
	"fmt"
	"github.com/swarmbit/spacemesh-state-api/config"
	"net/http"
	"strings"
	"sync"
	"time"
	"strconv"
)

const priceKey = "priceKey"

type PriceResolver struct {
	priceMap *sync.Map
	isXT     bool
}

func NewPriceResolver(config *config.Config) *PriceResolver {
	fetchTime := 15
	isXT := false
	if config.Price != nil {
		if config.Price.RefreshTime > 0 {
			fetchTime = config.Price.RefreshTime
		}
		if strings.ToLower(config.Price.Provider) == "xt" {
			isXT = true
		}
	}
	priceResolver := &PriceResolver{
		priceMap: &sync.Map{},
		isXT:     isXT,
	}

	priceResolver.fetchPrice()
	priceResolver.periodicPriceFetch(fetchTime)
	return priceResolver
}

func (p *PriceResolver) GetPrice() float64 {
	priceResponse, present := p.priceMap.Load(priceKey)
	if !present {
		return -1
	}
	return priceResponse.(*PriceCache).usdPrice
}

func (p *PriceResolver) periodicPriceFetch(refreshTime int) {
	ticker := time.NewTicker(time.Duration(refreshTime) * time.Minute)
	go func() {
		for range ticker.C {
			p.fetchPrice()
		}
	}()
}

func (p *PriceResolver) fetchPrice() {
	if p.isXT {
		if !p.fetchXT() {
			p.fetchCoinpaprikaPrice()
		}
	} else {
		if !p.fetchCoinpaprikaPrice() {
			p.fetchXT()
		}
	}

}

func (p *PriceResolver) fetchXT() bool {
	fmt.Println("Fetch price from XT")
	respXT, err := http.Get("https://www.xt.com/sapi/v4/market/public/ticker/24h?symbol=smh_usdt")
	if err != nil || respXT.StatusCode == 404 {
		return false
	}
	defer respXT.Body.Close()

	var xtResponce PriceXTResponse
	if err := json.NewDecoder(respXT.Body).Decode(&xtResponce); err != nil {
		p.priceMap.Delete(priceKey)
		fmt.Println("Error:", err)
		return false
	}

	if len(xtResponce.Result) > 0 {
		price, err := strconv.ParseFloat(xtResponce.Result[0].Current, 64)
		if err != nil {
			p.priceMap.Delete(priceKey)

			fmt.Println("Error no price on XT response")
			return false
		}
		p.priceMap.Store(priceKey, &PriceCache{
			usdPrice: price,
		})
		return true
	} else {
		p.priceMap.Delete(priceKey)
		fmt.Println("Error no price on XT response")

		return false
	}

}

func (p *PriceResolver) fetchCoinpaprikaPrice() bool {
	fmt.Println("Fetch price from coinpaprika")
	resp, err := http.Get("https://api.coinpaprika.com/v1/tickers/smh-spacemesh")
	if err != nil || resp.StatusCode == 404 {
		return false
	}
	defer resp.Body.Close()

	var response PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		p.priceMap.Delete(priceKey)
		fmt.Println("Error:", err)
		return false
	}

	value := response.Quotes["USD"]
	if value == nil {
		return false
	}

	p.priceMap.Store(priceKey, &PriceCache{
		usdPrice: value.Price,
	})
	return true

}

type PriceCache struct {
	usdPrice float64
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
