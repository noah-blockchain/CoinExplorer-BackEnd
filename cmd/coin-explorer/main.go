package main

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/api"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/core"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/internal/database"
)

func main() {
	// init environment
	env := core.NewEnvironment()

	// connect to database
	db := database.Connect(env)
	defer database.Close(db)

	// create explorer
	explorer := core.NewExplorer(db, env)

	// run api
	api.Run(db, explorer)
}
