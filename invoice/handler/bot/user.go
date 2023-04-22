package bot

import "haynay/haynay-invoice/invoice/entity"

type user struct {
	CurrentStep string
	Invoice     *entity.Invoice
}
