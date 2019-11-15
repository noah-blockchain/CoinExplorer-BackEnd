package coins

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/balance"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/transaction"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/coins"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/core"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/errors"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/tools"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

//const CacheCoinsCount = time.Duration(15)

//symbol, volume, reserve, crr, price
//ask, desc
type GetCoinsRequest struct {
	Page    string  `form:"page"     binding:"omitempty,numeric"`
	Symbol  *string `form:"symbol"   binding:"omitempty"`
	Filter  *string `form:"filter"   binding:"omitempty"`
	OrderBy *string `form:"order_by" binding:"omitempty"`
}

type GetCoinBySymbolRequest struct {
	Symbol string `uri:"symbol"`
}

type CacheCoinsData struct {
	Coins      []models.Coin
	Pagination tools.Pagination
}

func isModelsContain(value string, values []string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}

	return false
}

func getCoinsWithPagination(c *gin.Context, req GetCoinsRequest, pagination *tools.Pagination) []models.Coin {
	explorer := c.MustGet("explorer").(*core.Explorer)
	var data []models.Coin

	var field, orderBy *string
	if req.Filter != nil && isModelsContain(*req.Filter, []string{"crr", "volume", "reserve_balance", "symbol",
		"price", "capitalization", "delegated"}) {
		field = req.Filter
	}

	if req.OrderBy != nil && isModelsContain(*req.OrderBy, []string{"ASC", "DESC"}) {
		orderBy = req.OrderBy
	}

	getCoins := func() []models.Coin {
		return explorer.CoinRepository.GetPaginated(pagination, field, orderBy, req.Symbol)
	}

	// cache last blocks
	if pagination.GetCurrentPage() == 1 && pagination.GetPerPage() == tools.DefaultLimit {
		//cached := explorer.Cache.Get("coins", func() interface{} {
		//	return CacheCoinsData{getCoins(), pagination}
		//}, CacheCoinsCount).(CacheCoinsData)
		cached := CacheCoinsData{getCoins(), *pagination}
		data = cached.Coins
		*pagination = cached.Pagination
	} else {
		data = getCoins()
	}

	return data
}

// Get list of coins
func GetCoins(c *gin.Context) {
	//explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetCoinsRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var data []models.Coin
	pagination := tools.NewPagination(c.Request)
	data = getCoinsWithPagination(c, request, &pagination)

	// make response as empty array if no models found
	if len(data) == 0 {
		empty := make([]coins.Resource, 0)
		c.JSON(http.StatusOK, gin.H{"data": empty})
		return
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(data, coins.Resource{}, pagination))
}

// Get coin detail
func GetCoinBySymbol(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetCoinBySymbolRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch coin by symbol
	coin := explorer.CoinRepository.GetBySymbol(request.Symbol)

	// check coin to existing
	if coin == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Coin not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(coins.Resource).Transform(*coin),
	})
}

// Get list of transactions by noah address
func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetCoinBySymbolRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByCoin(request.Symbol, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.ResourceTransactionOutput{}, pagination))
}

// Get validator detail by public key
func GetValidators(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetCoinBySymbolRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	data := explorer.ValidatorRepository.GetValidatorsBySymbol(request.Symbol, &pagination)

	// check validator to existing
	if data == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Validator not found.", c)
		return
	}
	c.JSON(http.StatusOK,
		resource.TransformPaginatedCollection(data, validator.ResourceWithValidators{}, pagination),
	)
}

func GetAddressBalances(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetCoinBySymbolRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	balances := explorer.AddressRepository.GetBalancesByCoinSymbol(request.Symbol, &pagination)

	c.JSON(http.StatusOK,
		resource.TransformPaginatedCollection(balances, balance.ResourceCoinAddressBalances{}, pagination),
	)
}
