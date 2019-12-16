package data_resources

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/resource"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type EditCandidate struct {
	PubKey        string `json:"pub_key"`
	RewardAddress string `json:"reward_address"`
	OwnerAddress  string `json:"owner_address"`
}

func (EditCandidate) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.EditCandidateTxData)

	return EditCandidate{
		PubKey:        data.PubKey,
		RewardAddress: data.RewardAddress,
		OwnerAddress:  data.OwnerAddress,
	}
}
