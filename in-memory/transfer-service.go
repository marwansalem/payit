package inmemory

import (
	"fmt"

	"github.com/marwansalem/payit/models"
)

type TransferService struct {
	Accounts  *InMemoryAccountDataManager
	Transfers *InMemoryTransferDataManager
}

func (svc *TransferService) lockAndFetchAccount(accountID string) *inMemoryAccount {
	accoount, ok := svc.Accounts.getInMemoryByID(accountID)
	if !ok {
		return nil
	}
	accoount.Lock.Lock()
	accoount, _ = svc.Accounts.getInMemoryByID(accountID)

	return accoount
}

// TODO define the standard errors, so they can be used by other implementation of TransferService
func (svc *TransferService) MakeTransfer(senderID, receiverID string, amount float64) (*models.Transfer, error) {

	transfer := &models.Transfer{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
		Succeeded:  false,
	}

	transfer, err := svc.Transfers.Create(transfer)
	if err != nil {
		return transfer, err
	}
	if senderID == receiverID {
		return transfer, fmt.Errorf("cannot transfer to self")
	}

	if amount <= 0 {
		return transfer, fmt.Errorf("amount must be positive")
	}

	_, senderExists := svc.Accounts.getInMemoryByID(senderID)
	if !senderExists {
		return transfer, fmt.Errorf("could not find sender")
	}

	_, receiverExists := svc.Accounts.getInMemoryByID(receiverID)
	if !receiverExists {
		return transfer, fmt.Errorf("could not find receiver")

	}

	// refetch account after Locking to ensure we have the latest balance, that will not changed until the transaction is over
	senderAccount := svc.lockAndFetchAccount(senderID)
	if senderAccount.Account.Balance < amount {
		senderAccount.Lock.Unlock()
		return transfer, fmt.Errorf("sender %s does not have enough balance, transfer amount: %f, balance: %f", senderAccount.ID, amount, senderAccount.Balance)
	}
	senderAccount.Account.Balance -= amount
	svc.Accounts.Update(senderAccount.Account)
	senderAccount.Lock.Unlock()

	receiverAccount := svc.lockAndFetchAccount(receiverID)
	receiverAccount.Account.Balance += amount
	svc.Accounts.Update(receiverAccount.Account)
	receiverAccount.Lock.Unlock()

	transfer.Succeeded = true
	svc.Transfers.Update(transfer)
	return transfer, nil
}
