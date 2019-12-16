package blocks

import (
	"time"

	"github.com/noah-blockchain/noah-explorer-extender/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-extender/internal/resource"
	validatorMeta "github.com/noah-blockchain/noah-explorer-extender/internal/validator/meta"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Resource struct {
	ID          uint64               `json:"height"`
	Size        uint64               `json:"size"`
	NumTxs      uint32               `json:"txCount"`
	BlockTime   float64              `json:"blockTime"`
	CreatedAt   string               `json:"created_at"`
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
		CreatedAt:   block.CreatedAt.Format(time.RFC3339),
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
