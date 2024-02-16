package model

import (
	"io"

	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
	internalserviceerror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/internal_service_error"
)

type Account struct {
	AccountID string `json:"-" db:"account_id"`
	Username  string `json:"username" db:"username"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	Salt      string `json:"-" db:"password_salt"`
	IsActive  string `json:"-" db:"is_active"`
}

func (acc *Account) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) error {
	err := util.DecodeJSONBody(r, acc)

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	if checkRequired {
		return util.CheckParametersUnity(acc, requiredFields)
	}

	return nil
}

func (acc *Account) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
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
		if err = util.CheckParametersUnity(&tmp, requiredFields); err != nil {
			return nil, err
		}
	}

	return util.CreateJSONResponse(&tmp)
}

// func (acc *Account) Insert(db *sql.DB) error {
// 	query := fmt.Sprintf(
// 		`INSERT INTO %s (username, name, email, password, is_active)
// 		VALUES($1, $2, $3, $4, $5)`, tableName,
// 	)

// 	_, err := db.Exec(query, acc.Username, acc.Name, acc.Email, acc.Password, acc.IsActive)

// 	return err
// }

// func (acc *Account) Update(db *sql.DB, conditionColumns []string, conditionsValues []any) error {
// 	keys, values := getNonEmptyField(acc)

// 	updateFieldsString := transformNamesToUpdateQuery(keys, 1, ",")
// 	conditionFieldsString := transformNamesToUpdateQuery(conditionColumns, len(keys)+1, " AND ")

// 	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, updateFieldsString, conditionFieldsString)

// 	values = append(values, conditionsValues...)

// 	_, err := db.Exec(query, values...)

// 	return err
// }

// func (acc *Account) IsExists(db *sql.DB) (bool, error) {
// 	//https://stackoverflow.com/questions/32554400/efficiently-determine-if-any-rows-satisfy-a-predicate-in-postgres?rq=3

// 	var exists bool

// 	keys, values := getNonEmptyField(acc)
// 	conditionString := transformNamesToUpdateQuery(keys, 1, " AND ")

// 	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s)", tableName, conditionString)

// 	fmt.Println(query)

// 	err := db.QueryRow(query, values...).Scan(&exists)

// 	if err != nil {
// 		return false, err
// 	}

// 	return exists, nil
// }
