package main

import (
	"log"

	"github.com/marwansalem/payit/utils"
)

func main() {

	err := utils.LoadAccounts(nil, "./accounts.json")
	if err != nil {
		log.Printf("Failed to Load Accounts, %v", err)
	}
}
