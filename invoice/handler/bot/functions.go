package bot

import (
	"errors"
	"fmt"
	"haynay/haynay-invoice/invoice/entity"
	"strconv"
)

func noProcess(input string, user *user) error {
	return nil
}

func noRollback(user *user) error {
	return nil
}

func initializeInvoice(b *bot) processFunc {
	return func(input string, user *user) error {
		invoice := b.service.Create("LS")
		user.Invoice = &invoice
		return nil
	}
}

func printInvoice(b *bot) func(user *user) []string {
	return func(user *user) []string {
		return []string{
			"Berikut invoicenya, copy-paste ke WA",
			b.service.Print(*user.Invoice)}
	}
}

func saveInvoice(b *bot) processFunc {
	return func(input string, user *user) error {
		_, err := b.service.Save(user.Invoice)
		return err
	}
}

func deleteInvoice(b *bot) rollbackFunc {
	return func(user *user) error {
		_, err := b.service.Delete(user.Invoice)
		return err
	}
}

func setBuyer(input string, user *user) error {
	user.Invoice.Buyer = input
	return nil
}

func setPhone(input string, user *user) error {
	user.Invoice.Phone = input
	return nil
}

func setAddressAndInitializeOrder(input string, user *user) error {
	user.Invoice.Address = input
	return initializeOrder(input, user)
}

func initializeOrder(input string, user *user) error {
	user.Invoice.Orders = append(user.Invoice.Orders, entity.Order{})
	return nil
}

func removeLastOrder(user *user) error {
	user.Invoice.Orders = user.Invoice.Orders[:len(user.Invoice.Orders)-1]
	return nil
}

func setOrderItem(input string, user *user) error {
	nOrder := len(user.Invoice.Orders)
	user.Invoice.Orders[nOrder-1].Item = input
	return nil
}

func editFromOrderItem(b *bot) editedStepFunc {
	return func(user *user) *step {
		if len(user.Invoice.Orders) == 1 {
			return b.steps["address"]
		} else {
			return b.steps["order_more"]
		}
	}
}

func setOrderQuantity(input string, user *user) error {
	nOrder := len(user.Invoice.Orders)
	quantity, err := strconv.Atoi(input)
	if err != nil {
		return errInvalidNumber
	} else if quantity < 0 {
		return errNegative
	}

	user.Invoice.Orders[nOrder-1].Quantity = quantity
	return nil
}

func setOrderPrice(input string, user *user) error {
	nOrder := len(user.Invoice.Orders)
	price, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return errInvalidNumber
	} else if price < 0 {
		return errNegative
	}

	user.Invoice.Orders[nOrder-1].PricePerItem = price
	return nil
}

func setCourier(input string, user *user) error {
	user.Invoice.Courier = input
	return nil
}

func setCOD(input string, user *user) error {
	user.Invoice.COD = true
	return nil
}

func setFreeShippingFee(input string, user *user) error {
	user.Invoice.FreeShippingFee = true
	return nil
}

func rollbackCourierInfo(user *user) error {
	user.Invoice.COD = false
	user.Invoice.FreeShippingFee = false
	return nil
}

func setShippingFee(input string, user *user) error {
	price, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return errInvalidNumber
	} else if price < 0 {
		return errNegative
	}

	user.Invoice.ShippingFee = price
	return nil
}

func setPO(input string, user *user) error {
	user.Invoice.PO = true
	return nil
}

func setPODuration(input string, user *user) error {
	duration, err := strconv.Atoi(input)
	if err != nil {
		return errInvalidNumber
	} else if duration < 1 {
		return errors.New("tidak boleh kurang dari 1 pekan")
	}

	user.Invoice.PODurationWeek = duration
	return nil
}

func rollbackPO(user *user) error {
	user.Invoice.PO = false
	user.Invoice.PODurationWeek = 0
	return nil
}

func buyerOutput(user *user) []string {
	return []string{
		"Memulai bikin invoice.\n- Ketik \"Batal\" kapanpun untuk membatalkan invoice, atau \"Ubah\" untuk mengubah isian sebelumnya.\n- Jika terdapat pilihan, kamu bisa masukin angka, bahasa indonesia, ataupun bahasa inggris.",
		"Masukkan nama pembeli"}
}

func orderItemOutput(user *user) []string {
	return []string{fmt.Sprintf("Masukkan nama produk ke-%d", len(user.Invoice.Orders))}
}

func orderQuantityOutput(user *user) []string {
	return []string{fmt.Sprintf("Masukkan jumlah produk ke-%d", len(user.Invoice.Orders))}
}

func orderPriceOutput(user *user) []string {
	return []string{fmt.Sprintf("Masukkan harga satuan produk ke-%d", len(user.Invoice.Orders))}
}

func orderMoreOutput(user *user) []string {
	nOrder := len(user.Invoice.Orders) + 1
	return []string{fmt.Sprintf("Pilih langkah selanjutnya:\n1. Ketik \"OK\" untuk lanjut\n2. Ketik \"Tambah\" untuk memasukkan produk ke-%d", nOrder)}
}

func courierPriceOutput(user *user) []string {
	message := "Masukkan biaya ongkir"
	if user.Invoice.COD {
		message += "\n(meskipun COD, berguna untuk info ke pembeli)"
	} else if user.Invoice.FreeShippingFee {
		message += "\n(meskipun free ongkir, berguna untuk pencatatan)"
	}

	return []string{message}
}

func confirmationOutput(b *bot) func(user *user) []string {
	return func(user *user) []string {
		return []string{fmt.Sprintf("Konfirmasi:\n%s\nKetik \"OK\" untuk menyimpan invoice", b.service.ConfirmationMessage(*user.Invoice))}
	}
}

func editAction(step *step, currentUser *user) *action {
	if step == nil || step.EditedStep == nil {
		return &action{
			NextStep: step.Name,
			Process:  noProcess,
			Rollback: noRollback,
		}
	}

	editedStep := step.EditedStep(currentUser)

	// NOTE: this approach -using pointer to function type-
	//       MAY NOT works for dynamic rollbackFunc like deleteInvoice
	rollbacksMap := map[string]map[*rollbackFunc]bool{}
	if editedStep.FreeTextAction != nil {
		rollbacksMap[editedStep.FreeTextAction.NextStep] = map[*rollbackFunc]bool{
			&editedStep.FreeTextAction.Rollback: true,
		}
	}

	for _, action := range editedStep.ChoicesAction {
		if _, found := rollbacksMap[action.NextStep]; !found {
			rollbacksMap[action.NextStep] = map[*rollbackFunc]bool{}
		}

		if _, found := rollbacksMap[action.NextStep][&action.Rollback]; !found {
			rollbacksMap[action.NextStep][&action.Rollback] = true
		}
	}

	return &action{
		NextStep: editedStep.Name,
		Process: func(input string, user *user) error {
			if rollbacks, found := rollbacksMap[user.CurrentStep]; found {
				for rollback := range rollbacks {
					if err := (*rollback)(user); err != nil {
						return err
					}
				}
			}

			return nil
		},
		Rollback: noRollback,
	}
}
