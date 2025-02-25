package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

// Function to render the create account page
func showRegisterForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "register.gohtml", nil)
}

// Function to process form submission
func registrationHandler(w http.ResponseWriter, r *http.Request) {
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
	if action == "create" {
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if password != confirmPassword {
			tmpl.ExecuteTemplate(w, "register.gohtml", "Passwords do not match!")
			return
		}
		if !isValidPassword(password) {
			tmpl.ExecuteTemplate(w, "register.gohtml", "Password must be at least 8 characters long and contain only lowercase letters and digits.")
			return
		}

	}
	//hashes password for registering users
	pwd := r.FormValue("password")
	hashedPassword, err := HashPassword(pwd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "register.gohtml", "Cannot create user")
		return
	}

	// Create a new User instance with form data
	user := User{
		ID:        primitive.NewObjectID(),
		Name:      r.FormValue("name"),
		Email:     r.FormValue("email"),
		Password:  hashedPassword,
		IsNew:     action == "create",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := context.TODO()

	filter := bson.D{
		{"email", user.Email},
	}
	//checks if the user already exists
	count, err := usersCollection.CountDocuments(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "register.gohtml", "Cannot create user")
		return
	}

	if count > 0 {
		w.WriteHeader(http.StatusUnauthorized)
		tmpl.ExecuteTemplate(w, "register.gohtml", "You're registered !!!!!!!")
		return
	}

	res, err := usersCollection.InsertOne(ctx, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "register.gohtml", "Cannot create user")
		return
	}

	userDto := UserDTO{
		ID:   user.ID.Hex(),
		Name: user.Name,
	}

	log.Println("Registered ", res.InsertedID)
	tmpl.ExecuteTemplate(w, "home.gohtml", userDto)
}
