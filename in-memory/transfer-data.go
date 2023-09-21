package inmemory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/marwansalem/payit/models"
)

type InMemoryTransferDataManager struct {
	transfers *sync.Map
}

func (transferData *InMemoryTransferDataManager) Init() {
	transferData.transfers = &sync.Map{}
}

func (transferData *InMemoryTransferDataManager) Create(transfer *models.Transfer) (*models.Transfer, error) {
	transfer.ID = uuid.NewString()
	_, alreadyExists := transferData.transfers.LoadOrStore(transfer.ID, transfer)
	if alreadyExists {
		return nil, fmt.Errorf("transfer %v already exists", transfer.ID)
	}
	return transfer, nil
}

func (transferData *InMemoryTransferDataManager) GetByID(id string) (*models.Transfer, bool) {
	value, ok := transferData.transfers.Load(id)
	if ok {
		return value.(*models.Transfer), true
	}
	return nil, false
}

func (transferData *InMemoryTransferDataManager) List() *[]*models.Transfer {
	transfers := []*models.Transfer{}
	transferData.transfers.Range(func(_, value interface{}) bool {
		transfers = append(transfers, value.(*models.Transfer))
		return true
	})
	return &transfers
}

func (transferData *InMemoryTransferDataManager) Update(transfer *models.Transfer) error {
	_, alreadyExists := transferData.transfers.Load(transfer.ID)
	if !alreadyExists {
		return fmt.Errorf("transfer %v does not exist", transfer.ID)
	}

	transferData.transfers.Store(transfer.ID, transfer)
	return nil
}
