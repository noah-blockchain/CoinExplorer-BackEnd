package validators

import (
	"net/http"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	"github.com/noah-blockchain/noah-explorer-extender/internal/core"
	"github.com/noah-blockchain/noah-explorer-extender/internal/errors"
	h "github.com/noah-blockchain/noah-explorer-extender/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-extender/internal/resource"
	"github.com/noah-blockchain/noah-explorer-extender/internal/stake"
	"github.com/noah-blockchain/noah-explorer-extender/internal/tools"
	"github.com/noah-blockchain/noah-explorer-extender/internal/transaction"
	"github.com/noah-blockchain/noah-explorer-extender/internal/validator"
	"github.com/noah-blockchain/noah-explorer-extender/internal/validator/meta"
	"github.com/noah-blockchain/coinExplorer-tools/helpers"
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"gopkg.in/guregu/null.v3/zero"
)

type GetAggregatedValidatorRequest struct {
	Page    string  `form:"page"     binding:"omitempty,numeric"`
	Filter  *string `form:"filter"   binding:"omitempty"`
	OrderBy *string `form:"order_by" binding:"omitempty"`
}

type GetValidatorRequest struct {
	PublicKey string `uri:"publicKey"    binding:"required,noahPubKey"`
}

// TODO: replace string to int
type GetValidatorTransactionsRequest struct {
	Page       string  `form:"page"        binding:"omitempty,numeric"`
	StartBlock *string `form:"startblock"  binding:"omitempty,numeric"`
	EndBlock   *string `form:"endblock"    binding:"omitempty,numeric"`
}

type CacheValidatorsData struct {
	Validators []models.Validator
	Pagination tools.Pagination
}

// cache time
const CacheBlocksCount = time.Duration(15)

// Get list of transaction by validator public key
func GetValidatorTransactions(c *gin.Context) {
	var validatorRequest GetValidatorRequest
	var request GetValidatorTransactionsRequest

	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	err := c.ShouldBindUri(&validatorRequest)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	err = c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	publicKey := helpers.RemovePrefix(validatorRequest.PublicKey)
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByFilter(transaction.ValidatorFilter{
		ValidatorPubKey: publicKey,
		StartBlock:      request.StartBlock,
		EndBlock:        request.EndBlock,
	}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get validator detail by public key
func GetValidator(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetValidatorRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	data := explorer.ValidatorRepository.GetByPublicKey(helpers.RemovePrefix(request.PublicKey))

	// check validator to existing
	if data == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Validator not found.", c)
		return
	}

	// get array of active validator ids by last block
	activeValidatorIDs := getActiveValidatorIDs(explorer)
	// get total stake of active validators
	totalStake := getTotalStakeByActiveValidators(explorer, activeValidatorIDs)

	c.JSON(http.StatusOK, gin.H{
		"data": validator.Resource{}.Transform(*data, validator.Params{
			TotalStake:          totalStake,
			ActiveValidatorsIDs: activeValidatorIDs,
		}),
	})
}

// Get IDs of active validators
func getActiveValidatorIDs(explorer *core.Explorer) []uint64 {
	return explorer.Cache.Get("active_validators", func() interface{} {
		return explorer.ValidatorRepository.GetActiveValidatorIds()
	}, CacheBlocksCount).([]uint64)
}

// Get total stake of active validators
func getTotalStakeByActiveValidators(explorer *core.Explorer, validators []uint64) string {
	return explorer.Cache.Get("validators_total_stake", func() interface{} {
		return explorer.ValidatorRepository.GetTotalStakeByActiveValidators(validators)
	}, CacheBlocksCount).(string)
}

func getValidatorsWithPagination(c *gin.Context, req GetAggregatedValidatorRequest, pagination *tools.Pagination) []models.Validator {
	explorer := c.MustGet("explorer").(*core.Explorer)
	var data []models.Validator

	var field, orderBy *string
	if req.Filter != nil && h.IsModelsContain(*req.Filter, []string{
		"uptime", "total_stake", "commission", "count_delegators"}) {
		field = req.Filter
	}

	if req.OrderBy != nil && h.IsModelsContain(*req.OrderBy, []string{"ASC", "DESC"}) {
		orderBy = req.OrderBy
	}

	getValidators := func() []models.Validator {
		return explorer.ValidatorRepository.GetValidatorsWithPagination(pagination, field, orderBy)
	}

	// cache last blocks
	if pagination.GetCurrentPage() == 1 && pagination.GetPerPage() == tools.DefaultLimit {
		//cached := explorer.Cache.Get("coins", func() interface{} {
		//	return CacheCoinsData{getCoins(), pagination}
		//}, CacheCoinsCount).(CacheCoinsData)
		cached := CacheValidatorsData{getValidators(), *pagination}
		data = cached.Validators
		*pagination = cached.Pagination
	} else {
		data = getValidators()
	}

	return data
}

func GetAggregatedValidators(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	var request GetAggregatedValidatorRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	activeValidatorIDs := getActiveValidatorIDs(explorer)
	totalStakeActiveValidators := getTotalStakeByActiveValidators(explorer, activeValidatorIDs)

	pagination := tools.NewPagination(c.Request)
	data := getValidatorsWithPagination(c, request, &pagination)

	resources := make([]validator.ResourceAggregator, len(data))
	for i, d := range data {
		resources[i] = validator.ResourceAggregator{
			PublicKey:       d.GetPublicKey(),
			Status:          d.Status,
			Meta:            new(meta.Resource).Transform(d),
			Uptime:          d.Uptime,
			CreatedAt:       d.CreatedAt.Format(time.RFC3339),
			CountDelegators: d.CountDelegators,
		}

		if d.Commission != nil {
			resources[i].Commission = *d.Commission
		}

		if d.TotalStake != nil {
			resources[i].Stake = pointer.ToString(h.QNoahStr2Noah(zero.StringFromPtr(d.TotalStake).String))
		}

		part, _ := validator.GetValidatorPartAndStake(d, totalStakeActiveValidators, activeValidatorIDs)
		resources[i].Part = part
	}

	// add params to each model resource
	c.JSON(http.StatusOK,
		resource.TransformPaginatedCollection(resources, validator.ResourceAggregator{}, pagination),
	)
}

// Get validator detail by public key
func GetDelegators(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	// validate request
	var request GetValidatorRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	data := explorer.StakeRepository.GetPaginatedDelegatorsForValidator(helpers.RemovePrefix(request.PublicKey), &pagination)

	// check validator to existing
	if data == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Validator not found.", c)
		return
	}

	c.JSON(http.StatusOK,
		resource.TransformPaginatedCollection(data, stake.ResourceDelegatorsForValidator{}, pagination),
	)
}
