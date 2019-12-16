package reward

import (
	"time"

	"github.com/noah-blockchain/noah-explorer-extender/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-extender/internal/resource"
	validatorMeta "github.com/noah-blockchain/noah-explorer-extender/internal/validator/meta"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Resource struct {
	BlockID       uint64             `json:"block"`
	Role          string             `json:"role"`
	Amount        string             `json:"amount"`
	Address       string             `json:"address"`
	Validator     string             `json:"validator"`
	ValidatorMeta resource.Interface `json:"validator_meta"`
	Timestamp     string             `json:"timestamp"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	reward := model.(models.Reward)

	return Resource{
		BlockID:       reward.BlockID,
		Role:          reward.Role,
		Amount:        helpers.QNoahStr2Noah(reward.Amount),
		Address:       reward.Address.GetAddress(),
		Validator:     reward.Validator.GetPublicKey(),
		Timestamp:     reward.Block.CreatedAt.Format(time.RFC3339),
		ValidatorMeta: new(validatorMeta.Resource).Transform(*reward.Validator),
	}
}
