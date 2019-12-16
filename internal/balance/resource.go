package balance

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
)

type Resource struct {
	Coin   string `json:"coin"`
	Amount string `json:"amount"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	balance := model.(models.Balance)

	return Resource{
		Coin:   balance.Coin.Symbol,
		Amount: helpers.QNoahStr2Noah(balance.Value),
	}
}

type ResourceCoinAddressBalances struct {
	Address string `json:"coin"`
	Amount  string `json:"amount"`
}

func (ResourceCoinAddressBalances) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	balance := model.(models.Balance)

	return ResourceCoinAddressBalances{
		Address: balance.Address.GetAddress(),
		Amount:  helpers.QNoahStr2Noah(balance.Value),
	}
}
