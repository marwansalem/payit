package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/marwansalem/payit/controllers"
	inmemory "github.com/marwansalem/payit/in-memory"
	"github.com/marwansalem/payit/utils"
)

func main() {
	inMemoryAccountDataManager := &inmemory.InMemoryAccountDataManager{}
	inMemoryAccountDataManager.Init()
	err := utils.LoadAccounts(inMemoryAccountDataManager, "./accounts.json")
	if err != nil {
		log.Printf("Failed to Load Accounts, %v", err)
		os.Exit(-1)
	}
	log.Printf("Ready to receive requests")

	router := gin.Default()

	accountController := &controllers.AccountController{
		Accounts: inMemoryAccountDataManager,
	}
	accountRoutes := router.Group("/accounts")
	accountRoutes.GET("", accountController.GetAll)
	accountRoutes.GET(":id", accountController.Get)
	accountRoutes.POST("", accountController.Create)
	accountRoutes.PUT("", accountController.Update)
	router.Run(":8080")

}
