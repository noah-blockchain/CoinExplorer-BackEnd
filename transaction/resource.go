package transaction

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"time"

	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/transaction/data_resources"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Resource struct {
	Txn         uint64                 `json:"txn"`
	Hash        string                 `json:"hash"`
	Nonce       uint64                 `json:"nonce"`
	Block       uint64                 `json:"block"`
	Timestamp   string                 `json:"timestamp"`
	Fee         string                 `json:"fee"`
	Type        uint8                  `json:"type"`
	Payload     string                 `json:"payload"`
	From        string                 `json:"from"`
	Data        resource.ItemInterface `json:"data"`
	Gas         uint64                 `json:"gas"`
	GasPrice    uint64                 `json:"gas_price"`
	GasCoinName string                 `json:"gas_coin"`
	To          *string                `json:"to,omitempty"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	tx := model.(models.Transaction)

	res := Resource{
		Txn:       tx.ID,
		Hash:      tx.GetHash(),
		Nonce:     tx.Nonce,
		Block:     tx.BlockID,
		Timestamp: tx.CreatedAt.Format(time.RFC3339),
		Fee:       helpers.Fee2Noah(tx.GetFee()),
		Type:      tx.Type,
		Payload:   base64.StdEncoding.EncodeToString(tx.Payload[:]),
		From:      tx.FromAddress.GetAddress(),
		Data:      TransformTxData(tx),
		Gas:       tx.Gas,
		GasPrice:  tx.GasPrice,
	}

	if tx.GasCoin != nil {
		res.GasCoinName = tx.GasCoin.Symbol
	}

	if tx.Type == models.TxTypeSend {
		var sendTxData models.SendTxData
		if err := json.Unmarshal(tx.Data, &sendTxData); err == nil {
			res.To = &sendTxData.To
		}
	}

	return res
}

type TransformTxConfig struct {
	Model    resource.ItemInterface
	Resource resource.Interface
}

var transformConfig = map[uint8]TransformTxConfig{
	models.TxTypeSend:                {Model: new(models.SendTxData), Resource: data_resources.Send{}},
	models.TxTypeSellCoin:            {Model: new(models.SellCoinTxData), Resource: data_resources.SellCoin{}},
	models.TxTypeSellAllCoin:         {Model: new(models.SellAllCoinTxData), Resource: data_resources.SellAllCoin{}},
	models.TxTypeBuyCoin:             {Model: new(models.BuyCoinTxData), Resource: data_resources.BuyCoin{}},
	models.TxTypeCreateCoin:          {Model: new(models.CreateCoinTxData), Resource: data_resources.CreateCoin{}},
	models.TxTypeDeclareCandidacy:    {Model: new(models.DeclareCandidacyTxData), Resource: data_resources.DeclareCandidacy{}},
	models.TxTypeDelegate:            {Model: new(models.DelegateTxData), Resource: data_resources.Delegate{}},
	models.TxTypeUnbound:             {Model: new(models.UnbondTxData), Resource: data_resources.Unbond{}},
	models.TxTypeRedeemCheck:         {Model: new(models.RedeemCheckTxData), Resource: data_resources.RedeemCheck{}},
	models.TxTypeMultiSig:            {Model: new(models.CreateMultisigTxData), Resource: data_resources.CreateMultisig{}},
	models.TxTypeMultiSend:           {Model: new(models.MultiSendTxData), Resource: data_resources.Multisend{}},
	models.TxTypeEditCandidate:       {Model: new(models.EditCandidateTxData), Resource: data_resources.EditCandidate{}},
	models.TxTypeSetCandidateOnline:  {Model: new(models.SetCandidateTxData), Resource: data_resources.SetCandidate{}},
	models.TxTypeSetCandidateOffline: {Model: new(models.SetCandidateTxData), Resource: data_resources.SetCandidate{}},
}

func TransformTxData(tx models.Transaction) resource.Interface {
	config := transformConfig[tx.Type]

	val := reflect.New(reflect.TypeOf(config.Model).Elem()).Interface()
	err := json.Unmarshal(tx.Data, val)
	helpers.CheckErr(err)

	return config.Resource.Transform(val, tx)
}

type ResourceTransactionOutput struct {
	Txn         uint64  `json:"txn"`
	Hash        string  `json:"hash"`
	Nonce       uint64  `json:"nonce"`
	Block       uint64  `json:"block"`
	Timestamp   string  `json:"timestamp"`
	Fee         string  `json:"fee"`
	Type        uint8   `json:"type"`
	From        string  `json:"from"`
	Gas         uint64  `json:"gas"`
	GasPrice    uint64  `json:"gas_price"`
	GasCoinName string  `json:"gas_coin"`
	To          *string `json:"to,omitempty"`
	//Data        resource.ItemInterface `json:"data"`
}

func (ResourceTransactionOutput) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	txOutput := model.(models.TransactionOutput)

	res := ResourceTransactionOutput{
		Txn:       txOutput.ID,
		Hash:      txOutput.Transaction.GetHash(),
		Nonce:     txOutput.Transaction.Nonce,
		Block:     txOutput.Transaction.BlockID,
		Timestamp: txOutput.Transaction.CreatedAt.Format(time.RFC3339),
		Fee:       helpers.Fee2Noah(txOutput.Transaction.GetFee()),
		Type:      txOutput.Transaction.Type,
		From:      txOutput.Transaction.FromAddress.GetAddress(),
		Gas:       txOutput.Transaction.Gas,
		GasPrice:  txOutput.Transaction.GasPrice,
	}

	//if txOutput.Transaction != nil {
	//	res.Data = TransformTxData(*txOutput.Transaction)
	//}

	if txOutput.Transaction.GasCoin != nil {
		res.GasCoinName = txOutput.Transaction.GasCoin.Symbol
	}

	if txOutput.Transaction.Type == models.TxTypeSend {
		var sendTxData models.SendTxData
		if err := json.Unmarshal(txOutput.Transaction.Data, &sendTxData); err == nil {
			res.To = &sendTxData.To
		}
	}

	return res
}
