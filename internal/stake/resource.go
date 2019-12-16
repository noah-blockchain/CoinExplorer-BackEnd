package stake

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
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

type ResourceStakeDelegation struct {
	Address   string `json:"address"`
	Value     string `json:"value"`
	NoahValue string `json:"noah_value"`
	PublicKey string `json:"public_key"`
}

func (ResourceStakeDelegation) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	stake := model.(models.Stake)

	return ResourceStakeDelegation{
		Address:   stake.OwnerAddress.GetAddress(),
		PublicKey: stake.Validator.GetPublicKey(),
		Value:     helpers.QNoahStr2Noah(stake.Value),
		NoahValue: helpers.QNoahStr2Noah(stake.NoahValue),
	}
}

type ResourceDelegatorsForValidator struct {
	Address   string `json:"address"`
	Symbol    string `json:"symbol"`
	Value     string `json:"value"`
	NoahValue string `json:"noah_value"`
}

func (ResourceDelegatorsForValidator) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	stake := model.(models.Stake)

	return ResourceDelegatorsForValidator{
		Address:   stake.OwnerAddress.GetAddress(),
		Symbol:    stake.Coin.Symbol,
		Value:     helpers.QNoahStr2Noah(stake.Value),
		NoahValue: helpers.QNoahStr2Noah(stake.NoahValue),
	}
}
