package sqlrow

import (
	"database/sql"
	"reflect"
)

type bindErr string

func (be bindErr) Error() string { return string(be) }

const (
	// ErrMustPointer is returned if given instance for Bind is not pointer
	ErrMustPointer bindErr = `binding destination instance must be a pointer`
)

// Bind function for convenience
func Bind(row *sql.Row, i interface{}) error {
	return NewBinder(row).Bind(i)
}

// Binder binds sql row to golang type
type Binder struct {
	row *sql.Row
	err error
}

// NewBinder returns a new Binder
func NewBinder(row *sql.Row) *Binder {
	return &Binder{row: row}
}

// Bind binds row to given instance
func (binder *Binder) Bind(i interface{}) error {
	if binder.err != nil {
		return binder.err
	}
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		return ErrMustPointer
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return binder.row.Scan(i)
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	fi := []interface{}{}
	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).CanSet() {
			continue
		}
		fi = append(fi, v.Field(i).Addr().Interface())
	}

	return binder.row.Scan(fi...)
}
