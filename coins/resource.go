package coins

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type Resource struct {
	Crr            uint64 `json:"crr"`
	Volume         string `json:"volume"`
	ReserveBalance string `json:"reserveBalance"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	coin := model.(models.Coin)
	return Resource{
		Crr:            coin.Crr,
		Volume:         helpers.QNoahStr2Noah(coin.Volume),
		ReserveBalance: helpers.QNoahStr2Noah(coin.ReserveBalance),
		Name:           coin.Name,
		Symbol:         coin.Symbol,
	}
}
