package main

import (
	"fmt"
	"net/http"

	"github.com/ganeshbabugb/my-go-restify/internal/response"
	"github.com/go-chi/chi/v5"
)

// TODO:GB IMPLEMENT THIS
func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	fmt.Printf("userIDStr: %v\n", userIDStr)
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// TODO:GB IMPLEMENT THIS
func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	fmt.Printf("userIDStr: %v\n", userIDStr)
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func (app *application) Me(w http.ResponseWriter, r *http.Request) {
	user, _ := contextGetAuthenticatedUser(r)
	response.JSONWithHeaders(w, http.StatusOK, user, nil)
}
