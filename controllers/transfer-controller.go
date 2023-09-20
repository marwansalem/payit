package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	inmemory "github.com/marwansalem/payit/in-memory"
	"github.com/marwansalem/payit/models"
)

type TransferController struct {
	TransferService inmemory.TransferService
}

func (controller *TransferController) parseTransfer(c *gin.Context) (*models.Transfer, error) {
	var transfer models.Transfer
	if err := c.BindJSON(&transfer); err != nil {
		return nil, err
	}
	return &transfer, nil
}

func (controller *TransferController) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Printf("Failed to fetch transfer, id is not provided")
		BadRequest(c, "transfer id not found")
		return
	}
	transfer, found := controller.TransferService.Transfers.GetByID(id)
	if !found {
		NotFound(c, fmt.Sprintf("transfer %s not found", id))
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (controller *TransferController) GetAll(c *gin.Context) {
	transfers := controller.TransferService.Transfers.List()

	c.JSON(http.StatusOK, transfers)
}

func (controller *TransferController) Create(c *gin.Context) {
	transfer, err := controller.parseTransfer(c)
	if err != nil {
		BadRequest(c, "failed to parse transfer request")
		return
	}

	if transfer.ID != "" {
		BadRequest(c, "transfer id should not be set while making a request")
		return
	}

	if transfer.SenderID == "" {
		BadRequest(c, "sender id must be specified")
		return
	}

	if transfer.ReceiverID == "" {
		BadRequest(c, "receiver id must be specified")
		return
	}

	if transfer.Amount <= 0 {
		BadRequest(c, "transfer amount must greater than 0")
		return
	}
	transfer, err = controller.TransferService.MakeTransfer(transfer.SenderID, transfer.ReceiverID, transfer.Amount)
	if err != nil {
		ErrorResponse(c, err)
		return
	}
	if !transfer.Succeeded {
		ErrorResponse(c, fmt.Errorf("transfer failed"))
		return
	}
	c.JSON(http.StatusOK, transfer)
}
