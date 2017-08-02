package rest

//This file is auto-generated by go-raml
//Do not edit this file by hand since it will be overwritten during the next generation

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/zero-os/0-stor/server/db"
)

// NamespacesInterface is interface for /namespaces/{nsid} root endpoint
type NamespacesInterface interface { // UpdateReferenceList is the handler for PUT /namespaces/{nsid}/objects/{id}/references
	// Update reference list.
	// The reference list of the object will be update with the references from the request body
	UpdateReferenceList(http.ResponseWriter, *http.Request)
	// DeleteObject is the handler for DELETE /namespaces/{nsid}/objects/{id}
	// Delete object from the server
	DeleteObject(http.ResponseWriter, *http.Request)
	// GetObject is the handler for GET /namespaces/{nsid}/objects/{id}
	// Retrieve object from the server
	GetObject(http.ResponseWriter, *http.Request)
	// ListObjects is the handler for GET /namespaces/{nsid}/objects
	// List keys of the object in the namespace
	ListObjects(http.ResponseWriter, *http.Request)
	// CreateObject is the handler for POST /namespaces/{nsid}/objects
	// Set an object into the namespace
	CreateObject(http.ResponseWriter, *http.Request)
	// reservationsidGet is the handler for GET /namespaces/{nsid}/reservations/{id}
	// Return information about a reservation
	reservationsidGet(http.ResponseWriter, *http.Request)
	// UpdateReservation is the handler for PUT /namespaces/{nsid}/reservations/{id}
	// Renew an existing reservation
	UpdateReservation(http.ResponseWriter, *http.Request)
	// ListReservations is the handler for GET /namespaces/{nsid}/reservations
	// Return a list of all the existing reservation for the give resource
	ListReservations(http.ResponseWriter, *http.Request)
	// CreateReservation is the handler for POST /namespaces/{nsid}/reservations
	// Create a reservation for the given resource.
	CreateReservation(http.ResponseWriter, *http.Request)
	// GetNameSpace is the handler for GET /namespaces/{nsid}
	// Get detail view about namespace
	GetNameSpace(http.ResponseWriter, *http.Request)
}

// NamespacesInterfaceRoutes is routing for /namespaces/{nsid} root endpoint
func NamespacesInterfaceRoutes(r *mux.Router, i NamespacesInterface, db db.DB) {

	 iyo := NewOauth2itsyouonlineMiddleware(db).Handler

	 r.Handle("/namespaces/{nsid}/objects/{id}/references", alice.New(iyo).Then(http.HandlerFunc(i.UpdateReferenceList))).Methods("PUT")
	 r.Handle("/namespaces/{nsid}/objects/{id}", alice.New(iyo).Then(http.HandlerFunc(i.DeleteObject))).Methods("DELETE")
	 r.Handle("/namespaces/{nsid}/objects/{id}", alice.New(iyo).Then(http.HandlerFunc(i.GetObject))).Methods("GET")
	 r.Handle("/namespaces/{nsid}/objects", alice.New(iyo).Then(http.HandlerFunc(i.ListObjects))).Methods("GET")
	 r.Handle("/namespaces/{nsid}/objects", alice.New(iyo).Then(http.HandlerFunc(i.CreateObject))).Methods("POST")
	 r.Handle("/namespaces/{nsid}/reservations/{id}", alice.New(iyo).Then(http.HandlerFunc(i.reservationsidGet))).Methods("GET")
	 r.Handle("/namespaces/{nsid}/reservations/{id}", alice.New(iyo).Then(http.HandlerFunc(i.UpdateReservation))).Methods("PUT")
	 r.Handle("/namespaces/{nsid}/reservations", alice.New(iyo).Then(http.HandlerFunc(i.ListReservations))).Methods("GET")
	 r.Handle("/namespaces/{nsid}/reservations", alice.New(iyo).Then(http.HandlerFunc(i.CreateReservation))).Methods("POST")
	 r.Handle("/namespaces/{nsid}", alice.New(iyo).Then(http.HandlerFunc(i.GetNameSpace))).Methods("GET")
}