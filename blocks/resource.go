package blocks

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	validatorMeta "github.com/noah-blockchain/CoinExplorer-BackEnd/validator/meta"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"time"
)

type Resource struct {
	ID          uint64               `json:"height"`
	Size        uint64               `json:"size"`
	NumTxs      uint32               `json:"txCount"`
	BlockTime   float64              `json:"blockTime"`
	Timestamp   string               `json:"timestamp"`
	BlockReward string               `json:"reward"`
	Hash        string               `json:"hash"`
	Validators  []resource.Interface `json:"validators"`
}

// lastBlockId - uint64 pointer to the last block height, optional field.
func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	block := model.(models.Block)

	return Resource{
		ID:          block.ID,
		Size:        block.Size,
		NumTxs:      block.NumTxs,
		BlockTime:   helpers.Nano2Seconds(block.BlockTime),
		Timestamp:   block.CreatedAt.Format(time.RFC3339),
		BlockReward: helpers.QNoahStr2Noah(block.BlockReward),
		Hash:        block.GetHash(),
		Validators:  resource.TransformCollection(block.BlockValidators, ValidatorResource{}),
	}
}

type ValidatorResource struct {
	PublicKey     string             `json:"publicKey"`
	ValidatorMeta resource.Interface `json:"validator_meta"`
	Signed        bool               `json:"signed"`
}

func (ValidatorResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	blockValidator := model.(models.BlockValidator)

	return ValidatorResource{
		PublicKey:     blockValidator.Validator.GetPublicKey(),
		Signed:        blockValidator.Signed,
		ValidatorMeta: new(validatorMeta.Resource).Transform(blockValidator.Validator),
	}
}
