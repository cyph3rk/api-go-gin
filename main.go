package main

import (
	"github.com/cyph3rk/api-go-gin/database"
	"github.com/cyph3rk/api-go-gin/routes"
)

func main() {

	database.ConectaComBancoDeDados()

	routes.HandleRequests()

}
