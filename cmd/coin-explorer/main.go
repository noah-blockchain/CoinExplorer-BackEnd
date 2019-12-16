package main

import (
	"github.com/noah-blockchain/noah-explorer-api/internal/api"
	"github.com/noah-blockchain/noah-explorer-api/internal/core"
	"github.com/noah-blockchain/noah-explorer-api/internal/database"
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
