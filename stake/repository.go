package stake

import (
	"github.com/go-pg/pg"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/tools"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Get list of stakes by Noah address
func (repository Repository) GetByAddress(address string) []*models.Stake {
	var stakes []*models.Stake

	err := repository.db.Model(&stakes).
		Column("Coin", "OwnerAddress._").
		Where("owner_address.address = ?", address).
		Select()

	helpers.CheckErr(err)

	return stakes
}

// Get paginated list of stakes by Noah address
func (repository Repository) GetPaginatedByAddress(address string, pagination *tools.Pagination) []models.Stake {
	var stakes []models.Stake
	var err error

	pagination.Total, err = repository.db.Model(&stakes).
		Column("Coin.symbol", "Validator.public_key", "OwnerAddress._").
		Column("Validator.name", "Validator.description", "Validator.icon_url", "Validator.site_url").
		Where("owner_address.address = ?", address).
		Apply(pagination.Filter).
		SelectAndCount()

	helpers.CheckErr(err)

	return stakes
}

// Get total delegated noah value
func (repository Repository) GetSumInNoahValue() (string, error) {
	var sum string
	err := repository.db.Model(&models.Stake{}).ColumnExpr("SUM(noah_value)").Select(&sum)
	return sum, err
}

// Get total delegated sum by address
func (repository Repository) GetSumInNoahValueByAddress(address string) (string, error) {
	var sum string
	err := repository.db.Model(&models.Stake{}).
		Column("OwnerAddress._").
		ColumnExpr("SUM(noah_value)").
		Where("owner_address.address = ?", address).
		Select(&sum)

	return sum, err
}


// Get paginated list of stakes by Noah address
func (repository Repository) GetPaginatedStakeForCoin(coinSymbol string, pagination *tools.Pagination) []models.Stake {
	var stakes []models.Stake
	var err error

	pagination.Total, err = repository.db.Model(&stakes).
		Join("LEFT JOIN coins as c").
		JoinOn("c.id = stake.coin_id").
		Where("c.symbol = ?", coinSymbol).
		Column("Validator.public_key", "OwnerAddress.address").
		Column("value", "noah_value").
		Apply(pagination.Filter).
		SelectAndCount()

	helpers.CheckErr(err)

	return stakes
}

// Get paginated list of delegators by validator pubKey
func (repository Repository) GetPaginatedDelegatorsForValidator(pubKey string, pagination *tools.Pagination) []models.Stake {
	var stakeDelegators []models.Stake
	var err error

	pagination.Total, err = repository.db.Model(&stakeDelegators).
		Join("LEFT JOIN validators as v").
		JoinOn("v.id = stake.validator_id").
		Where("v.public_key = ?", pubKey).
		Column("OwnerAddress.address").
		Column("value", "noah_value").
		Apply(pagination.Filter).
		SelectAndCount()

	helpers.CheckErr(err)
	return stakeDelegators
}
