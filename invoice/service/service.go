package service

import (
	"fmt"
	"haynay/haynay-invoice/invoice/entity"
	"haynay/haynay-invoice/invoice/repository"
	"haynay/haynay-invoice/utility"
	"math"
	"time"
)

const (
	printFormat = "*INVOICE %s*                        %s\nBismillaah                                   Linabutik.id\n\n%s*List Order:*\n%s+ Ongkos Kirim: %s\n--------------------------------------------------------------\n*Total: %s*\n--------------------------------------------------------------\nStatus: ✅️ %s"
)

type service struct {
	repository            repository.Repository
	downPaymentMultiplier int64
}

func New(repository repository.Repository) Service {
	return &service{
		repository:            repository,
		downPaymentMultiplier: 50000,
	}
}

func (s *service) Create(code string) entity.Invoice {
	return entity.Invoice{
		Code:      code,
		Orders:    []entity.Order{},
		CreatedAt: time.Now(),
	}
}

func (s *service) Save(invoice *entity.Invoice) (*entity.Invoice, error) {
	s.calculate(invoice)
	return s.repository.Insert(invoice)
}

func (s *service) Delete(invoice *entity.Invoice) (*entity.Invoice, error) {
	return s.repository.Delete(invoice)
}

func (s *service) ConfirmationMessage(invoice entity.Invoice) string {
	buyer := fmt.Sprintf("%s %s\n%s\n", invoice.Buyer, invoice.Phone, invoice.Address)
	orders := ""
	for _, order := range invoice.Orders {
		orders += fmt.Sprintf("- %s %d pcs @ %s\n", order.Item, order.Quantity, utility.Rupiah(order.PricePerItem))
	}

	shipping := fmt.Sprintf("%s %s", invoice.Courier, utility.Rupiah(invoice.ShippingFee))
	if invoice.COD {
		shipping += ", COD"
	} else if invoice.FreeShippingFee {
		shipping += ", free ongkir"
	}
	shipping += "\n"

	stock := "Ready stock\n"
	if invoice.PO {
		stock = fmt.Sprintf("Preorder %d minggu\n", invoice.PODurationWeek)
	}

	return fmt.Sprintf("%s%s%s%s", buyer, orders, shipping, stock)
}

func (s *service) Print(invoice entity.Invoice) string {
	s.calculate(&invoice)
	id := fmt.Sprintf("%s-%04d", invoice.Code, invoice.ID)
	date := fmt.Sprintf("%02d%02d%d", invoice.CreatedAt.Day(), invoice.CreatedAt.Month(), invoice.CreatedAt.Year())
	buyer := fmt.Sprintf("_Dear,_\n*%s*\n%s\n\n%s\n\n", invoice.Buyer, invoice.Phone, invoice.Address)
	orders := ""
	for _, order := range invoice.Orders {
		orders += fmt.Sprintf("+ %dx %s\n    %s\n", order.Quantity, order.Item, utility.Rupiah(order.Price))
	}

	shipping := invoice.Courier
	if invoice.COD {
		shipping += " bayar di tempat"
	} else if invoice.FreeShippingFee {
		shipping = "Free Ongkir " + shipping
	}
	shipping = fmt.Sprintf("%s ( _%s_ )", utility.Rupiah(invoice.ShippingFee), shipping)

	stock := "*Ready*"
	if invoice.PO {
		stock = fmt.Sprintf("*PO ±%d Minggu*\n\n⭐️\n• _DP: *min %s_\n• _*Pelunasan: setelah barang ready*_", invoice.PODurationWeek, utility.Rupiah(invoice.MinDP))
	}

	return fmt.Sprintf(printFormat, id, date, buyer, orders, shipping, utility.Rupiah(invoice.Price), stock)
}

func (s *service) calculate(invoice *entity.Invoice) *entity.Invoice {
	invoice.Price = 0
	for i, order := range invoice.Orders {
		invoice.Orders[i].Price = int64(order.Quantity) * order.PricePerItem
		invoice.Price += invoice.Orders[i].Price
	}

	invoice.ActualShippingFee = invoice.ShippingFee
	if invoice.COD || invoice.FreeShippingFee {
		invoice.ActualShippingFee = 0
	}
	invoice.Price += invoice.ActualShippingFee

	if invoice.PO {
		invoice.MinDP = int64(math.Ceil(float64(invoice.Price/2)/float64(s.downPaymentMultiplier))) * s.downPaymentMultiplier
		if invoice.MinDP > invoice.Price {
			invoice.MinDP = invoice.Price
		}
	}

	return invoice
}
