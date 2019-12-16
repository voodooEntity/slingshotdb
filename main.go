package main

import (
	"fmt"
	"slingshot/api"
	"slingshot/storage"
)

func main() {
	fmt.Println("> Starting SlingshotDB ")
	storage.Boot()
	api.Start()
}
