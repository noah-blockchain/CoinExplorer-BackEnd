package address

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/balance"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/resource"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Resource struct {
	Address  string               `json:"address"`
	Balances []resource.Interface `json:"balances"`
}

func (r Resource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	address := model.(models.Address)
	result := Resource{
		Address:  address.GetAddress(),
		Balances: resource.TransformCollection(address.Balances, balance.Resource{}),
	}

	return result
}
