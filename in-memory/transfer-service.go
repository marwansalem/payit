package inmemory

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/marwansalem/payit/models"
)

type TransferService struct {
	Accounts  *InMemoryAccountDataManager
	Transfers *inMemoryTransferDataManager
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

func (svc *TransferService) MakeTransfer(senderID, receiverID string, amount float64) (string, error) {
	transferID := uuid.NewString()

	transfer := &models.Transfer{
		ID:         transferID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
		Succeeded:  false,
	}

	go svc.Transfers.Create(transfer)
	if senderID == receiverID {
		return transfer.ID, fmt.Errorf("cannot transfer to self")
	}

	if amount == 0 {
		return transferID, fmt.Errorf("cannot transfer 0")
	}

	if amount < 0 {
		return transferID, fmt.Errorf("amount must be positive")
	}

	senderAccount, ok := svc.Accounts.getInMemoryByID(senderID)
	if !ok {
		return transferID, fmt.Errorf("could not find sender %s", senderID)
	}

	receiverAccount, ok := svc.Accounts.getInMemoryByID(senderID)
	if !ok {
		return transferID, fmt.Errorf("could not find receiver %s", receiverID)

	}

	// refetch account after Locking to ensure we have the latest balance, that will not changed until the transaction is over
	senderAccount = svc.lockAndFetchAccount(senderID)
	if senderAccount.Account.Balance < amount {
		senderAccount.Lock.Unlock()
		return transferID, fmt.Errorf("sender %s does not have enough balance, transfer amount: %f, balance: %f", senderAccount.ID, amount, senderAccount.Balance)
	}
	senderAccount.Account.Balance -= amount
	svc.Accounts.Update(senderAccount.Account)
	senderAccount.Lock.Unlock()

	receiverAccount = svc.lockAndFetchAccount(receiverID)
	receiverAccount.Account.Balance += amount
	svc.Accounts.Update(receiverAccount.Account)
	receiverAccount.Lock.Unlock()

	transfer.Succeeded = true
	go svc.Transfers.Update(transfer)
	return transferID, nil
}
