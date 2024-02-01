package database

import (
	"database/sql"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Model interface {
	FromJSON(r io.Reader, checkRequired bool) error
	ToJSON(checkRequired bool) ([]byte, error)

	Insert(db *sql.DB) error
	Update(db *sql.DB, condition map[string]string) error
	IsExists(db *sql.DB) (bool, error)

	// Query(db *sql.DB) error
	// JoinQuery(db *sql.DB) error
}

func getNonEmptyField(v interface{}) ([]string, []any) {
	s := reflect.ValueOf(v)
	typeOfS := s.Type()

	names := make([]string, 8)
	values := make([]any, 8)

	for i := 0; i < typeOfS.NumField(); i++ {
		field := typeOfS.Field(i)
		jsonTag := field.Tag.Get("json")

		// Gatekeep conditional
		if jsonTag == "-" || jsonTag == "" {
			continue
		}

		k := strings.SplitAfter(jsonTag, ",")[0]
		v := s.Field(i).Interface()

		// Check if a field is empty/has value of "zero"
		if v != reflect.Zero(s.Field(i).Type()).Interface() {
			names = append(names, k)
			values = append(values, v)
		}
	}

	return names, values
}

func transformNamesToUpdateQuery(names []string, start int) string {
	q := ""
	c := start

	for k := range names {
		q += fmt.Sprintf("%s=$%d", k, c)
		c++
	}

	return q
}
