package main

import (
	"server/api"
	_ "server/api/user"
)

func main() {
	api.NewServer().Run(":8082")
}
