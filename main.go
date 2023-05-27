package main

import (
	"api-cariprice/config"
	"api-cariprice/routes"
	"flag"
)

func main() {
	config.LoadEnv()
	db := config.InitialDB()
	// config.Migration(db)

	flag.Parse()
	arg := flag.Arg(0)

	if arg != "" {
		config.InitCommands(db)
	} else {
		routes.WebRouter(db)
	}
}
