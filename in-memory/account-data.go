package inmemory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
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

func (accountData *InMemoryAccountDataManager) Create(account *models.Account) (*models.Account, error) {
	if account.ID == "" {
		account.ID = uuid.NewString()
	}
	_, alreadyExists := accountData.accounts.LoadOrStore(account.ID, toInternalRepresentation(account))
	if alreadyExists {
		return nil, fmt.Errorf("account %v already exists", account.ID)
	}
	return account, nil
}

func (accountData *InMemoryAccountDataManager) CreateBulk(accounts *[]*models.Account) error {
	for _, account := range *accounts {
		_, err := accountData.Create(account)
		if err != nil {
			return err
		}
	}
	return nil
}

func (accountData *InMemoryAccountDataManager) GetByID(id string) (*models.Account, bool) {
	inMemoryAccount, exists := accountData.getInMemoryByID(id)
	if !exists {
		return nil, false
	}
	return inMemoryAccount.Account, true
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

func (accountData *InMemoryAccountDataManager) getInMemoryByID(id string) (*inMemoryAccount, bool) {
	value, ok := accountData.accounts.Load(id)
	if ok {
		return value.(*inMemoryAccount), true
	}
	return nil, false
}
