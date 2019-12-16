package delegation

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
	validatorMeta "github.com/noah-blockchain/noah-explorer-api/internal/validator/meta"
)

type Resource struct {
	Coin          string             `json:"coin"`
	Value         string             `json:"value"`
	NoahValue     string             `json:"noah_value"`
	PubKey        string             `json:"pub_key"`
	ValidatorMeta resource.Interface `json:"validator_meta"`
}

func (resource Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:          stake.Coin.Symbol,
		PubKey:        stake.Validator.GetPublicKey(),
		Value:         helpers.QNoahStr2Noah(stake.Value),
		NoahValue:     helpers.QNoahStr2Noah(stake.NoahValue),
		ValidatorMeta: new(validatorMeta.Resource).Transform(*stake.Validator),
	}
}
