package stake

import (
	"github.com/go-pg/pg"
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/tools"
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
		Column("Coin.symbol", "Validator.id", "Validator.public_key",
			"Validator.commission", "Validator.total_stake",
			"Validator.name", "Validator.description",
			"Validator.icon_url", "Validator.site_url",
			"OwnerAddress._").
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
		Column("OwnerAddress.address", "Coin.symbol").
		Column("value", "noah_value").
		Apply(pagination.Filter).
		SelectAndCount()

	helpers.CheckErr(err)
	return stakeDelegators
}

func (repository Repository) GetStakesForAddress(address string) (*[]models.Stake, error) {
	var stakes []models.Stake
	var err error

	_, err = repository.db.Query(&stakes, `
		SELECT s.noah_value, s.created_at, v.total_stake, v.commission, v.id, v.public_key, c.symbol
			FROM public.stakes as s 
			LEFT JOIN public.addresses as a on a.id = s.owner_address_id
			LEFT JOIN public.coins as c on c.id = s.coin_id
			LEFT JOIN public.validators as v on v.id = s.validator_id
			WHERE a.address=? AND v.status=? AND v.commission < 100 AND v.total_stake IS NOT NULL;
	`, helpers.RemoveNoahPrefix(address), models.ValidatorStatusReady)
	if err != nil {
		return nil, err
	}
	return &stakes, nil
}
