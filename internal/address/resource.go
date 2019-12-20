package address

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/balance"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
	"sort"
)

type Resource struct {
	Address  string               `json:"address"`
	Balances []resource.Interface `json:"balances"`
}

type ResourceTopAddresses struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type ByBalance []ResourceTopAddresses

func (a ByBalance) Len() int { return len(a) }

func (a ByBalance) Less(i, j int) bool {
	x, _ := helpers.NewFloat(0, 100).SetString(a[i].Balance)
	y, _ := helpers.NewFloat(0, 100).SetString(a[j].Balance)
	return x.Cmp(y) == 1
}

func (a ByBalance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (r Resource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	address := model.(models.Address)
	result := Resource{
		Address:  address.GetAddress(),
		Balances: resource.TransformCollection(address.Balances, balance.Resource{}),
	}

	return result
}

func (r ResourceTopAddresses) Transform(model []models.Address) []ResourceTopAddresses {
	top := make([]ResourceTopAddresses, len(model))
	for i, address := range model {
		balans := helpers.NewFloat(0, 100)
		for _, b := range address.Balances {
			if b.Coin.Symbol == "NOAH" {
				amount, _ := helpers.NewFloat(0, 100).SetString(b.Value)
				balans.Add(balans, amount)
			} else {
				price := helpers.GetPrice(b.Value, b.Coin.Price)
				balans.Add(balans, price)
			}
		}
		result := ResourceTopAddresses{
			Address: address.GetAddress(),
			Balance: balans.String(),
		}
		top[i] = result
	}
	sort.Sort(ByBalance(top))
	return top
}
