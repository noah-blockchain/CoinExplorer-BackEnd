package validator

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/stake"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/validator/meta"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Resource struct {
	PublicKey      string                `json:"public_key"`
	Status         *uint8                `json:"status"`
	Meta           resource.Interface    `json:"meta"`
	Stake          *string               `json:"stake"`
	Part           *string               `json:"part"`
	DelegatorCount *int                  `json:"delegator_count,omitempty"`
	DelegatorList  *[]resource.Interface `json:"delegator_list,omitempty"`
}

type Params struct {
	TotalStake           string // total stake of current active validator ids (by last block)
	ActiveValidatorsIDs  []uint64
	IsDelegatorsRequired bool
}

// Required extra params: object type of Params.
func (r Resource) Transform(model resource.ItemInterface, values ...resource.ParamInterface) resource.Interface {
	validator := model.(models.Validator)
	params := values[0].(Params)
	part, validatorStake := r.getValidatorPartAndStake(validator, params.TotalStake, params.ActiveValidatorsIDs)

	result := Resource{
		PublicKey: validator.GetPublicKey(),
		Status:    validator.Status,
		Stake:     validatorStake,
		Part:      part,
		Meta:      new(meta.Resource).Transform(validator),
	}

	if params.IsDelegatorsRequired {
		result.DelegatorList, result.DelegatorCount = r.getDelegatorsListAndCount(validator)
	}

	return result
}

// return validator stake and part of the total (%)
func (r Resource) getValidatorPartAndStake(validator models.Validator, totalStake string, validators []uint64) (*string, *string) {
	var part, stake *string

	if helpers.InArray(validator.ID, validators) && validator.TotalStake != nil {
		val := helpers.CalculatePercent(*validator.TotalStake, totalStake)
		part = &val
	}

	if validator.TotalStake != nil {
		val := helpers.QNoahStr2Noah(*validator.TotalStake)
		stake = &val
	}

	return part, stake
}

// return list of delegators and count
func (r Resource) getDelegatorsListAndCount(validator models.Validator) (*[]resource.Interface, *int) {
	delegatorsCount := len(validator.Stakes)
	delegators := resource.TransformCollection(validator.Stakes, stake.Resource{})

	return &delegators, &delegatorsCount
}
