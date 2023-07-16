package routes

import(
	"github.com/gorilla/mux"
	"assignment/controllers"
	"assignment/middleware"
)

func UserRoutes(r *mux.Router){
	router:=r.PathPrefix("/users").Subrouter()

	router.Use(middleware.Auth)
	router.HandleFunc("/deposit", controllers.Deposit).Methods("POST")
	router.HandleFunc("/createInvoice", controllers.CreateInvoice).Methods("POST")
	router.HandleFunc("/transfer", controllers.Transfer).Methods("POST")

	router.HandleFunc("/getInvoice", controllers.GetInvoice).Methods("POST")
	router.HandleFunc("/withDraw", controllers.WithDraw).Methods("POST")
	router.HandleFunc("/logout", controllers.Logout).Methods("POST")
}