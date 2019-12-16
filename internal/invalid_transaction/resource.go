package invalid_transaction

import (
	"time"

	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
)

type Resource struct {
	Hash      string `json:"hash"`
	Block     uint64 `json:"block"`
	Timestamp string `json:"timestamp"`
	Type      uint8  `json:"type"`
	From      string `json:"from"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	tx := model.(models.InvalidTransaction)

	return Resource{
		Hash:      tx.GetHash(),
		Block:     tx.BlockID,
		Timestamp: tx.CreatedAt.Format(time.RFC3339),
		Type:      tx.Type,
		From:      tx.FromAddress.GetAddress(),
	}
}
