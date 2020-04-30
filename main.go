package main

import (
	"slingshot/api"
	"slingshot/config"
	"slingshot/storage"
)

func main() {
	config.Logger.Println("> Starting SlingshotDB ")
	storage.Boot()
	api.Start()
}
