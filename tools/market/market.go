package market

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/noahdev"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/core/config"
	"time"
)

type Service struct {
	PriceChange PriceChange
	api         *noahdev.Api
	baseCoin    string
}

type PriceChange struct {
	Price  float64
	Change float64
}

func NewService(noahdevApi *noahdev.Api, basecoin string) *Service {
	return &Service{
		api:         noahdevApi,
		baseCoin:    basecoin,
		PriceChange: PriceChange{Price: 0, Change: 0},
	}
}

func (s *Service) Run() {
	for {
		response, err := s.api.GetCurrentPrice()
		if err == nil {
			s.PriceChange = PriceChange{
				Price:  response.Data.Price / 10000,
				Change: response.Data.Delta,
			}
		}

		time.Sleep(time.Duration(config.MarketPriceUpdatePeriodInMin * time.Minute))
	}
}
