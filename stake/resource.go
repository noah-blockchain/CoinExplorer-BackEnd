package stake

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type Resource struct {
	Coin      string `json:"coin"`
	Address   string `json:"address"`
	Value     string `json:"value"`
	NoahValue string `json:"noah_value"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:      stake.Coin.Symbol,
		Address:   stake.OwnerAddress.GetAddress(),
		Value:     helpers.QNoahStr2Noah(stake.Value),
		NoahValue: helpers.QNoahStr2Noah(stake.NoahValue),
	}
}
