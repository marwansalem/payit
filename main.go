package main

import (
	"log"
	"os"

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
}
