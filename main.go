package main

import (
	"controllers"
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/rs/cors"
)

func main() {
	controllers.INITKey()

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("views/public"))))
	//Public Endpoint
	mux.HandleFunc("/", controllers.IndexHandler)

	mux.HandleFunc("/login", controllers.LoginHandler)

	//Protected Endpoints
	mux.Handle("/entitypermissions", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.FetchEntityPermissionsHandler)))

	mux.Handle("/userpermissions", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.FetchUserPermissionsHandler)))

	mux.Handle("/fetchclient", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.FetchClientHandler)))

	mux.Handle("/deleteclient", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.DeleteClientHandler)))

	mux.Handle("/fetchallclients", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.FetchAllClientsHandler)))

	mux.Handle("/saveclient", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.SaveClientHandler)))

	mux.Handle("/getclienttype", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.ClientTypeHandler)))

	mux.Handle("/fetchallpermissions", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.FetchAllPermissionsHandler)))

	mux.Handle("/savepermission", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.SavePermissionHandler)))

	mux.Handle("/fetchallroles", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.FetchAllRolesHandler)))

	mux.Handle("/savenewuser", negroni.New(
		negroni.HandlerFunc(controllers.ValidateTokenMiddleware),
		negroni.WrapFunc(controllers.SaveUserHandler)))

	// corsHandler := cors.Default().Handler(mux)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept - Encoding", "X - CSRF - Token", "Authorization"},
	}).Handler(mux)

	fmt.Println("Listening on port 1147")
	err := http.ListenAndServe(":1147", corsHandler)
	log.Printf("Server Error : %#v\n", err)
}
