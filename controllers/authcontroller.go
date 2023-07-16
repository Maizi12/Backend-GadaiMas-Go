package controllers

import(
	"net/http"
	"assignment/models"
	"encoding/json"
	"assignment/config"
	"assignment/helpers"
	"fmt"
	_"strconv"
	_"strings"

)

func Signup(w http.ResponseWriter, r *http.Request){
	var signup models.Signup

	if err:=json.NewDecoder(r.Body).Decode(&signup);err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}
	defer r.Body.Close()

	if signup.Password!=signup.PasswordConfirm{
		helpers.Response(w,400,"Password not match",nil,nil)
		return
	}

	passwordHash,err:=helpers.HashPassword(signup.Password)

	if err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
	}

	user:=models.User{
		ID_User: Numbering("users", "CUST"),
		Name: signup.Name,
		Balance: signup.Balance,
		Phone:signup.Phone,
		Email: signup.Email,
		Password: passwordHash,
	}

	user_dis:=models.Regis{
		Name:user.Name,
		Email:signup.Email,
	}

	if err:=config.DB.Create(&user).Error;err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}

	helpers.Response(w,201,"Register successfully",user_dis,nil)

}

func Login(w http.ResponseWriter, r *http.Request){
	var login models.Login

	if err:=json.NewDecoder(r.Body).Decode(&login);err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}

	var user models.User
	if err:=config.DB.First(&user,"email=?",login.Email).Error;err != nil {
		helpers.Response(w,404,"Wrong Email",nil,nil)
		return 
	}

	if err:=helpers.VerifyPassword(user.Password,login.Password);err != nil {
		helpers.Response(w,404,"Wrong Password",nil,nil)
		return 
	}

	token,err:=helpers.CreateToken(&user)

	if err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}
	
	helpers.Response(w,200,"Successfully Login",user,token)
	fmt.Println(token)
}

func Logout(w http.ResponseWriter, r *http.Request){
	var logout models.Logout
	//users:=r.Context().Value("userinfo").(*helpers.MyCustomClaims)

	if err:=json.NewDecoder(r.Body).Decode(&logout);err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}

	var user models.User
	if err:=config.DB.First(&user,"email=?",logout.Email).Error;err != nil {
		helpers.Response(w,404,"Wrong Email",nil,nil)
		return 
	}

	if err:=helpers.VerifyPassword(user.Password,logout.Password);err != nil {
		helpers.Response(w,404,"Wrong Password",nil,nil)
		return 
	}

	
	helpers.Response(w,200,"Thank you "+user.Name+", See you next time!",nil,nil)
	fmt.Println(logout.Token)
}