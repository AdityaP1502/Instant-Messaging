package model

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
)

// table name is account

var tableName string = "account"

type Account struct {
	Username string `json:"username" column:"username"`
	Name     string `json:"name" column:"name"`
	Email    string `json:"email" column:"email"`
	Password string `json:"password" column:"password"`
	IsActive bool   `json:"-" column:"is_active"`
}

func (acc *Account) FromJSON(r io.Reader, checkRequired bool) error {
	err := util.DecodeJSONBody(r, acc)

	if err != nil {
		return err
	}

	if checkRequired {
		return util.CheckParametersUnity(acc)
	}

	return nil
}

func (acc *Account) ToJSON(checkRequired bool) ([]byte, error) {
	var err error

	var tmp struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}

	tmp.Username = acc.Username
	tmp.Name = acc.Name
	tmp.Email = acc.Email

	if checkRequired {
		if err = util.CheckParametersUnity(&tmp); err != nil {
			return nil, err
		}
	}

	return util.CreateJSONResponse(&tmp)
}

func (acc *Account) Insert(db *sql.DB) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (username, name, email, password, is_active)
		VALUES($1, $2, $3, $4, $5)`, tableName,
	)

	_, err := db.Exec(query, acc.Username, acc.Name, acc.Email, acc.Password, acc.IsActive)

	return err
}

func (acc *Account) Update(db *sql.DB, conditionColumns []string, conditionsValues []any) error {
	keys, values := getNonEmptyField(acc)

	updateFieldsString := transformNamesToUpdateQuery(keys, 1, ",")
	conditionFieldsString := transformNamesToUpdateQuery(conditionColumns, len(keys)+1, " AND ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, updateFieldsString, conditionFieldsString)

	values = append(values, conditionsValues...)

	_, err := db.Exec(query, values...)

	return err
}

func (acc *Account) IsExists(db *sql.DB) (bool, error) {
	//https://stackoverflow.com/questions/32554400/efficiently-determine-if-any-rows-satisfy-a-predicate-in-postgres?rq=3

	var exists bool

	keys, values := getNonEmptyField(acc)
	conditionString := transformNamesToUpdateQuery(keys, 1, " AND ")

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s)", tableName, conditionString)

	fmt.Println(query)

	err := db.QueryRow(query, values...).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
