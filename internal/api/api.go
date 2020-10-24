package api

import (
	"net/http"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
)

const MaxJSONBodySize = 2097152

// Mount initialises the API
func Mount(srv *mattrax.Server) {
	r := srv.Router.PathPrefix("/api").Subrouter()
	r.Use(Headers(srv))
	r.Use(mux.CORSMethodMiddleware(r))

	r.HandleFunc("/login", Login(srv)).Methods(http.MethodPost, http.MethodOptions)

	rAuthed := r.PathPrefix("/").Subrouter()
	rAuthed.Use(RequireAuthentication(srv))

	rAuthed.HandleFunc("/devices", Devices(srv)).Methods(http.MethodGet, http.MethodOptions)
	rAuthed.HandleFunc("/device/{id}", Device(srv)).Methods(http.MethodGet, http.MethodOptions)
	rAuthed.HandleFunc("/device/{id}/info", DeviceInformation(srv)).Methods(http.MethodGet, http.MethodOptions)
	rAuthed.HandleFunc("/device/{id}/scope", DeviceScope(srv)).Methods(http.MethodGet, http.MethodOptions)
	rAuthed.HandleFunc("/policies", Policies(srv)).Methods(http.MethodGet, http.MethodOptions)
	rAuthed.HandleFunc("/policy/{id}", Policy(srv)).Methods(http.MethodGet, http.MethodOptions)
	rAuthed.HandleFunc("/users", Users(srv)).Methods(http.MethodGet, http.MethodOptions)
	rAuthed.HandleFunc("/user/{upn}", User(srv)).Methods(http.MethodGet, http.MethodOptions)
}
