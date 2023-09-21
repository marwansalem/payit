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

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	accountController := &controllers.AccountController{
		Accounts: inMemoryAccountDataManager,
	}
	accountRoutes := router.Group("/accounts")
	accountRoutes.GET("", accountController.GetAll)
	accountRoutes.GET(":id", accountController.Get)
	accountRoutes.POST("", accountController.Create)
	accountRoutes.PUT(":id", accountController.Update)

	inMemoryTransferDataManager := &inmemory.InMemoryTransferDataManager{}
	inMemoryTransferDataManager.Init()
	transferController := &controllers.TransferController{
		TransferService: inmemory.TransferService{
			Accounts:  inMemoryAccountDataManager,
			Transfers: inMemoryTransferDataManager,
		},
	}
	transferRoutes := router.Group("/transfers")

	transferRoutes.GET("", transferController.GetAll)
	transferRoutes.GET(":id", transferController.Get)
	transferRoutes.POST("", transferController.Create)

	router.Run(":8080")

}
