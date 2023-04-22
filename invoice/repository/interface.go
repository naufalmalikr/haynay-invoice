package repository

import "haynay/haynay-invoice/invoice/entity"

type Repository interface {
	Insert(invoice *entity.Invoice) (*entity.Invoice, error)
	Delete(invoice *entity.Invoice) (*entity.Invoice, error)
}
