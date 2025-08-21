package database

import (
	"blogapi/models"
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	// "errors"
	"fmt"
	"io"
)

type Store struct {
	db *sql.DB
}

func CreateDatabaseStore() (*Store, error) {
	store := &Store{}
	DBhandle, err := sql.Open("pgx", "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, err
	}

	err = DBhandle.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("database connection established")

	store.db = DBhandle
	return store, nil
}

func (store *Store) StreamArticlesByUser(w io.Writer, userId int64) error {

	var article models.Article

	rows, err := store.db.Query("select id, title, body, author_id from articles where author_id = $1", userId)
	if err != nil {
		return err
	}

	defer rows.Close()

	//start creating the json array stream manually
	_, err = w.Write([]byte("["))
	if err != nil {
		return err
	}

	var isFirstRow bool = true

	for rows.Next() {

		if !isFirstRow {
			_, err = w.Write([]byte(","))
			if err != nil {
				return err
			}
		}
		isFirstRow = false

		err = rows.Scan(&article.Id, &article.Title, &article.Body, &article.AuthorId)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			return err
		}

	}
	_, err = w.Write([]byte("]"))
	if err != nil {
		return err
	}

	return rows.Err()

}

func (store *Store) RegisterUser(user models.CreateUser) error {

	_, err := store.db.Exec("insert into users(username, passwordhash) values ($1, $2)", user.UserName, user.PasswordHash)
	if err != nil {
		return err
	}

	return nil

}

func (store *Store) CreateArticleInDb(article models.Article) error {

	_, err := store.db.Exec("insert into articles (title, body, author_id) values ($1, $2, $3)", article.Title, article.Body, article.AuthorId)
	if err != nil {
		return err
	}
	return nil

}

func (store *Store) StreamAllArticles(w io.Writer) error {

	rows, err := store.db.Query("select id, title, body, author_id from articles")
	if err != nil {
		return err
	}

	defer rows.Close()

	var article models.Article

	isFirstRow := true

	_, err = w.Write([]byte("["))
	if err != nil {
		return err
	}

	for rows.Next() {

		if !isFirstRow {
			_, err = w.Write([]byte(","))
			if err != nil {
				return err
			}
		}
		isFirstRow = false

		err = rows.Scan(&article.Id, &article.Title, &article.Body, &article.AuthorId)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte("]"))
	if err != nil {
		return err
	}

	return rows.Err()
}

func (store *Store) VerifyUserRegistered(username string, password string) (int64, string, error) {

	var userId int64
	var passwordHash string

	row := store.db.QueryRow("select id, passwordHash from users where username = $1", username)
	err := row.Scan(&userId, &passwordHash)
	if err != nil {
		return -1, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return -1, "", err
	}

	return userId, passwordHash, nil

}

func (store *Store) RegisterEditedArticle(article models.Article) error {

	_, err := store.db.Exec("update articles set title = $1,body = $2 where id = $3", article.Title, article.Body, article.Id)
	if err != nil {
		return err
	}
	return nil
}

func (store *Store) DeleteArticle(article models.Article) error {

	_, err := store.db.Exec("delete from articles where id = $1", article.Id)
	if err != nil {
		return err
	}
	return nil
}

func (store *Store) CheckAndEditArticle(article models.Article) (bool, error) {

	rowsInfo, err := store.db.Exec("update articles set body = $1, title = $2 where id = $3 and author_id = $4", article.Body, article.Title, article.Id, article.AuthorId)
	if err != nil {
		return false, err
	}

	rowsAffected, err := rowsInfo.RowsAffected()
	if err != nil {
		return false, nil
	}

	if rowsAffected == 0 {
		return false, nil
	}else {
		return true, nil
	}
}
