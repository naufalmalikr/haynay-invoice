package repository

import (
	"encoding/json"
	"haynay/haynay-invoice/invoice/entity"
	"log"
	"sync"
	"time"
)

type repository struct {
	invoiceIDSeq int64
	idMutex      sync.Mutex
}

func New(invoiceIDSeq int64) Repository {
	return &repository{invoiceIDSeq: invoiceIDSeq}
}

func (r *repository) Insert(invoice *entity.Invoice) (*entity.Invoice, error) {
	r.idMutex.Lock()
	defer r.idMutex.Unlock()

	r.invoiceIDSeq++
	invoice.ID = r.invoiceIDSeq
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = invoice.CreatedAt
	invoiceJson, err := json.Marshal(invoice)
	if err != nil {
		r.invoiceIDSeq--
		return invoice, err
	}

	log.Printf("Invoice inserted: %s\n", string(invoiceJson))
	return invoice, nil
}

func (r *repository) Delete(invoice *entity.Invoice) (*entity.Invoice, error) {
	invoiceJson, err := json.Marshal(invoice)
	if err != nil {
		return invoice, err
	}

	log.Printf("Invoice deleted: %s\n", string(invoiceJson))
	return invoice, nil
}

func (r *repository) Shutdown() {
	log.Printf("Last invoice ID: %d\n", r.invoiceIDSeq)
}
