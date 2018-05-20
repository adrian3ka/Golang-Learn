package main

import (
	"log"
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
)

func init() {
	fmt.Println("Init Main")
	// Assign a user store
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = store

	// Assign a session store
	sessionStore, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	globalSessionStore = sessionStore

	// Assign a sql database
	db, err := NewMySQLDB("root:@tcp(127.0.0.1:3306)/gophr")
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db

	// Assign an image store
	globalImageStore = NewDBImageStore()
}

func NewApp() Middleware{
	router := NewRouter()

	router.Handle("GET", "/", HandleHome)
	
	// Add the route handler
	router.Handle("POST", "/register", HandleUserCreate) //handle_user.go
	router.Handle("GET", "/register", HandleUserNew) //handle_user.go
	
	router.Handle("GET", "/login", HandleSessionNew)
	router.Handle("POST", "/login", HandleSessionCreate)
	router.Handle("GET", "/image/:imageID", HandleImageShow)
	router.Handle("GET", "/user/:userID", HandleUserShow)
	
	
	router.ServeFiles(
		"/im/*filepath",
		http.Dir("data/images/"),
	)
	
	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)
	
	
	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/sign-out", HandleSessionDestroy)
	secureRouter.Handle("GET", "/account", HandleUserEdit)
	secureRouter.Handle("POST", "/account", HandleUserUpdate)
	secureRouter.Handle("GET", "/images/new", HandleImageNew)
	secureRouter.Handle("POST", "/images/new", HandleImageCreate)



	middleware := Middleware{}
	middleware.Add(router)
	
	middleware.Add(http.HandlerFunc(RequireLogin))
	middleware.Add(secureRouter)

	return middleware
}

// Creates a new router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	return router
}


func main(){
	log.Fatal(http.ListenAndServe(":3000", NewApp()))
}