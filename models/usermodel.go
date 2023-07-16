package models
import(
	_"time"
)
type User struct{
	Number int `json:"-"`
	ID_User string `json:",omitempty"`
	Name string 
	Balance int `json:",omitempty"`
	Phone int
	Email string `json:"-"`
	Password string `json:"-"`
}

type Signup struct{
	Name string				`json:"name" binding:"required"`
	Balance int				`json:"deposit"`
	Phone int				`json:"phone" binding:"required"`
	Email string			`json:"email" binding:"required"`
	Password string			`json:"password" binding:"required"`
	PasswordConfirm string	`json:"passwordConfirm" binding:"required"`
}

type Display struct{
	ID_Invoice string
}

type Regis struct{
	Name string
	Email string
}

type Deposit struct{
	ID_User string
	Balance int
	Phone int
}

type Login struct{
	Email string	`json:"email"`
	Password string		`json:"password`
}

type WithDraw struct{
	Balance int
}

type Logout struct{
	Email string	`json:"email"`
	Password string		`json:"password`
	Token string `json:"-"`
}