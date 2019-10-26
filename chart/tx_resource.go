package chart

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/transaction"
	"time"
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
