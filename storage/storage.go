// Package storage can store different types by using a simple
// api. Just using Get() and Set() and Remove()
//    store := storage.New()
//    store.Set("MyString","This is my string")
//    store.Set("MyInt", 1234)
//    // Set type when use Get()
//    i := store.Get("MyInt").(int)
package storage

import (
	"encoding/json"
	"errors"
)

type dataType int

const (
	STRING dataType = iota
	STRINGSLICE
	STRINGPOINTER
	INT
	INTPOINTER
	FLOAT64
)

// ErrTypeNotSupported is returned by Set() when the given type could
// not be stored.
var ErrTypeNotSupported = errors.New("Type is not supported!")

// ErrKeyNotFound is returned by GetWithErrors() when the key is not found
var ErrKeyNotFound = errors.New("The given key was not found inside storage!")

// Data is the type were everything is stored and which can be used for
// Marshaling to json.
type Data struct {
	Values       map[string]dataType `json:"values"`
	Strings      map[string]string   `json:"strings"`
	PStrings     map[string]*string  `json:"pstrings"`
	StringSlices map[string][]string `json:"stringSlices"`
	Ints         map[string]int      `json:"ints"`
	PInts        map[string]*int     `json:"pints"`
	Floats64     map[string]float64  `json:"floats64"`
}

// New returns a pointer to a new empty storage
func New() *Data {
	return &Data{
		Values:       make(map[string]dataType),
		Strings:      make(map[string]string),
		PStrings:     make(map[string]*string),
		StringSlices: make(map[string][]string),
		Ints:         make(map[string]int),
		PInts:        make(map[string]*int),
		Floats64:     make(map[string]float64),
	}
}

// Set is used to set the data of a type.
func (d *Data) Set(key string, i interface{}) error {
	switch i := i.(type) {
	case string:
		d.Values[key] = STRING
		d.Strings[key] = i
	case []string:
		d.Values[key] = STRINGSLICE
		d.StringSlices[key] = i
	case *string:
		d.Values[key] = STRINGPOINTER
		d.PStrings[key] = i
	case int:
		d.Values[key] = INT
		d.Ints[key] = i
	case *int:
		d.Values[key] = INTPOINTER
		d.PInts[key] = i
	case float64:
		d.Values[key] = FLOAT64
		d.Floats64[key] = i
	default:
		return ErrTypeNotSupported
	}
	return nil

}

// Get enables a simple api, which just returns the value
// if no value is found nil is returned. If it is clear which
// type the value is it is possible to set the type when calling
// that function.
//    i := store.Get("myInt").(int)
// This usage is dangerous, because it can cause a runtime error,
// when the found value is another type then expected.
func (d *Data) Get(key string) interface{} {
	i, _ := d.GetWithErrors(key)
	return i
}

// GetWithErrors get a value and returns an error. Because of the
// api there is a Get() and a GetWithErrors() function.
func (d *Data) GetWithErrors(key string) (interface{}, error) {
	dType, ok := d.Values[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	switch dType {
	case STRING:
		return d.Strings[key], nil
	case STRINGSLICE:
		return d.StringSlices[key], nil
	case STRINGPOINTER:
		return d.PStrings[key], nil
	case INT:
		return d.Ints[key], nil
	case INTPOINTER:
		return d.PInts[key], nil
	case FLOAT64:
		return d.Floats64[key], nil
	}
	return nil, nil
}

// Remove just deletes the key from the storage. True is returned
// when something could be removed. If the key not exists inside
// the storage nothing could be deleted, so false is returned.
func (d *Data) Remove(key string) bool {
	dType, ok := d.Values[key]
	if !ok {
		return false
	}
	delete(d.Values, key)
	switch dType {
	case STRING:
		delete(d.Strings, key)
	case STRINGSLICE:
		delete(d.StringSlices, key)
	case STRINGPOINTER:
		delete(d.PStrings, key)
	case INT:
		delete(d.Ints, key)
	case INTPOINTER:
		delete(d.PInts, key)
	case FLOAT64:
		delete(d.Floats64, key)
	default:
		return false
	}
	return true
}

// Marshal the storage into json format
func (d *Data) Marshal() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

// Unmarshal a slice of bytes into the data
func Unmarshal(b []byte) (*Data, error) {
	data := New()
	err := json.Unmarshal(b, data)
	return data, err
}
