package entity

import (
	"time"
)

type Invoice struct {
	ID                int64
	Code              string
	Buyer             string
	Phone             string
	Address           string
	Orders            []Order
	Courier           string
	ShippingFee       int64
	ActualShippingFee int64
	COD               bool
	FreeShippingFee   bool
	Price             int64
	PO                bool
	PODurationWeek    int
	MinDP             int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
