package data_resources

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
)

type Send struct {
	Coin  string `json:"coin"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func (Send) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.SendTxData)

	return Send{
		Coin:  data.Coin,
		To:    data.To,
		Value: helpers.QNoahStr2Noah(data.Value),
	}
}
