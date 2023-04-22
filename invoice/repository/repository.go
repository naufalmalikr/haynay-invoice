package repository

import (
	"encoding/json"
	"haynay/haynay-invoice/invoice/entity"
	"log"
	"os"
	"strconv"
	"time"
)

type repository struct {
	invoiceIdSeqFile string
}

func New() Repository {
	return &repository{
		invoiceIdSeqFile: "invoice_id_seq.txt",
	}
}

func (r *repository) Insert(invoice *entity.Invoice) (*entity.Invoice, error) {
	id, err := r.nextID()
	if err != nil {
		return nil, err
	}

	invoice.ID = id
	invoice.UpdatedAt = time.Now()
	invoiceJson, _ := json.Marshal(invoice)
	log.Printf("Invoice inserted: %s\n", string(invoiceJson))
	return invoice, nil
}

func (r *repository) nextID() (int64, error) {
	lastID, err := os.ReadFile(r.invoiceIdSeqFile)
	if err != nil {
		return 0, err
	}

	id, _ := strconv.ParseInt(string(lastID), 10, 64)
	id++
	file, err := os.Create(r.invoiceIdSeqFile)
	if err != nil {
		return id, err
	}
	defer file.Close()

	_, err = file.WriteString(strconv.FormatInt(id, 10))
	return id, err
}

func (r *repository) Delete(invoice *entity.Invoice) (*entity.Invoice, error) {
	invoiceJson, _ := json.Marshal(invoice)
	log.Printf("Invoice deleted: %s\n", string(invoiceJson))
	return invoice, nil
}
