package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

// Function to render the login page
func showLoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// extract from the form
	// eml, pwd
	//

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
	if action != "login" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// todo validate

	filter := bson.D{
		{"email", email},
	}

	ctx := context.TODO()
	var user User
	res := usersCollection.FindOne(ctx, filter)
	if res.Err() != nil {
		tmpl.ExecuteTemplate(w, "login.gohtml", "User not found")
		return
	}

	err = res.Decode(&user)
	if err != nil {
		tmpl.ExecuteTemplate(w, "login.gohtml", "User could not be found")
		return
	}

	// compare hash
	isCorrectPassword := CheckPasswordHash(password, user.Password)
	if !isCorrectPassword {
		tmpl.ExecuteTemplate(w, "home.gohtml", "Password incorrect")
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)

}
