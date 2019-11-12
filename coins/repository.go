package coins

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/tools"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Repository struct {
	DB             *pg.DB
	baseCoinSymbol string
}

func NewRepository(db *pg.DB, baseCoinSymbol string) *Repository {
	return &Repository{
		DB:             db,
		baseCoinSymbol: baseCoinSymbol,
	}
}

type CustomCoinsStatusData struct {
	ReserveSum string
	Count      uint
}

// Get custom coins data for status page
func (repository *Repository) GetCustomCoinsStatusData() (CustomCoinsStatusData, error) {
	var data CustomCoinsStatusData

	err := repository.DB.
		Model(&models.Coin{}).
		ColumnExpr("SUM(reserve_balance) as reserve_sum, COUNT(*) as count").
		Where("symbol != ?", repository.baseCoinSymbol).
		Select(&data)

	return data, err
}

// Get paginated list of blocks
func (repository Repository) GetPaginated(pagination *tools.Pagination, field *string, orderBy *string, symbol *string) []models.Coin {
	var coins []models.Coin
	var err error
	fieldSql := "reserve_balance"
	orderBySql := "DESC"

	if field != nil {
		fieldSql = *field
	}

	if orderBy != nil {
		orderBySql = *orderBy
	}

	query := repository.DB.Model(&coins).
		Column("coin.crr", "coin.volume", "coin.reserve_balance", "coin.name", "coin.symbol", "coin.price", "coin.delegated", "coin.updated_at", "coin.created_at", "coin.capitalization", "a.address").
		Apply(pagination.Filter).
		Join("LEFT JOIN addresses AS a ON a.id = coin.creation_address_id")

	if symbol != nil {
		query = query.Where("coin.symbol LIKE ?", fmt.Sprintf("%%%s%%", *symbol))
	}

	query = query.Where("coin.deleted_at IS NULL").
		Order(fmt.Sprintf("coin.%s %s", fieldSql, orderBySql))

	pagination.Total, err = query.SelectAndCount()
	helpers.CheckErr(err)

	return coins
}

// Get coin by symbol
func (repository Repository) GetBySymbol(symbol string) *models.Coin {
	var coin models.Coin

	err := repository.DB.Model(&coin).
		Column("coin.crr", "coin.volume", "coin.reserve_balance", "coin.name", "coin.symbol", "coin.price", "coin.delegated", "coin.updated_at", "coin.created_at", "coin.capitalization", "a.address").
		Join("LEFT JOIN addresses AS a ON a.id = coin.creation_address_id").
		Where("coin.symbol LIKE ?", fmt.Sprintf("%%%s%%", symbol)).
		Select()

	if err != nil {
		return nil
	}

	return &coin
}
