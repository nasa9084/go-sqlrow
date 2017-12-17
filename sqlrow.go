package sqlrow

import (
	"reflect"
)

type bindErr string

func (be bindErr) Error() string { return string(be) }

const (
	// ErrMustPointer is returned if given instance for Bind is not pointer
	ErrMustPointer bindErr = `binding destination instance must be a pointer`
)

// Scanner is wrapper interface of type like *sql.Row.
type Scanner interface {
	Scan(...interface{}) error
}

// Bind function for convenience
func Bind(sc Scanner, i interface{}) error {
	return NewBinder(sc).Bind(i)
}

// Binder binds sql row to golang type
type Binder struct {
	sc  Scanner
	err error
}

// NewBinder returns a new Binder
func NewBinder(sc Scanner) *Binder {
	return &Binder{sc: sc}
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
		return binder.sc.Scan(i)
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

	return binder.sc.Scan(fi...)
}
