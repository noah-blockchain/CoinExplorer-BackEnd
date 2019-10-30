package coins

import (
	"github.com/gin-gonic/gin"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/coins"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/core"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/errors"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/tools"
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"net/http"
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
	if req.Filter != nil && isModelsContain(*req.Filter, []string{"crr", "volume", "reserve_balance", "symbol", "price"}) {
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

	c.JSON(http.StatusOK, gin.H{"data": resource.TransformPaginatedCollection(data, coins.Resource{}, pagination)})
}
