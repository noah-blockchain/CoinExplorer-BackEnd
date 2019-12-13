package core

import (
	"github.com/go-pg/pg"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/address"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/balance"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/blocks"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/coins"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/invalid_transaction"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/noahdev"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/reward"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/slash"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/stake"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/tools/cache"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/tools/market"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/transaction"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/validator"
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
	MarketService                *market.Service
	BalanceService               *balance.Service
}

func NewExplorer(db *pg.DB, env *Environment) *Explorer {
	marketService := market.NewService(noahdev.NewApi(env.NoahDevApiHost), env.BaseCoin)
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
		MarketService:                marketService,
		BalanceService:               balance.NewService(env.BaseCoin, marketService),
	}
}
