package main

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/api"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/core"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/database"
)

func main() {
	// init environment
	env := core.NewEnvironment()

	// connect to database
	db := database.Connect(env)
	defer database.Close(db)

	// create explorer
	explorer := core.NewExplorer(db, env)

	// run market price update
	go explorer.MarketService.Run()

	// create ws extender
	//extender := core.NewExtenderWsClient(explorer)
	//defer extender.Close()

	// subscribe to channel and add cache handler
	//sub := extender.CreateSubscription(explorer.Environment.WsBlocksChannel)
	//sub.OnPublish(explorer.Cache)
	//extender.Subscribe(sub)

	// run api
	api.Run(db, explorer)
}
