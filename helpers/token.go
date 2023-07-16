package helpers

import(
	"github.com/golang-jwt/jwt/v5"
	"assignment/models"
	"time"
	"fmt"
)

var mySigningKey=[]byte("mysecretkey")

type MyCustomClaims struct{
	ID_User string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone int `json:"phone"`
	jwt.RegisteredClaims
}

func CreateToken(user *models.User)(string,error){
	claims:=MyCustomClaims{
		user.ID_User,
		user.Name,
		user.Email,
		user.Phone,
		jwt.RegisteredClaims{
			ExpiresAt:jwt.NewNumericDate(time.Now().Add(100 * time.Minute)),
			IssuedAt:jwt.NewNumericDate(time.Now()),
			NotBefore:jwt.NewNumericDate(time.Now()),
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	ss,err:=token.SignedString(mySigningKey)
	return ss,err
}

func ValidateToken(tokenString string)(any,error){
	token,err:=jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token)(interface{},error){
		return mySigningKey,nil
	})
	if err != nil {
		return nil,fmt.Errorf("Unauthorized")
	}
	claims,ok:=token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid{
		return nil,fmt.Errorf("Unauthorized")
	}
	return claims,nil
}