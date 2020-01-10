package reward

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/aggregated_reward"
	"github.com/noah-blockchain/noah-explorer-api/internal/events"
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

// Get filtered list of rewards by Noah address
func (repository Repository) GetPaginatedByAddress(filter events.SelectFilter, pagination *tools.Pagination) []models.Reward {
	var rewards []models.Reward
	var err error

	// get count of rewards
	pagination.Total, err = repository.db.Model(&rewards).
		Column("Address.address").
		Apply(filter.Filter).
		Count()
	helpers.CheckErr(err)

	if pagination.Total == 0 {
		return nil
	}

	// get rewards
	err = repository.db.Model(&rewards).
		Column("Address.address", "Validator.public_key", "Block.created_at").
		Column("Validator.name", "Validator.description", "Validator.icon_url", "Validator.site_url").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("block.id DESC").
		Order("reward.amount").
		Select()
	helpers.CheckErr(err)

	return rewards
}

type ChartData struct {
	Time   time.Time `json:"time"`
	Amount string    `json:"amount"`
}

// Get filtered chart data by Noah address
func (repository Repository) GetChartData(address string, filter tools.Filter) []ChartData {
	var rewards models.Reward
	var chartData []ChartData

	err := repository.db.Model(&rewards).
		Column("Address._").
		ColumnExpr("SUM(amount) as amount").
		Where("address.address = ?", address).
		Apply(filter.Filter).
		Select(&chartData)

	helpers.CheckErr(err)

	return chartData
}

func (repository Repository) GetAggregatedChartData(filter aggregated_reward.SelectFilter) []ChartData {
	var rewards models.AggregatedReward
	var chartData []ChartData

	err := repository.db.Model(&rewards).
		Column("Address._").
		ColumnExpr("SUM(amount) as amount").
		ColumnExpr("time_id as time").
		Group("time").
		Order("time").
		Apply(filter.Filter).
		Select(&chartData)

	helpers.CheckErr(err)

	return chartData
}

func (repository Repository) GetPaginatedAggregatedByAddress(filter aggregated_reward.SelectFilter, pagination *tools.Pagination) []models.AggregatedReward {
	var rewards []models.AggregatedReward
	var err error

	// get rewards
	pagination.Total, err = repository.db.Model(&rewards).
		Column("Address.address", "Validator").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("time_id DESC").
		Order("amount").
		SelectAndCount()

	helpers.CheckErr(err)

	return rewards
}

func (repository Repository) GetSumRewardForValidator(validatorId uint64, createdAt time.Time) string {
	var total string

	// get total stake of active validators
	err := repository.db.Model((*models.Reward)(nil)).
		ColumnExpr("SUM(amount)").
		Where("role = ?", "Validator").
		Where("created_at >= ?", createdAt).
		Where("validator_id = ?", validatorId).
		Select(&total)
	if err != nil {
		return "0"
	}

	return total
}
