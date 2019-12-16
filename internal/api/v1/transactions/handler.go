package transactions

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noah-blockchain/noah-explorer-extender/internal/blocks"
	"github.com/noah-blockchain/noah-explorer-extender/internal/core"
	"github.com/noah-blockchain/noah-explorer-extender/internal/errors"
	"github.com/noah-blockchain/noah-explorer-extender/internal/invalid_transaction"
	"github.com/noah-blockchain/noah-explorer-extender/internal/resource"
	"github.com/noah-blockchain/noah-explorer-extender/internal/tools"
	"github.com/noah-blockchain/noah-explorer-extender/internal/transaction"
	"github.com/noah-blockchain/coinExplorer-tools/helpers"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

// TODO: replace string in StartBlock, EndBlock, Page to int
type GetTransactionsRequest struct {
	Addresses  []string `form:"addresses[]" binding:"omitempty,noahAddress"`
	Page       string   `form:"page"        binding:"omitempty,numeric"`
	StartBlock *string  `form:"startblock"  binding:"omitempty,numeric"`
	EndBlock   *string  `form:"endblock"    binding:"omitempty,numeric"`
}

type GetTransactionRequest struct {
	Hash string `uri:"hash" binding:"noahTxHash"`
}

// Transaction cache helpers
const CacheBlocksCount = time.Duration(15)

type CacheTxData struct {
	Transactions []models.Transaction
	Pagination   tools.Pagination
}

// Get list of transactions
func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetTransactionsRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// remove Noah wallet prefix from each address
	noahAddresses := make([]string, len(request.Addresses))
	for key, addr := range request.Addresses {
		noahAddresses[key] = helpers.RemovePrefix(addr)
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)

	var txs []models.Transaction
	if len(noahAddresses) > 0 {
		txs = explorer.TransactionRepository.GetPaginatedTxsByAddresses(noahAddresses, transaction.BlocksRangeSelectFilter{
			StartBlock: request.StartBlock,
			EndBlock:   request.EndBlock,
		}, &pagination)
	} else {
		// prepare retrieving models
		getTxsFunc := func() []models.Transaction {
			return explorer.TransactionRepository.GetPaginatedTxsByFilter(blocks.RangeSelectFilter{
				StartBlock: request.StartBlock,
				EndBlock:   request.EndBlock,
			}, &pagination)
		}

		// cache last transactions
		if len(c.Request.URL.Query()) == 0 {
			cached := explorer.Cache.Get("transactions", func() interface{} {
				return CacheTxData{getTxsFunc(), pagination}
			}, CacheBlocksCount).(CacheTxData)

			txs = cached.Transactions
			pagination = cached.Pagination
		} else {
			txs = getTxsFunc()
		}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get transaction detail by hash
func GetTransaction(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetTransactionRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	hash := helpers.RemovePrefix(request.Hash)
	tx := explorer.TransactionRepository.GetTxByHash(hash)
	if tx == nil {
		invalidTx := explorer.InvalidTransactionRepository.GetTxByHash(hash)
		if invalidTx == nil {
			errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Transaction not found.", c)
			return
		}

		c.JSON(http.StatusPartialContent, gin.H{
			"data": new(invalid_transaction.Resource).Transform(*invalidTx),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(transaction.Resource).Transform(*tx),
	})
}
