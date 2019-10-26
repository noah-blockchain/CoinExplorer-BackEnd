package chart

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/reward"
	"time"
)

type RewardResource struct {
	Time   string `json:"time"`
	Amount string `json:"amount"`
}

func (RewardResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := model.(reward.ChartData)

	return RewardResource{
		Time:   data.Time.Format(time.RFC3339),
		Amount: helpers.QNoahStr2Noah(data.Amount),
	}
}
