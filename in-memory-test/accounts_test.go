package inmemory_test

import (
	"testing"

	inmemory "github.com/marwansalem/payit/in-memory"
	"github.com/marwansalem/payit/models"
	"github.com/stretchr/testify/assert"
)

func TestGivenAccountThenCreatingAccountShouldSucceed(t *testing.T) {
	accountData := inmemory.InMemoryAccountDataManager{}
	accountData.Init()

	newAccount := models.Account{
		Name:    "PayIt Customer",
		Balance: 1000.0,
	}

	_, err := accountData.Create(&newAccount)
	assert.NoError(t, err)
}

func TestGivenAccountWithoutIDThenCreatingAccountShouldAssignID(t *testing.T) {
	accountData := inmemory.InMemoryAccountDataManager{}
	accountData.Init()

	newAccount := models.Account{
		Name:    "PayIt Customer",
		Balance: 1000.0,
	}

	account, _ := accountData.Create(&newAccount)
	assert.NotEmpty(t, account.ID)
}

func TestGivenAccountWithIDThenCreatingAccountShouldNotAssignID(t *testing.T) {
	accountData := inmemory.InMemoryAccountDataManager{}
	accountData.Init()

	originalID := "some id"
	newAccount := models.Account{
		ID:      originalID,
		Name:    "PayIt Customer",
		Balance: 1000.0,
	}

	account, _ := accountData.Create(&newAccount)
	assert.Equal(t, originalID, account.ID)
}

func TestGivenCreatedAccountThenGetAccountShouldSucceeed(t *testing.T) {
	accountData := inmemory.InMemoryAccountDataManager{}
	accountData.Init()

	createdAccount := models.Account{
		Name:    "PayIt Customer",
		Balance: 10.0,
	}
	account, _ := accountData.Create(&createdAccount)

	_, exists := accountData.GetByID(account.ID)
	assert.True(t, exists)
}

func TestGivenNoAccountThenGetAccountShouldFail(t *testing.T) {
	accountData := inmemory.InMemoryAccountDataManager{}
	accountData.Init()

	_, exists := accountData.GetByID("some id")
	assert.False(t, exists)
}

func TestGivenAccountThenUpdatingShouldFetchReflectChange(t *testing.T) {
	accountData := inmemory.InMemoryAccountDataManager{}
	accountData.Init()

	originalAccount := models.Account{
		Name:    "PayIt Customer",
		Balance: 10.0,
	}
	account, _ := accountData.Create(&originalAccount)
	account.Balance = 100.0

	accountData.Update(account)

	updated, _ := accountData.GetByID(account.ID)

	assert.Equal(t, updated.Balance, 100.0)
}
