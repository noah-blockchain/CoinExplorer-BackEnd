package data_resources

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
)

type DeclareCandidacy struct {
	Address    string `json:"address"`
	PubKey     string `json:"pub_key"`
	Commission string `json:"commission"`
	Coin       string `json:"coin"`
	Stake      string `json:"stake"`
}

func (DeclareCandidacy) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.DeclareCandidacyTxData)

	return DeclareCandidacy{
		Address:    data.Address,
		PubKey:     data.PubKey,
		Commission: data.Commission,
		Coin:       data.Coin,
		Stake:      helpers.QNoahStr2Noah(data.Stake),
	}
}
