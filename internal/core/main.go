package core

import (
	"github.com/go-pg/pg"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/address"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/blocks"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/coins"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/invalid_transaction"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/reward"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/slash"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/stake"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/tools/cache"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/transaction"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/validator"
)

type Explorer struct {
	CoinRepository               coins.Repository
	BlockRepository              blocks.Repository
	AddressRepository            address.Repository
	TransactionRepository        transaction.Repository
	InvalidTransactionRepository invalid_transaction.Repository
	RewardRepository             reward.Repository
	SlashRepository              slash.Repository
	ValidatorRepository          validator.Repository
	StakeRepository              stake.Repository
	Environment                  Environment
	Cache                        *cache.ExplorerCache
}

func NewExplorer(db *pg.DB, env *Environment) *Explorer {
	return &Explorer{
		CoinRepository:               *coins.NewRepository(db, env.BaseCoin),
		BlockRepository:              *blocks.NewRepository(db),
		AddressRepository:            *address.NewRepository(db),
		TransactionRepository:        *transaction.NewRepository(db),
		InvalidTransactionRepository: *invalid_transaction.NewRepository(db),
		RewardRepository:             *reward.NewRepository(db),
		SlashRepository:              *slash.NewRepository(db),
		ValidatorRepository:          *validator.NewRepository(db),
		StakeRepository:              *stake.NewRepository(db),
		Environment:                  *env,
		Cache:                        cache.NewCache(),
	}
}
