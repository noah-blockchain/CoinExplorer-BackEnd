package address

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/balance"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
	"github.com/noah-blockchain/noah-explorer-api/internal/tools"
	"sort"
)

const precision = 100

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
	x, _ := helpers.NewFloat(0, precision).SetString(a[i].Balance)
	y, _ := helpers.NewFloat(0, precision).SetString(a[j].Balance)
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

func (r ResourceTopAddresses) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	return nil
}

func (r ResourceTopAddresses) TransformCollection(model []models.Address, pagination tools.Pagination) resource.PaginationResource {
	top := make([]ResourceTopAddresses, len(model))
	for i, address := range model {
		uBalance := helpers.NewFloat(0, precision)
		for _, b := range address.Balances {
			if b.Coin.Symbol == "NOAH" {
				amount, _ := helpers.NewFloat(0, precision).SetString(b.Value)
				uBalance.Add(uBalance, amount)
			} else {
				price := helpers.GetPrice(b.Value, b.Coin.Price)
				uBalance.Add(uBalance, price)
			}
		}
		result := ResourceTopAddresses{
			Address: address.GetAddress(),
			Balance: uBalance.String(),
		}
		top[i] = result
	}

	sort.Sort(ByBalance(top))
	result := make([]resource.Interface, len(top))
	for i, v := range top {
		result[i] = v
	}

	return resource.PaginationResource{
		Data: result,
		Links: resource.PaginationLinksResource{
			First: pagination.GetFirstPageLink(),
			Last:  pagination.GetLastPageLink(),
			Prev:  pagination.GetPrevPageLink(),
			Next:  pagination.GetNextPageLink(),
		},
		Meta: resource.PaginationMetaResource{
			CurrentPage: pagination.GetCurrentPage(),
			LastPage:    pagination.GetLastPage(),
			Path:        pagination.GetPath(),
			PerPage:     pagination.GetPerPage(),
			Total:       pagination.Total,
			Additional:  nil,
		},
	}

}
