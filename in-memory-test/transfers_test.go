package inmemory_test

import (
	"testing"

	inmemory "github.com/marwansalem/payit/in-memory"
	"github.com/marwansalem/payit/models"
	"github.com/stretchr/testify/assert"
)

func prepareTransferService() *inmemory.TransferService {
	accountData := inmemory.InMemoryAccountDataManager{}
	accountData.Init()

	transferData := inmemory.InMemoryTransferDataManager{}
	transferData.Init()

	transferService := inmemory.TransferService{
		Accounts:  &accountData,
		Transfers: &transferData,
	}
	return &transferService
}

func TestGivenNewTransferThenWithNegativeShouldFail(t *testing.T) {
	svc := prepareTransferService()

	_, err := svc.MakeTransfer("1", "2", -1.0)
	assert.EqualError(t, err, "amount must be positive")
}

func TestGivenNewTransferThenWithSenderSameAsReceiverShouldFail(t *testing.T) {
	svc := prepareTransferService()

	_, err := svc.MakeTransfer("1", "1", -1.0)
	assert.EqualError(t, err, "cannot transfer to self")
}

func TestGivenNewTransferThenWithNonExistingSenderShouldFail(t *testing.T) {
	svc := prepareTransferService()
	receiver := models.Account{
		ID:      "2",
		Name:    "Receiver",
		Balance: 1000,
	}
	svc.Accounts.Create(&receiver)

	_, err := svc.MakeTransfer("1", "2", 3.0)
	assert.EqualError(t, err, "could not find sender")
}

func TestGivenNewTransferThenWithNonExistingReceiverShouldFail(t *testing.T) {
	svc := prepareTransferService()
	sender := models.Account{
		ID:      "1",
		Name:    "sender",
		Balance: 1000,
	}
	svc.Accounts.Create(&sender)

	_, err := svc.MakeTransfer("1", "2", 3.0)
	assert.EqualError(t, err, "could not find receiver")
}

func TestGivenTwoValidAccountsAndAmountThenTransferShouldSucceed(t *testing.T) {
	svc := prepareTransferService()
	sender := models.Account{
		ID:      "1",
		Name:    "sender",
		Balance: 1000,
	}

	receiver := models.Account{
		ID:      "2",
		Name:    "receiver",
		Balance: 1000,
	}
	svc.Accounts.Create(&sender)
	svc.Accounts.Create(&receiver)

	_, err := svc.MakeTransfer("1", "2", 3.0)
	assert.NoError(t, err)
}

func TestGivenSenderAndReceiverThenTransferShouldModifyReceiverBalance(t *testing.T) {
	svc := prepareTransferService()
	sender := models.Account{
		ID:      "1",
		Name:    "sender",
		Balance: 1000,
	}

	receiver := models.Account{
		ID:      "2",
		Name:    "receiver",
		Balance: 1000,
	}
	svc.Accounts.Create(&sender)
	svc.Accounts.Create(&receiver)

	svc.MakeTransfer("1", "2", 3.0)
	receiverUpdated, _ := svc.Accounts.GetByID("2")
	assert.Equal(t, 1003.0, receiverUpdated.Balance)
}

func TestGivenSenderoAndReceiverThenAmountGreaterThanSenderBalanceShouldFail(t *testing.T) {
	svc := prepareTransferService()
	sender := models.Account{
		ID:      "1",
		Name:    "sender",
		Balance: 1000,
	}

	receiver := models.Account{
		ID:      "2",
		Name:    "receiver",
		Balance: 1000,
	}
	svc.Accounts.Create(&sender)
	svc.Accounts.Create(&receiver)

	_, err := svc.MakeTransfer("1", "2", sender.Balance+1)
	assert.Contains(t, err.Error(), "does not have enough balance")
}

// receiver := models.Account{
// 	ID:      "2",
// 	Name:    "Receiver",
// 	Balance: 1000,
// }
// svc.Accounts.Create(&sender)
// svc.Accounts.Create(&receiver)
