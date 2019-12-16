package main

import (
	"github.com/noah-blockchain/noah-explorer-extender/internal/api"
	"github.com/noah-blockchain/noah-explorer-extender/internal/core"
	"github.com/noah-blockchain/noah-explorer-extender/internal/database"
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
