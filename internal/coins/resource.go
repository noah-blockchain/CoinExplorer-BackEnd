package coins

import (
	"time"

	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
)

type Resource struct {
	Crr                 uint64 `json:"crr"`
	Volume              string `json:"volume"`
	ReserveBalance      string `json:"reserve_balance"`
	Name                string `json:"name"`
	Symbol              string `json:"symbol"`
	Price               string `json:"price"`
	StartPrice          string `json:"start_price"`
	StartVolume         string `json:"start_volume"`
	StartReserveBalance string `json:"start_reserve_balance"`
	Capitalization      string `json:"capitalization"`
	Delegated           uint64 `json:"delegated"`
	CreatedAt           string `json:"created_at"`
	Creator             string `json:"creator"`
	Description         string `json:"description"`
	IconURL             string `json:"icon_url"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	coin := model.(models.Coin)

	res := Resource{
		Crr:                 coin.Crr,
		Volume:              helpers.QNoahStr2Noah(coin.Volume),
		ReserveBalance:      helpers.QNoahStr2Noah(coin.ReserveBalance),
		Price:               helpers.QNoahStr2Noah(coin.Price),
		Capitalization:      helpers.ConvertCapitalizationQNoahToNoah(coin.Capitalization),
		StartVolume:         helpers.QNoahStr2Noah(coin.StartVolume),
		StartReserveBalance: helpers.QNoahStr2Noah(coin.StartReserveBalance),
		StartPrice:          helpers.QNoahStr2Noah(coin.StartPrice),
		Name:                coin.Name,
		Symbol:              coin.Symbol,
		Delegated:           coin.Delegated,
		CreatedAt:           coin.CreatedAt.Format(time.RFC3339),
		Creator:             coin.GetAddress(),
		Description:         coin.Description,
		IconURL:             coin.IconURL,
	}

	return res
}
