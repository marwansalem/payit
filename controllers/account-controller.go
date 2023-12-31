package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/marwansalem/payit/data"
	"github.com/marwansalem/payit/models"
)

type AccountController struct {
	Accounts data.AccountData
}

func (controller *AccountController) parseAccount(c *gin.Context) (*models.Account, error) {
	type accountJSON struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Balance string `json:"balance"`
	}

	var accountJson accountJSON
	if err := c.BindJSON(&accountJson); err != nil {
		return nil, err
	}
	balance, err := strconv.ParseFloat(accountJson.Balance, 64)
	if err != nil {
		return nil, err
	}
	return &models.Account{
		ID:      accountJson.ID,
		Name:    accountJson.Name,
		Balance: balance,
	}, nil

}

func (controller *AccountController) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Printf("Failed to fetch account, id is not provided")
		BadRequest(c, "invalid account id")
		return
	}
	account, found := controller.Accounts.GetByID(id)
	if !found {
		NotFound(c, fmt.Sprintf("account %s not found", id))
		return
	}

	c.JSON(http.StatusOK, account)
}

func (controller *AccountController) GetAll(c *gin.Context) {
	accounts := controller.Accounts.List()
	c.JSON(http.StatusOK, accounts)
}

func (controller *AccountController) Create(c *gin.Context) {
	account, err := controller.parseAccount(c)
	if err != nil {
		log.Printf("Failed to create account, cannot parse request body %v", err)
		BadRequest(c, err.Error())
		return
	}

	if account.ID != "" {
		message := "ID is not allowed in request body, will be assigned after creation"
		log.Print(message)
		BadRequest(c, message)
		return
	}

	if account.Balance < 0 {
		log.Printf("Failed to create account %s, balance %f is less than 0", account.ID, account.Balance)
		BadRequest(c, "Balance cannot be less than 0")
		return
	}

	account, err = controller.Accounts.Create(account)
	if err != nil {
		log.Printf("Error occurred while creating account %s, %v", account.ID, err)

		ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, account)
}

func (controller *AccountController) Update(c *gin.Context) {
	id := c.Param("id")
	account, err := controller.parseAccount(c)

	if id != account.ID {
		BadRequest(c, "id in endpoint does not match id in request body")
		return
	}
	if err != nil {
		log.Printf("Failed to update account %s, invalid account: %v", account.ID, err)
		BadRequest(c, err.Error())
		return
	}

	if account.Balance < 0 {
		log.Printf("Failed to update account %s, balance %f is less than 0", account.ID, account.Balance)
		BadRequest(c, "Balance cannot be less than 0")
		return
	}

	err = controller.Accounts.Update(account)
	if err != nil {
		log.Printf("Error occurred while updating account %s, %v", account.ID, err)

		ErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}
