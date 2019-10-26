package address

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"github.com/go-pg/pg"
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
