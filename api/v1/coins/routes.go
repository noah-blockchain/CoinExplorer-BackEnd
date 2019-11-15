package coins

import (
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	coins := r.Group("/coins")
	{
		coins.GET("", GetCoins)
		coins.GET("/:symbol", GetCoinBySymbol)
		coins.GET("/:symbol/transactions", GetTransactions)
		coins.GET("/:symbol/validators", GetValidators)
	}
}
