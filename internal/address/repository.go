package address

import (
	"github.com/go-pg/pg"
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/tools"
)

type Repository struct {
	DB *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// Get address model by address
func (repository Repository) GetByAddress(noahAddress string) *models.Address {
	var address models.Address

	err := repository.DB.Model(&address).Column("Balances", "Balances.Coin").
		Where("address = ?", noahAddress).Select()
	if err != nil {
		return nil
	}

	return &address
}

// Get list of addresses models
func (repository Repository) GetByAddresses(noahAddresses []string) []models.Address {
	var addresses []models.Address

	err := repository.DB.Model(&addresses).Column("Balances", "Balances.Coin").
		WhereIn("address IN (?)", pg.In(noahAddresses)).Select()

	helpers.CheckErr(err)

	return addresses
}

// Get address model by address
func (repository Repository) GetBalancesByCoinSymbol(coinSymbol string, pagination *tools.Pagination) []models.Balance {
	var balances []models.Balance
	var err error

	pagination.Total, err = repository.DB.Model(&balances).
		Join("LEFT JOIN coins as c").
		JoinOn("balance.coin_id = c.id").
		Where("c.symbol=?", coinSymbol).
		Column("Address.address", "balance.value").
		Apply(pagination.Filter).
		SelectAndCount()

	helpers.CheckErr(err)
	return balances
}
