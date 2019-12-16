package chart

import (
	"time"

	"github.com/noah-blockchain/noah-explorer-extender/internal/resource"
	"github.com/noah-blockchain/noah-explorer-extender/internal/transaction"
)

type TransactionResource struct {
	Date    string `json:"date"`
	TxCount uint64 `json:"txCount"`
}

func (TransactionResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := model.(transaction.TxCountChartData)

	return TransactionResource{
		Date:    data.Time.Format(time.RFC3339),
		TxCount: data.Count,
	}
}
