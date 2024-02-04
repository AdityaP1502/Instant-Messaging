package model

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
	Delete(db *sql.DB) error
	IsExists(db *sql.DB) (bool, error)

	// Query(db *sql.DB) error
	// JoinQuery(db *sql.DB) error
}

func getNonEmptyField(v interface{}) ([]string, []any) {
	s := reflect.ValueOf(v).Elem()
	typeOfS := s.Type()

	names := make([]string, 0, 8)
	values := make([]any, 0, 8)

	for i := 0; i < typeOfS.NumField(); i++ {
		field := typeOfS.Field(i)
		columnTag := field.Tag.Get("column")

		// Gatekeep conditional
		if columnTag == "-" || columnTag == "" {
			continue
		}

		k := strings.SplitAfter(columnTag, ",")[0]
		v := s.Field(i).Interface()

		// Check if a field is empty/has value of "zero"
		if v != reflect.Zero(s.Field(i).Type()).Interface() {
			names = append(names, k)
			values = append(values, v)
		}
	}

	return names, values
}

func transformNamesToUpdateQuery(names []string, start int, sep string) string {
	fmt.Println(len(names))
	q := ""
	c := start

	for _, k := range names {
		q += fmt.Sprintf("%s=$%d%s", k, c, sep)
		c++
	}

	return q[:len(q)-len(sep)]
}
