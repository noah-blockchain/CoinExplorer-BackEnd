package validator

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/stake"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/validator/meta"
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"gopkg.in/guregu/null.v3/zero"
)

type Resource struct {
	PublicKey       string                `json:"public_key"`
	Status          *uint8                `json:"status"`
	Meta            resource.Interface    `json:"meta"`
	Stake           *string               `json:"stake"`
	Commission      uint64                `json:"commission"`
	Part            *string               `json:"part"`
	Uptime          *float64              `json:"uptime"`
	CountDelegators *uint64               `json:"count_delegators"`
	DelegatorCount  *int                  `json:"delegator_count,omitempty"`
	DelegatorList   *[]resource.Interface `json:"delegator_list,omitempty"`
}

type Params struct {
	TotalStake          string // total stake of current active validator ids (by last block)
	ActiveValidatorsIDs []uint64
}

// Required extra params: object type of Params.
func (r Resource) Transform(model resource.ItemInterface, values ...resource.ParamInterface) resource.Interface {
	validator := model.(models.Validator)
	params := values[0].(Params)
	part, validatorStake := GetValidatorPartAndStake(validator, params.TotalStake, params.ActiveValidatorsIDs)

	res := Resource{
		PublicKey:       validator.GetPublicKey(),
		Status:          validator.Status,
		Stake:           validatorStake,
		Part:            part,
		Uptime:          validator.Uptime,
		Meta:            new(meta.Resource).Transform(validator),
		CountDelegators: validator.CountDelegators,
	}

	if validator.Commission != nil {
		res.Commission = *validator.Commission
	}

	return res
}

// return validator stake and part of the total (%)
func GetValidatorPartAndStake(validator models.Validator, totalStake string, validators []uint64) (*string, *string) {
	var part, stakeFull *string

	if helpers.InArray(validator.ID, validators) && validator.TotalStake != nil {
		val := helpers.CalculatePercent(*validator.TotalStake, totalStake)
		part = &val
	}

	if validator.TotalStake != nil {
		val := helpers.QNoahStr2Noah(*validator.TotalStake)
		stakeFull = &val
	}

	return part, stakeFull
}

// return list of delegators and count
func (r Resource) getDelegatorsListAndCount(validator models.Validator) (*[]resource.Interface, *int) {
	delegatorsCount := len(validator.Stakes)
	delegators := resource.TransformCollection(validator.Stakes, stake.Resource{})

	return &delegators, &delegatorsCount
}

type ResourceWithValidators struct {
	PublicKey   string `json:"public_key"`
	Name        string `json:"name"`
	SiteUrl     string `json:"site_url"`
	IconUrl     string `json:"icon_url"`
	Description string `json:"description"`
}

func (ResourceWithValidators) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	validator := model.(models.Validator)

	return ResourceWithValidators{
		PublicKey:   validator.GetPublicKey(),
		Name:        zero.StringFromPtr(validator.Name).String,
		SiteUrl:     zero.StringFromPtr(validator.SiteUrl).String,
		IconUrl:     zero.StringFromPtr(validator.IconUrl).String,
		Description: zero.StringFromPtr(validator.Description).String,
	}
}

type ResourceAggregator struct {
	PublicKey       string             `json:"public_key"`
	Stake           *string            `json:"stake"`
	Part            *string            `json:"part"`
	Uptime          *float64           `json:"uptime"`
	Commission      uint64             `json:"commission"`
	Status          *uint8             `json:"status"`
	CreatedAt       string             `json:"created_at"`
	CountDelegators *uint64            `json:"count_delegators"`
	Meta            resource.Interface `json:"meta"`
}

func (ResourceAggregator) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	validator := model.(ResourceAggregator)
	return validator
}
