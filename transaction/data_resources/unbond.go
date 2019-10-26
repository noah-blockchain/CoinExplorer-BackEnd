package data_resources

import (

	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type Unbond struct {
	PubKey string `json:"pub_key"`
	Coin   string `json:"coin"`
	Value  string `json:"value"`
}

func (Unbond) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.UnbondTxData)

	return Unbond{
		PubKey: data.PubKey,
		Coin:   data.Coin,
		Value:  helpers.QNoahStr2Noah(data.Value),
	}
}
