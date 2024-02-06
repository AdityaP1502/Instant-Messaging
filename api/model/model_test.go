package model

import (
	"testing"
)

func TestFieldToQueryOnce(t *testing.T) {
	keys := []string{"username", "name", "password", "is_active"}

	query := transformNamesToUpdateQuery(keys, 1, ",")
	expected := "username=$1,name=$2,password=$3,is_active=$4"

	if query != expected {
		t.Errorf("Expected %s received %s", expected, query)
	}

	t.Logf("%s is correct", query)
}

func TestFieldToQueryTwice(t *testing.T) {
	keys := []string{"username", "name", "password", "is_active"}
	condition := []string{"foo", "bar"}

	_ = transformNamesToUpdateQuery(keys, 1, ",")
	query2 := transformNamesToUpdateQuery(condition, len(keys)+1, " AND ")

	if query2 != "foo=$5 AND bar=$6" {
		t.Errorf("Expected foo=$5 AND bar=$6 received %s", query2)
	}

	t.Logf("%s is correct", query2)
}

func TestGetNonEmptyField(t *testing.T) {
	data := &Account{
		Username: "asep",
		Password: "1234",
	}

	keys, _, _ := getNonEmptyField(data)

	t.Log(keys)

	query := transformNamesToUpdateQuery(keys, 1, ",")

	expected := "username=$1,password=$2"

	if query != expected {
		t.Errorf("Expected %s received %s", expected, query)
	}

	t.Logf("%s is correct", query)

}
