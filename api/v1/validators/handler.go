package validators

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/core"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/errors"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/tools"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/transaction"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/validator"
	"github.com/noah-blockchain/noah-explorer-tools/helpers"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type GetValidatorRequest struct {
	PublicKey string `uri:"publicKey"    binding:"required,noahPubKey"`
}

// TODO: replace string to int
type GetValidatorTransactionsRequest struct {
	Page       string  `form:"page"        binding:"omitempty,numeric"`
	StartBlock *string `form:"startblock"  binding:"omitempty,numeric"`
	EndBlock   *string `form:"endblock"    binding:"omitempty,numeric"`
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
			TotalStake:           totalStake,
			ActiveValidatorsIDs:  activeValidatorIDs,
			IsDelegatorsRequired: true,
		}),
	})
}

// Get list of validators
func GetValidators(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// fetch validators
	validators := explorer.Cache.Get("validators", func() interface{} {
		return explorer.ValidatorRepository.GetValidators()
	}, CacheBlocksCount).([]models.Validator)

	// get array of active validator ids by last block
	activeValidatorIDs := getActiveValidatorIDs(explorer)
	// get total stake of active validators
	totalStake := getTotalStakeByActiveValidators(explorer, activeValidatorIDs)

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		return resource.ParamsInterface{validator.Params{
			TotalStake:           totalStake,
			ActiveValidatorsIDs:  activeValidatorIDs,
			IsDelegatorsRequired: false,
		}}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollectionWithCallback(
			validators,
			validator.Resource{},
			resourceCallback,
		),
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
