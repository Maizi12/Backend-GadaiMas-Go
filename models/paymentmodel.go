package models

import(
	_"time"
)

type Payment struct{
	Number	int
	ID_Invoice	string
	ID_Payment	string
	Nominal int
	MissedAmounts int
	CreatedAt	string
} 

type PayAmount struct{
	Number int
	Nominal int
	MissedAmounts int
}