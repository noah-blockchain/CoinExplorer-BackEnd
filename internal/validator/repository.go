package validator

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/blocks"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/tools"
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

func (repository Repository) GetByPublicKey(publicKey string) *models.Validator {
	var validator models.Validator

	err := repository.db.Model(&validator).
		Column("Stakes", "Stakes.Coin", "Stakes.OwnerAddress").
		Where("public_key = ?", publicKey).
		Select()

	if err != nil {
		return nil
	}

	return &validator
}

func (repository Repository) GetTotalStakeByActiveValidators(ids []uint64) string {
	var total string

	// get total stake of active validators
	err := repository.db.Model((*models.Validator)(nil)).
		ColumnExpr("SUM(total_stake)").
		Where("id IN (?)", pg.In(ids)).
		Select(&total)

	helpers.CheckErr(err)
	return total
}

func (repository Repository) GetActiveValidatorIds() []uint64 {
	var blockValidator models.BlockValidator
	var ids []uint64

	// get active validators by last block
	err := repository.db.Model(&blockValidator).
		Column("validator_id").
		Where("block_id = ?", blocks.NewRepository(repository.db).GetLastBlock().ID).
		Select(&ids)

	helpers.CheckErr(err)
	return ids
}

// Get active candidates count
func (repository Repository) GetActiveCandidatesCount() int {
	var validator models.Validator

	count, err := repository.db.Model(&validator).
		Where("status = ?", models.ValidatorStatusReady).
		Count()

	helpers.CheckErr(err)
	return count
}

// Get validators
func (repository Repository) GetValidators() []models.Validator {
	var validators []models.Validator

	err := repository.db.Model(&validators).Select()

	helpers.CheckErr(err)
	return validators
}

func (repository Repository) GetValidatorsBySymbol(coinSymbol string, pagination *tools.Pagination) []models.Validator {
	var validators []models.Validator
	var err error

	pagination.Total, err = repository.db.Model(&validators).
		Join("INNER JOIN stakes as s").
		JoinOn("s.validator_id = validator.id").
		Join("INNER JOIN coins as c").
		JoinOn("s.coin_id = c.id").
		Where("c.symbol=?", coinSymbol).
		ColumnExpr("DISTINCT validator.public_key").
		Column("validator.name", "validator.site_url", "validator.icon_url", "validator.description").
		Apply(pagination.Filter).
		SelectAndCount()

	helpers.CheckErr(err)
	return validators
}

func (repository Repository) GetValidatorsWithPagination(pagination *tools.Pagination, field *string, orderBy *string) []models.Validator {
	var validators []models.Validator
	var err error
	fieldSql := "uptime"
	orderBySql := "DESC"

	if field != nil {
		fieldSql = *field
	}

	if orderBy != nil {
		orderBySql = *orderBy
	}

	query := repository.db.Model(&validators).
		Apply(pagination.Filter).
		Order(fmt.Sprintf("validator.%s %s", fieldSql, orderBySql))

	pagination.Total, err = query.SelectAndCount()

	helpers.CheckErr(err)
	return validators
}
