package inmemory

import (
	"fmt"
	"sync"

	"github.com/marwansalem/payit/models"
)

type inMemoryAccount struct {
	*models.Account
	Lock sync.RWMutex
}

type InMemoryAccountDataManager struct {
	accounts *sync.Map
}

func (accountData *InMemoryAccountDataManager) Init() {
	accountData.accounts = &sync.Map{}
}

func toInternalRepresentation(account *models.Account) *inMemoryAccount {
	return &inMemoryAccount{
		Account: account,
	}
}
func (accountData *InMemoryAccountDataManager) Create(account *models.Account) error {
	_, alreadyExists := accountData.accounts.LoadOrStore(account.ID, toInternalRepresentation(account))
	if alreadyExists {
		return fmt.Errorf("account %v already exists", account.ID)
	}
	return nil
}

func (accountData *InMemoryAccountDataManager) CreateBulk(accounts *[]*models.Account) error {
	for _, account := range *accounts {
		err := accountData.Create(account)
		if err != nil {
			return err
		}
	}
	return nil
}

func (accountData *InMemoryAccountDataManager) GetByID(id string) (*models.Account, bool) {
	value, ok := accountData.accounts.Load(id)
	if ok {
		return value.(*inMemoryAccount).Account, true
	}
	return nil, false
}

func (accountData *InMemoryAccountDataManager) List() *[]*models.Account {
	accounts := []*models.Account{}
	accountData.accounts.Range(func(_, value interface{}) bool {
		accounts = append(accounts, value.(*inMemoryAccount).Account)
		return true
	})
	return &accounts
}

func (accountData *InMemoryAccountDataManager) Update(account *models.Account) error {
	_, alreadyExists := accountData.accounts.Load(account.ID)
	if !alreadyExists {
		return fmt.Errorf("account %v does not exist", account.ID)
	}

	accountData.accounts.Store(account.ID, toInternalRepresentation(account))
	return nil
}
