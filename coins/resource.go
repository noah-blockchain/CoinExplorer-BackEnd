package coins

import (
	"fmt"
	"time"

	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Resource struct {
	Crr            uint64 `json:"crr"`
	Volume         string `json:"volume"`
	ReserveBalance string `json:"reserveBalance"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	Price          string `json:"price"`
	Capitalization string `json:"capitalization"`
	Delegated      uint64 `json:"delegated"`
	Timestamp      string `json:"timestamp"`
	Creator        string `json:"creator"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	coin := model.(models.Coin)

	res := Resource{
		Crr:            coin.Crr,
		Volume:         helpers.QNoahStr2Noah(coin.Volume),
		ReserveBalance: helpers.QNoahStr2Noah(coin.ReserveBalance),
		Name:           coin.Name,
		Symbol:         coin.Symbol,
		Price:          coin.Price,
		Delegated:      coin.Delegated,
		Timestamp:      coin.UpdatedAt.Format(time.RFC3339),
		Creator:        fmt.Sprintf("NOAHx%s", coin.Address),
		Capitalization: coin.Capitalization,
	}

	return res
}
