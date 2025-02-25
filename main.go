package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"log"
	"net/http"
	"time"
)

// User struct to store form data
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name,omitempty"`
	Email     string             `bson:"email" json:"email,omitempty"`
	Password  string             `bson:"password" json:"password,omitempty"`
	IsNew     bool               `bson:"isNew" json:"isNew,omitempty"` // Determines if it's a new account creation
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type UserDTO struct {
	ID   string
	Name string
}

// Global variable to store parsed templates
var tmpl *template.Template
var usersCollection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	var ctx = context.TODO()

	var err error
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("connection failed", err)
	}

	db := client.Database("forum")
	usersCollection = db.Collection("users")
	log.Println("Connected to MongoDB!")

	tmpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	var err error

	router := mux.NewRouter()
	router.HandleFunc("/", showDefault)
	router.HandleFunc("/register", showRegisterForm)
	router.HandleFunc("/login", showLoginForm).Methods(http.MethodGet)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	router.HandleFunc("/submit", registrationHandler)
	router.HandleFunc("/home", showHome)
	router.HandleFunc("/update/{userID}", updateHandler).Methods(http.MethodPost) // POST
	router.HandleFunc("/update/{userID}", showUpdateForm).Methods(http.MethodGet) //GET

	// Route to handle form submission
	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Server starting at port 8080")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
