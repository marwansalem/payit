package data

import "github.com/marwansalem/payit/models"

type AccountData interface {
	Create(*models.Account) error
	CreateBulk(*[]*models.Account) error
	Update(*models.Account) error
	GetByID(id string) (*models.Account, bool)
	List() *[]*models.Account
}

type TransferData interface {
	Create(*models.Transfer) (*models.Transfer, error)
	GetByID(id string) (*models.Transfer, bool)
	List() *[]*models.Transfer
}
