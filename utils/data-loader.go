package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/marwansalem/payit/data"
	"github.com/marwansalem/payit/models"
)

func loadAccountsFromFile(filePath string) (*[]*models.Account, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Define a custom type to perform the balance conversion
	type accountJSON struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Balance string `json:"balance"`
	}

	// Unmarshal the JSON data into a slice of *Account pointers with custom balance conversion
	var accountsJSON []accountJSON
	if err := json.Unmarshal(data, &accountsJSON); err != nil {
		return nil, err
	}

	var accounts []*models.Account
	for _, accJSON := range accountsJSON {
		balance, err := strconv.ParseFloat(accJSON.Balance, 64)
		if err != nil {
			return nil, err
		}

		account := &models.Account{
			ID:      accJSON.ID,
			Name:    accJSON.Name,
			Balance: balance,
		}

		accounts = append(accounts, account)
	}

	return &accounts, nil
}

func LoadAccounts(accountData data.AccountData, filePath string) error {
	accounts, err := loadAccountsFromFile(filePath)
	if err != nil {
		return err
	}
	for _, acc := range *accounts {
		fmt.Printf("%v\n", acc)
	}
	return accountData.CreateBulk(accounts)
}
