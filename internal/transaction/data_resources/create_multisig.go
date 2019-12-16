package data_resources

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/resource"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type CreateMultisig struct {
	Threshold string   `json:"threshold"`
	Weights   []string `json:"weights"`
	Addresses []string `json:"addresses"`
}

func (CreateMultisig) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.CreateMultisigTxData)

	return CreateMultisig{
		Threshold: data.Threshold,
		Weights:   data.Weights,
		Addresses: data.Addresses,
	}
}
