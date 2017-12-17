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

type testStructWithUnexported struct {
	Name string
	age  int
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

func TestBinderWithUnexported(t *testing.T) {
	var ts testStructWithUnexported
	row := getRow(
		[]string{"name"},
		[]driver.Value{"something"},
	)
	if err := sqlrow.NewBinder(row).Bind(&ts); err != nil {
		t.Errorf("err: %s", err)
		return
	}
	if ts.Name != "something" {
		t.Errorf(`"%s" != "something"`, ts.Name)
		return
	}
}

func TestBinderStruct(t *testing.T) {
	var ts testStruct
	row := getRow(
		[]string{"name", "age"},
		[]driver.Value{"something", 20},
	)
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

func TestBinderType(t *testing.T) {
	var tt testType
	row := getRow(
		[]string{"value"},
		[]driver.Value{"something"},
	)
	if err := sqlrow.NewBinder(row).Bind(&tt); err != nil {
		t.Errorf("err: %s", err)
		return
	}
	if string(tt) != "something" {
		t.Errorf(`"%s" != "something"`, string(tt))
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

type benchmarkStruct struct {
	A string
	B int
	C float64
	D bool
}

func BenchmarkBinder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		row := getRow(
			[]string{"a", "b", "c", "d"},
			[]driver.Value{"something", 20, 100.1, true},
		)
		var bs benchmarkStruct
		if err := sqlrow.NewBinder(row).Bind(&bs); err != nil {
			panic(err)
		}
	}
}

func BenchmarkRawScan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		row := getRow(
			[]string{"a", "b", "c", "d"},
			[]driver.Value{"something", 20, 100.1, true},
		)
		var bs benchmarkStruct
		if err := row.Scan(&bs.A, &bs.B, &bs.C, &bs.D); err != nil {

			panic(err)
		}
	}
}
