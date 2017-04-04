// Package storage can store different types by using a simple
// api
package storage

type dataType int

const (
	STRING dataType = iota
	STRINGSLICE
	STRINGPOINTER
	INT
	INTPOINTER
	FLOAT64
)

type Data struct {
	values       map[string]dataType
	strings      map[string]string
	pstrings     map[string]*string
	stringSlices map[string][]string
	ints         map[string]int
	pints        map[string]*int
	floats64     map[string]float64
}

func New() *Data {
	return &Data{
		values:       make(map[string]dataType),
		strings:      make(map[string]string),
		pstrings:     make(map[string]*string),
		stringSlices: make(map[string][]string),
		ints:         make(map[string]int),
		pints:        make(map[string]*int),
		floats64:     make(map[string]float64),
	}
}

func (d *Data) Set(key string, i interface{}) error {
	switch i := i.(type) {
	case string:
		d.values[key] = STRING
		d.strings[key] = i
	case []string:
		d.values[key] = STRINGSLICE
		d.stringSlices[key] = i
	case *string:
		d.values[key] = STRINGPOINTER
		d.pstrings[key] = i
	case int:
		d.values[key] = INT
		d.ints[key] = i
	case *int:
		d.values[key] = INTPOINTER
		d.pints[key] = i
	case float64:
		d.values[key] = FLOAT64
		d.floats64[key] = i
	}
	return nil
}

func (d *Data) Get(key string) (interface{}, error) {
	dType := d.values[key]
	switch dType {
	case STRING:
		return d.strings[key], nil
	case STRINGSLICE:
		return d.stringSlices[key], nil
	case STRINGPOINTER:
		return d.pstrings[key], nil
	case INT:
		return d.ints[key], nil
	case INTPOINTER:
		return d.pints[key], nil
	case FLOAT64:
		return d.floats64[key], nil
	}
	return nil, nil
}
