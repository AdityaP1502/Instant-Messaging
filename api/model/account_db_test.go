package model

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 3000
	user     = "instant"
	password = "4jBWgQ7qpmYq19+0y07Gc/VAts4QyBKrv1/UeORklQc="
	dbname   = "account_db"
)

var newUser = &Account{
	Username: "aditya",
	Name:     "I Made Aditya",
	Email:    "aditya@example.com",
	Password: "hello,world",
	IsActive: false,
}

func TestInsertUserData(t *testing.T) {
	// open connection to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	err = newUser.Insert(db)

	if err != nil {
		t.Error(err)
		return
	}

	var exist bool

	exist, err = newUser.IsExists(db)

	if err != nil {
		t.Error(err)
		return
	}

	if !exist {
		t.Errorf("Expected true received false")
		return
	}

	t.Log("Success")
}
