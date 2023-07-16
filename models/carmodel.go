package models

import(
	"time"
)

type Car struct{
	Number int `json:"-"`
	Carid string `json:",omitempty"`
	BrandName string
	GroupModelName string
	ModalName string
	Year	int
	Price	int
	CreatedAt	time.Time `json:,omitempty`
}