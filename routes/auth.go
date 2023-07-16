package routes

import(
	"github.com/gorilla/mux"
	"assignment/controllers"
)

func AuthRoutes(r *mux.Router){
	//router:=r.PathPrefix("/auth").Subrouter()
	router:=r.PathPrefix("").Subrouter()

	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	
}