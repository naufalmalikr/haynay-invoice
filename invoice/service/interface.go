package service

import "haynay/haynay-invoice/invoice/entity"

type Service interface {
	Create(code string) entity.Invoice
	Save(invoice *entity.Invoice) (*entity.Invoice, error)
	Delete(invoice *entity.Invoice) (*entity.Invoice, error)
	ConfirmationMessage(invoice entity.Invoice) string
	Print(invoice entity.Invoice) string
}
