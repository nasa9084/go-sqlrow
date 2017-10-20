package sqlrow_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	sqlmock "github.com/nasa9084/go-sql-mock"
	sqlrow "github.com/nasa9084/go-sqlrow"
)

type testType string

type testStruct struct {
	Name string
	Age  int
}

func getRow(cols []string, vals []driver.Value) *sql.Row {
	db, err := sql.Open("sqlmock", "")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlmock.ExpectedRows(
		sqlmock.Columns(cols),
		sqlmock.ValuesList(vals),
	)
	return db.QueryRow("")
}

func TestBinderStruct(t *testing.T) {
	var ts testStruct
	row := getRow(
		[]string{"name", "age"},
		[]driver.Value{"something", 20},
	)
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

func ExampleBinder() {
	db, _ := sql.Open("mysql", "...")
	row := db.QueryRow("SELECT * FROM ...")
	var ts testStruct
	if err := sqlrow.NewBinder(row).Bind(&ts); err != nil {
		panic(err)
	}
}
