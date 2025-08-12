package database

import (
	"blogapi/models"
	"database/sql"
	"encoding/json"
	"errors"
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

func (store *Store) StreamArticlesByUser(w io.Writer, username string) error {

	var userId uint64
	var article models.Article

	row := store.db.QueryRow("select id from users where username = $1", username)

	err := row.Scan(&userId)
	if err != nil {
		return err
	}


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

	_, err := store.db.Exec("insert into users(username, passwordhash) values ($1, $2)", user.Usename, user.PasswordHash)
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

		err = rows.Scan(&article.Id,&article.Title, &article.Body, &article.AuthorId)
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

func (store *Store) VerifyUserRegistered(username string) (bool, error) {

	var placeholder string

	row := store.db.QueryRow("select username from users where username = $1", username)
	err := row.Scan(&placeholder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}

func (store *Store) RegisterEditedArticle(article models.Article) error {

	_,err := store.db.Exec("update articles set title = $1,body = $2 where id = $3", article.Title, article.Body, article.Id)
	if err != nil {
		return err
	}
	return nil;
}


func (store* Store) DeleteArticle(article models.Article) error {

	_,err := store.db.Exec("delete from articles where id = $1", article.Id)
	if err != nil {
		return err
	}
	return nil;
}
