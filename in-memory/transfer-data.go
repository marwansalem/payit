package inmemory

import (
	"fmt"
	"sync"

	"github.com/marwansalem/payit/models"
)

type inMemoryTransfer struct {
	*models.Transfer
}

type inMemoryTransferDataManager struct {
	transfers *sync.Map
}

func (transferData *inMemoryTransferDataManager) Create(transfer *models.Transfer) error {
	_, alreadyExists := transferData.transfers.LoadOrStore(transfer.ID, transfer)
	if !alreadyExists {
		return fmt.Errorf("transfer %v does not exist", transfer.ID)
	}
	return nil
}

func (transferData *inMemoryTransferDataManager) GetByID(id string) (*models.Transfer, bool) {
	value, ok := transferData.transfers.Load(id)
	if ok {
		return value.(*inMemoryTransfer).Transfer, true
	}
	return nil, false
}

func (transferData *inMemoryTransferDataManager) List() *[]*models.Transfer {
	transfers := []*models.Transfer{}
	transferData.transfers.Range(func(_, value interface{}) bool {
		transfers = append(transfers, value.(*inMemoryTransfer).Transfer)
		return true
	})
	return &transfers
}

func (transferData *inMemoryTransferDataManager) Update(transfer *models.Transfer) error {
	_, alreadyExists := transferData.transfers.Load(transfer.ID)
	if !alreadyExists {
		return fmt.Errorf("transfer %v does not exist", transfer.ID)
	}

	transferData.transfers.Store(transfer.ID, transfer)
	return nil
}
