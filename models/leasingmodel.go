package models

import(
	
)
type Leasing struct{
LeasingID	string //`json:"-"`
LeasingName	string
Rates		float32
Terms		int `json:"-"`
}