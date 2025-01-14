package routes

import (
	"api-gateway/proxy"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	// Auth service routes
	r.HandleFunc("/login", proxy.ProxyRequest("http://auth-service:8081/login")).Methods("POST")
	r.HandleFunc("/validate", proxy.ProxyRequest("http://auth-service:8081/validate")).Methods("GET")

	// Registration service routes
	r.HandleFunc("/register", proxy.ProxyRequest("http://registration-service:8082/register")).Methods("POST")

}
