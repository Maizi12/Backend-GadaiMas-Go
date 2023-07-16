package main

import(
	"assignment/config"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"assignment/routes"
)

func main(){
	
	config.ConnectDB()
	r:=mux.NewRouter()
	router:=r.PathPrefix("/").Subrouter()
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", router)
	
}