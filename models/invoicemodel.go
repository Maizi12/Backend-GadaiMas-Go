package models

import(
	
)
type Transfer struct{
	ID_Invoice string
	Payment int
}
type Invoice struct{
	Number int
	ID_Invoice string
	ID_User string
	Carid string
	LeasingID string
	Hutang int
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}
type Inv_out struct{
	InvoiceId string
	CustomerId string
	LeasingID string
	CarId string
	Customer User
	Leasing Leasing
	Car Car
	LoanPrinciple int
	LoanTotal int
	Term int
	AmountsPaid int		`json:",omitempty"`
	MissedAmounts int	`json:",omitempty"`
	NextPayment int		
	PaymentsDue string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}