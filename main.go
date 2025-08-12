package main

import (
	"blogapi/database"
	"blogapi/middlewares"
	"blogapi/models"
	"blogapi/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

// global db handle, it maintains a connection pool of database connections
var GlobalDb *sql.DB

// required struct types
type server struct {
	store *database.Store
}

type user struct {
	Username string `json:"username"`
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

	var newUser models.CreateUser
	json.NewDecoder(r.Body).Decode(&newUser)
	passworHash, _ := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
	newUser.PasswordHash = string(passworHash) //we are using the original password and passwordHash with the same name passwordHash
	err := server.store.RegisterUser(newUser)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"userRegistered" : true}`))

}

// handles getting specific articles by id
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

func (server *server) handleGetAllArticles(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	err := server.store.StreamAllArticles(w)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (server *server) handleUserLogin(w http.ResponseWriter, r *http.Request) {

	var user user

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return
	}
	userRegistered, err := server.store.VerifyUserRegistered(user.Username)
	if err != nil {
		fmt.Println(err)
		return
	}

	if userRegistered {
		token, err := utils.GenerateJwt(user.Username)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"failure in generating the token"}`))
			return
		}

		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(15 * time.Minute),
			HttpOnly: true,
			Path:     "/",
		}

		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message":"user logged in"}`))

	} else {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"user login failed"}`))
	}
}

func (server *server) handleEditArticle(w http.ResponseWriter, r *http.Request) {

	var article models.Article

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = server.store.RegisterEditedArticle(article)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"updated":true}`))

}

func (server *server) handleDeleteArticle(w http.ResponseWriter, r *http.Request) {

	var article models.Article

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = server.store.DeleteArticle(article)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusGone)
	w.Write([]byte(`{"deleted":true}`))

}

func main() {

	r := chi.NewRouter()

	store, err := database.CreateDatabaseStore()
	if err != nil {
		fmt.Println(err)
		return
	}

	var server server = server{store}

	r.Route("/api", func(r chi.Router) {

		r.Post("/login", server.handleUserLogin)
		r.Post("/signup", server.handleUserRegistration)
		r.Get("/articles", server.handleGetAllArticles)

		r.Group(func(r chi.Router) {

			r.Use(middlewares.AuthenticationMiddleware)

			r.Get("/articles/{username}", server.handleGetArticlesByUser)
			r.Post("/articles", server.handlecreateArticle)
			r.Patch("/articles", server.handleEditArticle)
			r.Delete("/articles", server.handleDeleteArticle)

		})

	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}

}
