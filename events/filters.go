package events

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/blocks"
	"github.com/go-pg/pg/orm"
)

type SelectFilter struct {
	Address    string
	StartBlock *string
	EndBlock   *string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	blocksRange := blocks.RangeSelectFilter{StartBlock: f.StartBlock, EndBlock: f.EndBlock}

	return q.Where("address.address = ?", f.Address).Apply(blocksRange.Filter), nil
}
