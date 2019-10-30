package data_resources

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/coinExplorer-tools/models"
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
