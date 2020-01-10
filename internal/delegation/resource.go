package delegation

import (
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
)

type Resource struct {
	Coin           string             `json:"coin"`
	Value          string             `json:"value"`
	NoahValue      string             `json:"noah_value"`
	PubKey         string             `json:"pub_key"`
	ProfitReceived string             `json:"profit_received"`
	ValidatorMeta  resource.Interface `json:"validator_meta"`
}

func (resource Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	return model.(Resource)
}
