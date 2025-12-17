package main

import (
	"backend/app"
	"backend/core/config"
	"backend/core/constants"
)

func main() {
	app.Start(config.Load(constants.ServiceScheduler))
}
