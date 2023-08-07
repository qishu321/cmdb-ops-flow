package main

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/router"
)

func main() {
	models.InitDb()
	router.InitRouter()
}
