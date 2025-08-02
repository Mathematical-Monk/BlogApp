package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"blogapi/database"
	"blogapi/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// global db handle, it maintains a connection pool of database connections
var GlobalDb *sql.DB

// required struct types
type server struct {
	store *database.Store
}

type user struct {
	Username     string `json:"username"`
	Password string `json:"passwordHash"`
}

type article struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// handles the creation of new articles
func (server *server) handlecreateArticle(w http.ResponseWriter, r *http.Request) {

	var article models.Article

	json.NewDecoder(r.Body).Decode(&article)
	err := server.store.CreateArticleInDb(article)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"articleCreated":true}`))
}

// handles the registration of a new user
func (server *server) handleUserRegistration(w http.ResponseWriter, r *http.Request) {

	var newUser user
	json.NewDecoder(r.Body).Decode(&newUser)
	passworHdash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	err := server.store.RegisterUser(newUser.Username, string(passworHdash))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"userRegistered" : true}`))

}

//handles getting specific articles by id
func (server *server) handleGetArticlesByUser(w http.ResponseWriter, r *http.Request) {

	username := chi.URLParam(r, "username")
	fmt.Println(username)
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	err := server.store.StreamArticlesByUser(w, username)
	if err != nil {
		fmt.Println(err)
		return
	}

}

// func (server *server) handleGetAllArticles(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-type", "application/json")

// 	err := server.store.StreamAllArticles(w)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

func (server* server) handleUserLogin(w http.ResponseWriter, r* http.Request){

	var user user

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return
	}
	userRegistered,err := server.store.VerifyUserRegistered(user.Username)
	if err != nil {
		fmt.Println(err)
		return
	}

	if userRegistered {
		//token logic when user is registered

	}else{
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"userRegistered":false}`))
	}
}



func main() {

	r := chi.NewRouter()

	store, err := database.CreateDatabaseStore()
	if err != nil {
		fmt.Println(err)
		return
	}

	var server server = server{store}

	r.Get("/api/articles/{username}", server.handleGetArticlesByUser)
	// r.Get("/api/articles", server.handleGetAllArticles)

	r.Post("/api/articles", server.handlecreateArticle)
	r.Post("/api/signup", server.handleUserRegistration)

	r.Post("/api/login", server.handleUserLogin)

	// r.Patch("/api/articles", server.handleEditArticle)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}

}
