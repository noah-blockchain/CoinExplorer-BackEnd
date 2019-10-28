package coins

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"github.com/noah-blockchain/noah-go-node/core/types"
	"github.com/noah-blockchain/noah-go-node/math"
)

type Resource struct {
	Crr            uint64 `json:"crr"`
	Volume         string `json:"volume"`
	ReserveBalance string `json:"reserveBalance"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	Price          string `json:"price"`
	Timestamp      string `json:"timestamp"`
	Creator        string `json:"creator"`
}

const (
	precision = 100
)

func newFloat(x float64) *big.Float {
	return big.NewFloat(x).SetPrec(precision)
}

func convertStringToBigInt(value string) (*big.Int, error) {
	newValue := new(big.Int)
	newValue, ok := newValue.SetString(value, 10)
	if !ok {
		fmt.Println("SetString: error")
		return nil, errors.New("Can't convert string to big.Int (" + value + ")")
	}

	return newValue, nil
}

//reserve * (math.pow(1 + 1 / volume, 100 / crr) - 1)
func CalculatePurchaseAmount(supply *big.Int, reserve *big.Int, crr uint, wantReceive *big.Int) string {
	if wantReceive.Cmp(types.Big0) == 0 {
		return big.NewInt(0).String()
	}

	if crr == 100 {
		result := big.NewInt(0).Mul(wantReceive, reserve)
		return result.Div(result, supply).String()
	}

	tSupply := newFloat(0).SetInt(supply)
	tReserve := newFloat(0).SetInt(reserve)
	tWantReceive := newFloat(0).SetInt(wantReceive)

	res := newFloat(0).Add(tWantReceive, tSupply)   // reserve + supply
	res.Quo(res, tSupply)                           // (reserve + supply) / supply
	res = math.Pow(res, newFloat(100/float64(crr))) // ((reserve + supply) / supply)^(100/c)
	res.Sub(res, newFloat(1))                       // (((reserve + supply) / supply)^(100/c) - 1)
	res.Mul(res, tReserve)                          // reserve * (((reserve + supply) / supply)^(100/c) - 1)

	result, _ := res.Int(nil)

	return helpers.QNoahStr2Noah(result.String())
}

func getTokenPrice(volumeStr string, reserveStr string, crr uint64) string {
	volume, err := convertStringToBigInt(volumeStr)
	if err != nil {
		log.Println(err)
		return "0"
	}

	reserve, err := convertStringToBigInt(reserveStr)
	if err != nil {
		log.Println(err)
		return "0"
	}

	return CalculatePurchaseAmount(volume, reserve, uint(crr), big.NewInt(1))
}

func noahAddressFormatting(address string) string {
	if address == "" {
		return ""
	}

	return fmt.Sprintf("NOAHx%s", address)
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	coin := model.(models.Coin)

	res := Resource{
		Crr:            coin.Crr,
		Volume:         helpers.QNoahStr2Noah(coin.Volume),
		ReserveBalance: helpers.QNoahStr2Noah(coin.ReserveBalance),
		Name:           coin.Name,
		Symbol:         coin.Symbol,
		Price:          getTokenPrice(coin.Volume, coin.ReserveBalance, coin.Crr),
		Timestamp:      coin.UpdatedAt.Format(time.RFC3339),
		Creator:        noahAddressFormatting(coin.NoahAddress),
	}

	return res
}
