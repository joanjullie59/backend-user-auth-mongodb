package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

func showUpdateForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "update.gohtml", nil)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	// Check if the form submission is for account creation or login
	action := r.FormValue("action")

	// If it's a new account, validate passwords
	if action == "edit" {
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if password != confirmPassword {
			tmpl.ExecuteTemplate(w, "update.gohtml", "Passwords do not match!")
			return
		}
	}

	userIDParam := mux.Vars(r)["userID"]
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "update.gohtml", "Cannot update user")
		return
	}

	// check if user exists
	filter := bson.D{
		{"_id", userID},
	}

	ctx := context.TODO()
	count, err := usersCollection.CountDocuments(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "update.gohtml", "Cannot create user")
		return
	}

	if count != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "update.gohtml", "You dont belong here !")
		return
	}
	//hashes password for updated passwords
	pwd := r.FormValue("password")
	hashedPassword, err := HashPassword(pwd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "update.gohtml", "Cannot create user")
		return
	}

	// Create a new User instance with form data
	user := User{
		ID:        userID,
		Name:      r.FormValue("name"),
		Email:     r.FormValue("email"),
		Password:  hashedPassword,
		IsNew:     action == "create",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	filter = bson.D{
		{"email", user.Email},
	}

	userIDFilter := bson.D{
		{"_id", userID},
	}
	res, err := usersCollection.ReplaceOne(ctx, userIDFilter, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "update.gohtml", "Cannot update user")
		log.Println(err)
		return
	}

	log.Println("Registered ", res.UpsertedID)
	tmpl.ExecuteTemplate(w, "home.gohtml", user) // Show account creation success page

}
