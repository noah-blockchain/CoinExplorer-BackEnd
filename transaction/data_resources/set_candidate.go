package data_resources

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type SetCandidate struct {
	PubKey string `json:"pub_key"`
}

func (SetCandidate) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.SetCandidateTxData)

	return SetCandidate{
		PubKey: data.PubKey,
	}
}
