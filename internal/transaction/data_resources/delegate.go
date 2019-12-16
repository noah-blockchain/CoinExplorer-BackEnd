package data_resources

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
)

type Delegate struct {
	PubKey string `json:"pub_key"`
	Coin   string `json:"coin"`
	Value  string `json:"value"`
}

func (Delegate) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.DelegateTxData)

	return Delegate{
		PubKey: data.PubKey,
		Coin:   data.Coin,
		Value:  helpers.QNoahStr2Noah(data.Value),
	}
}
