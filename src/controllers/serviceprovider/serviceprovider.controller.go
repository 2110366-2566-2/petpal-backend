package controllers

import (
	"encoding/json"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/serviceprovider"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

// GetSVCPsHandler handles the fetching of all service providers
func GetSVCPsHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	params := r.URL.Query()

	// set default values for page and per
	if !params.Has("page") { params.Set("page", "1") }
	if !params.Has("per") { params.Set("per", "10") }

	// fetch page and per from request query
	page, err_page := strconv.ParseInt(params.Get("page"), 10, 64)
	per, err_per := strconv.ParseInt(params.Get("per"), 10, 64)
	if err_page != nil || err_per != nil{
		http.Error(w, "Failed to parse request query params", http.StatusBadRequest)
		return
	}

	// get all users, no filters for now
	svcps, err := utills.GetSVCPs(db, bson.D{}, page - 1, per)
	if err != nil {
		println(err.Error())
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(svcps)
}

// GetUserByIDHandler handles the fetching of a user by ID
func GetSVCPByIDHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB, id string) {
	user, err := utills.GetSVCPByID(db, id)
	if err != nil {
		println(err.Error())
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}