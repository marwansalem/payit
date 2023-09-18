package data

import "github.com/marwansalem/payit/models"

type AccountData interface {
	Create() (*models.Account, error)
	CreateBulk(*[]*models.Account) error
	GetByID(id string) (*models.Account, error)
	List() ([]*models.Account, error)
	Update(*models.Account) error
}

type TransferData interface {
	Create() (*models.Transfer, error)
	GetByID(id string) (*models.Transfer, error)
	List() ([]*models.Transfer, error)
}
