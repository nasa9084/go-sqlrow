package sqlrow_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	sqlmock "github.com/nasa9084/go-sql-mock"
	sqlrow "github.com/nasa9084/go-sqlrow"
)

type testStruct struct {
	Name string
	Age  int
}

type NilStruct struct{}

func getRow() *sql.Row {
	db, err := sql.Open("sqlmock", "")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlmock.ExpectedRows(
		sqlmock.Columns([]string{"name", "age"}),
		sqlmock.ValuesList([]driver.Value{"something", 20}),
	)
	return db.QueryRow("")
}

func TestBinder(t *testing.T) {
	var ts testStruct
	row := getRow()
	t.Log(row)
	if err := sqlrow.NewBinder(row).Bind(&ts); err != nil {
		t.Errorf("err: %s", err)
		return
	}
	if ts.Name != "something" {
		t.Errorf(`"%s" != "something"`, ts.Name)
		return
	}
	if ts.Age != 20 {
		t.Errorf("%d != 20", ts.Age)
		return
	}
}
