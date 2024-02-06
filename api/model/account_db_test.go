package model

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/jmoiron/sqlx"
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
	IsActive: strconv.FormatBool(false),
}

func connectToDB(t *testing.T) *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestInsertUserData(t *testing.T) {
	// open connection to database
	db := connectToDB(t)

	querynator := Querynator{TableName: "account"}
	err := querynator.Insert(newUser, db.DB)

	if err != nil {
		t.Error(err)
		return
	}

	var exist bool

	exist, err = querynator.IsExists(newUser, db.DB)

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

func TestUpdateUserData(t *testing.T) {
	db := connectToDB(t)

	querynator := Querynator{TableName: "account"}
	// err = querynator.Insert(newUser, db)

	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	var exist bool

	// exist, err = querynator.IsExists(newUser, db)

	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// if !exist {
	// 	t.Errorf("Expected true received false")
	// 	return
	// }

	// update user email
	userUpdatedFields := &Account{
		Email:    "newemail@new.domain.com",
		IsActive: strconv.FormatBool(true),
	}

	err := querynator.Update(
		userUpdatedFields,
		[]string{"username"},
		[]any{newUser.Username},
		db.DB,
	)

	if err != nil {
		t.Error(err)
		return
	}

	exist, err = querynator.IsExists(newUser, db.DB)

	if err != nil {
		t.Error(err)
		return
	}

	if exist {
		t.Errorf("There are an updated fields, expected to be false. received")
		return
	}

	t.Log("Success")
}

func TestFind(t *testing.T) {
	var err error

	db := connectToDB(t)
	accounts := []interface{}{
		&Account{Username: "ab", Email: "test1@gmail.com", Name: "Lucas", Password: "abc", IsActive: strconv.FormatBool(false)},
		&Account{Username: "abc", Email: "tes2@gmail.com", Name: "Lucas2", Password: "abc", IsActive: strconv.FormatBool(false)},
		&Account{Username: "ab", Email: "test3@gmail.com", Name: "Lucas3", Password: "abc", IsActive: strconv.FormatBool(false)},
		&Account{Username: "abde", Email: "test4@gmail.com", Name: "Lucas4", Password: "abc", IsActive: strconv.FormatBool(false)},
	}

	querynator := Querynator{TableName: "account"}
	for _, account := range accounts {
		err = querynator.Insert(account, db.DB)
		if err != nil {
			t.Error(err)
		}
	}

	result := []Account{}

	err = querynator.Find(&Account{IsActive: strconv.FormatBool(false)}, &result, 2, db, "account_id", "email", "name")

	if err != nil {
		t.Error(err)
		return
	}

	if len(result) != 2 {
		t.Errorf("Expected slice of length 2 received %d", len(result))
		return
	}

	t.Log(result)
}

func TestTest(t *testing.T) {
	db := connectToDB(t)

	dest := Account{}
	db.Get(&dest, "SELECT account_id,email,name FROM account WHERE is_active=$1 LIMIT 2", strconv.FormatBool(false))

	t.Log(dest)

}

func TestFindOne(t *testing.T) {
	var err error

	db := connectToDB(t)

	accounts := []interface{}{
		&Account{Username: "you_are_geh", Email: "gehemail@gmail.com", Name: "Lucas Geh", Password: "abc", IsActive: strconv.FormatBool(false)},
	}

	querynator := Querynator{TableName: "account"}
	for _, account := range accounts {
		err = querynator.Insert(account, db.DB)
		if err != nil {
			t.Error(err)
		}
	}

	result := &Account{}

	err = querynator.FindOne(&Account{Email: "gehemail@gmail.com"}, result, db, "account_id", "username", "email", "name")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)

	t.Log(len(result.Username))

	if result.Username != "you_are_geh" {
		t.Errorf("Invalid query result, expected username to be you_are_geh received %s", result.Username)
		return
	}

	t.Log("Success")
}

func TestDeleteAccount(t *testing.T) {
	db := connectToDB(t)

	querynator := Querynator{TableName: "account"}

	user := &Account{
		Username: "my_guy_is_geh69",
		Email:    "guygeh69@gmail.com",
		Name:     "Geh Person",
		Password: "696969696",
		IsActive: "True",
	}

	err := querynator.Insert(user, db.DB)

	if err != nil {
		t.Error(err)
		return
	}

	err = querynator.Delete(&Account{Username: "my_guy_is_geh69"}, db.DB)

	if err != nil {
		t.Error(err)
		return
	}

	var exist bool
	exist, err = querynator.IsExists(user, db.DB)

	if exist {
		t.Error("Expected false when check deleted data exist")
		return
	}

	t.Log("Success")
}
