package controllers

import(
	"net/http"
	"assignment/helpers"
	"encoding/json"
	"assignment/models"
	"assignment/config"
	_"context"
	"strings"
	"strconv"
	"time"
	_"fmt"
)

	var inv models.Invoice
	var user models.User
	var car models.Car
	var lease models.Leasing
	var trf models.Transfer
	var Next_Payment int
	var missedAmounts int
	var pay_count int
	var display models.Display
	var pay_time models.PayAmount

func InvOut(inv models.Invoice,user models.User,lease models.Leasing,car models.Car,Loan int,loanTotal int,Next_Payment int, keyword string)models.Inv_out {
	var paymentsDue string
	var createdAt string
	var updatedAt string
	var id_inv string
	var newTerms int
	
	config.DB.Raw("SELECT COUNT(id_invoice) FROM payments where id_invoice=?",inv.ID_Invoice).Scan(&pay_count)
	now:=time.Now()
	
	paymentsDue=now.AddDate(0,1,0).Format("2006-01-02 15:04:05.000")
	updatedAt=now.Format("2006-01-02 15:04:05.000")
	
	if keyword=="create"{
	paymentsDue=now.AddDate(0,1,0).Format("2006-01-02 15:04:05.000")
	createdAt=now.Format("2006-01-02 15:04:05.000")
	newTerms=lease.Terms
	id_inv=Numbering("invoices", "INV")
	}

	if keyword=="pay"{
	paymentsDue=now.AddDate(0,1,0).Format("2006-01-02 15:04:05.000")
	createdAt=inv.CreatedAt
	updatedAt=now.Format("2006-01-02 15:04:05.000")
	id_inv=inv.ID_Invoice
	newTerms=lease.Terms-1-pay_count
	}

	if keyword=="display"{
	createdAt=inv.CreatedAt
	updatedAt=inv.UpdatedAt
	t1, _ := time.Parse("2006-01-02 15:04:05.000",updatedAt)
	t2 := t1.AddDate(0, 1, 0)
	paymentsDue=t2.Format("2006-01-02 15:04:05.000")
	id_inv=inv.ID_Invoice
	newTerms=lease.Terms-pay_count
	}

	inv_out:=models.Inv_out{
		InvoiceId:id_inv,
		CustomerId:user.ID_User,
		LeasingID:inv.LeasingID,
		CarId:inv.Carid,
		Customer: models.User{
			Name:user.Name,
			Phone:user.Phone,
		},
		Leasing: models.Leasing{
			LeasingID:lease.LeasingID,
			LeasingName: lease.LeasingName,
			Rates: (lease.Rates*100),
		},
		Car: models.Car{
			BrandName: car.BrandName,
			GroupModelName: car.GroupModelName,
			ModalName: car.ModalName,
			Year:car.Year,
			Price:car.Price,
			CreatedAt: car.CreatedAt,
		},
		LoanPrinciple: Loan,
		LoanTotal: int(loanTotal),
		Term: newTerms,
		NextPayment: Next_Payment,
		PaymentsDue: paymentsDue,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt:"",
	}
	return inv_out
}

func Numbering(tbl_name, keyword string)(string){
	var iterate string
	var invid string
	config.DB.Raw("SELECT number FROM "+tbl_name+" ORDER BY number DESC LIMIT 1").Scan(&invid)
	inv_id,_:=strconv.Atoi(strings.Trim(invid,keyword))
	
	if inv_id>=10 && inv_id<100{
		iterate="0"
	}
	if inv_id>=100{
			iterate=""
	}
	if inv_id<10{
		iterate="00"
	}
	id_inv:=keyword+iterate+strconv.Itoa(inv_id+1)
	return id_inv
}

func Deposit(w http.ResponseWriter, r *http.Request){
	var users models.User
	user:=r.Context().Value("userinfo").(*helpers.MyCustomClaims)
	if err:=json.NewDecoder(r.Body).Decode(&users);err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}
	defer r.Body.Close()
	userinf:=&models.Deposit{
		ID_User: user.ID_User,
		Balance: users.Balance,
		Phone:	user.Phone,
	}
	config.DB.Exec("Update users SET Balance=? WHERE ID_User=?", userinf.Balance, userinf.ID_User)
	helpers.Response(w,200,"Deposit",userinf,nil)
}

func WithDraw(w http.ResponseWriter, r *http.Request){
	var withDraw models.WithDraw
	var user models.User
	users:=r.Context().Value("userinfo").(*helpers.MyCustomClaims)
	if err:=json.NewDecoder(r.Body).Decode(&withDraw);err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}
	config.DB.Raw("SELECT * from users WHERE id_user=?",users.ID_User).Scan(&user)
	user.Balance=user.Balance+(withDraw.Balance*-1)
	config.DB.Exec("Update users SET Balance=? WHERE ID_User=?", user.Balance, user.ID_User)
	helpers.Response(w,200,"WithDraw",user,nil)
}

func CreateInvoice(w http.ResponseWriter, r *http.Request){

	users:=r.Context().Value("userinfo").(*helpers.MyCustomClaims)

	id_inv:=Numbering("invoices", "INV")
	if err:=json.NewDecoder(r.Body).Decode(&inv);err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}

	user=models.User{
		ID_User: users.ID_User,
		Name:users.Name,
		Phone:users.Phone,
	}

	config.DB.Raw("SELECT balance from users WHERE id_user=?",user.ID_User).Scan(&user.Balance)
	config.DB.Raw("SELECT brand_name, group_model_name, modal_name, year, price, created_at from cars WHERE carid=?",inv.Carid).Scan(&car)
	config.DB.Raw("SELECT LeasingName, rates, terms FROM leasings WHERE leasing_id=?",inv.LeasingID).Scan(&lease)
	defer r.Body.Close()

	Loan:=car.Price-user.Balance
	loanTotal:=int((Loan+(int(lease.Rates*100)*(lease.Terms/12)*(Loan))/100))

	Next_Payment:=int(loanTotal)/lease.Terms
	lease.LeasingID=inv.LeasingID
	if car.Price<user.Balance{
		Balance:=user.Balance-car.Price
		Loan=0
		loanTotal=0
		Next_Payment=0
		lease.Terms=0
		config.DB.Exec("Update users SET Balance=? WHERE ID_User=?", Balance, user.ID_User)
	}
	
	id_inv=Numbering("invoices", "INV")
	
	inv_out:=InvOut(inv,user,lease,car,Loan,loanTotal,Next_Payment,"create")

	Create_Inv:=models.Invoice{
	ID_Invoice :id_inv,
	ID_User :users.ID_User,
	Carid :inv.Carid,
	LeasingID :inv.LeasingID,
	Hutang : Loan,
	CreatedAt :time.Now().Format("2006-01-02 15:04:05.000"),
	UpdatedAt :time.Now().Format("2006-01-02 15:04:05.000"),
	}

	if err:=config.DB.Create(&Create_Inv).Error;err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}

	helpers.Response(w,201,"Register successfully",inv_out,nil)

}

func Transfer(w http.ResponseWriter, r *http.Request){

	users:=r.Context().Value("userinfo").(*helpers.MyCustomClaims)

	if err:=json.NewDecoder(r.Body).Decode(&trf);err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}
	
	user=models.User{
		ID_User: users.ID_User,
		Name:users.Name,
		Phone:users.Phone,
	}

	config.DB.Raw("SELECT * FROM INVOICES WHERE id_invoice=?", trf.ID_Invoice).Scan(&inv)
	config.DB.Raw("SELECT balance from users WHERE id_user=?",user.ID_User).Scan(&user.Balance)
	config.DB.Raw("SELECT carid, brand_name, group_model_name, modal_name, year, price, created_at from cars WHERE carid=?",inv.Carid).Scan(&car)
	config.DB.Raw("SELECT leasing_id,LeasingName, rates, terms FROM leasings WHERE leasing_id=?",inv.LeasingID).Scan(&lease)
	
	defer r.Body.Close()

	Loan:=inv.Hutang
	
	loanTotal:=int((Loan+(int(lease.Rates*100)*(lease.Terms/12)*(Loan))/100))
	//loanTotal:=0
	Compound_lease:=int(loanTotal)/lease.Terms
	var kurangBayar int
	config.DB.Raw("SELECT missed_amounts FROM payments WHERE id_invoice=? AND nominal<? ORDER By number DESC LIMIT 1", trf.ID_Invoice,Compound_lease).Scan(&kurangBayar)
	
	if trf.Payment<Compound_lease{
		Next_Payment=int(float32(Compound_lease+(Compound_lease-trf.Payment))+(float32(Compound_lease-trf.Payment)*0.02))+(int(float32(kurangBayar)*0.02)+kurangBayar)//5.100.000
		missedAmounts=Compound_lease-trf.Payment
	}else{
		Next_Payment=int(float32(int(loanTotal)/lease.Terms))+kurangBayar
		missedAmounts=0
	}
	lease.LeasingID=inv.LeasingID

	inv_out:=InvOut(inv,user,lease,car,Loan,loanTotal,Next_Payment,"pay")

	pay_id:=Numbering("payments","PAY")
	pay:=models.Payment{
		ID_Invoice: trf.ID_Invoice,
		ID_Payment:	pay_id,
		Nominal: trf.Payment,
		MissedAmounts: (int(float32(missedAmounts)*0.02)+missedAmounts)+int(float32(kurangBayar)*0.02)+kurangBayar,
		CreatedAt:	time.Now().Format("2006-01-02 15:04:05.000"),
	}
	if err:=config.DB.Create(&pay).Error;err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}
	helpers.Response(w,201,"Register successfully",inv_out,nil)
}

func GetInvoice(w http.ResponseWriter, r *http.Request){
	
	if err:=json.NewDecoder(r.Body).Decode(&display);err != nil {
		helpers.Response(w,500,err.Error(),nil,nil)
		return 
	}

	config.DB.Raw("SELECT * FROM INVOICES WHERE id_invoice=?", display.ID_Invoice).Scan(&inv)
	config.DB.Raw("SELECT * from users WHERE id_user=?",inv.ID_User).Scan(&user)
	config.DB.Raw("SELECT carid, brand_name, group_model_name, modal_name, year, price, created_at from cars WHERE carid=?",inv.Carid).Scan(&car)
	config.DB.Raw("SELECT leasing_id,LeasingName, rates, terms FROM leasings WHERE leasing_id=?",inv.LeasingID).Scan(&lease)
	config.DB.Raw("SELECT COUNT(number) FROM payments WHERE id_invoice=?",display.ID_Invoice).Scan(&pay_time.Number)
	config.DB.Raw("SELECT SUM(nominal) FROM payments WHERE id_invoice=?",display.ID_Invoice).Scan(&pay_time.Nominal)
	config.DB.Raw("SELECT SUM(missed_amounts) FROM payments WHERE id_invoice=?",display.ID_Invoice).Scan(&pay_time.MissedAmounts)

	defer r.Body.Close()

	Loan:=inv.Hutang

	loanTotal:=int(float32(1+float32(lease.Rates*float32(lease.Terms/12)))*float32(Loan))

	Next_Payment=int(float32(int(loanTotal)/lease.Terms))+pay_time.MissedAmounts
	missedAmounts=pay_time.MissedAmounts

	lease.LeasingID=inv.LeasingID

	inv_out:=InvOut(inv,user,lease,car,Loan,loanTotal,Next_Payment,"display")

	helpers.Response(w,201,"Register successfully",inv_out,nil)
}