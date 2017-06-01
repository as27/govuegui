/*
Package storage defines a object, which can be exported and imported
as a JSON object. The data is stored with the type information, so inside
Go the type logic can be used. */
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type DataType string

// Option holds the one option of a element
type Option struct {
	Option string
	Values []string
}

const (
	STRING        DataType = "STRING"
	STRINGSLICE            = "STRINGARRAY"
	STRINGPOINTER          = "STRINGPOINTER"
	STRINGTABLE            = "STRINGTABLE"
	INT                    = "INT"
	INTPOINTER             = "INTPOINTER"
	FLOAT64                = "FLOAT64"
	OPTION                 = "OPTION"
	FUNCPOINTER            = "FUNCPOINTER"
)

// isExportedType defines, which types are exportet to JS
func isExportedType(d DataType) bool {
	if d == FUNCPOINTER {
		return false
	}
	return true
}

// ErrTypeNotSupported is returned by Set() when the given type could
// not be stored.
var ErrTypeNotSupported = errors.New("Type is not supported!")

// ErrKeyNotFound is returned by GetWithErrors() when the key is not found
var ErrKeyNotFound = errors.New("The given key was not found inside storage!")

// Data is the type were everything is stored and which can be used for
// Marshaling to json. Every entry uses a unique key, which is a string.
type Data struct {
	// Values is a map which know the internal DataType every key value
	// pair is exported here.
	Values map[string]DataType `json:"values"`
	// Data conatains all the data, which can be exported
	Data map[string]interface{} `json:"data"`
	// uneportedData contains all the types, which can not be marshaled
	// into json format.
	unexportedData map[string]interface{}
	// cache is used, when modifiing a value of pointer
	cache map[string]interface{}
}

// NewStorage returns a pointer to a new empty storage
func New() *Data {
	return &Data{
		Values:         make(map[string]DataType),
		Data:           make(map[string]interface{}),
		unexportedData: make(map[string]interface{}),
		cache:          make(map[string]interface{}),
	}
}

// Set is used to set the data of a type. Everytime a new value is set
// the type of the newest type is used.
func (d *Data) Set(key string, i interface{}) error {
	switch i.(type) {
	case string:
		d.Values[key] = STRING
	case []string:
		d.Values[key] = STRINGSLICE
	case [][]string:
		d.Values[key] = STRINGTABLE
	case *string:
		d.Values[key] = STRINGPOINTER
	case int:
		d.Values[key] = INT
	case *int:
		d.Values[key] = INTPOINTER
	case float64:
		d.Values[key] = FLOAT64
	case []*Option:
		d.Values[key] = OPTION
	case *func():
		d.Values[key] = FUNCPOINTER
	default:
		return ErrTypeNotSupported
	}
	if d.isExportedData(key) {
		d.Data[key] = i
	} else {
		d.unexportedData[key] = i
	}
	d.cache[key] = i

	return nil
}

// Get enables a simple api, which just returns the value
// if no value is found nil is returned. If it is clear which
// type the value is it is possible to set the type when calling
// that function.
//    i := store.Get("myInt").(int)
// This usage is dangerous, because it can cause a runtime error,
// when the found value is another type then expected.
// But it enables a very simple API, so if you use it be sure that
// the type you are getting is really fix
// To be sure at this point the function GetWithErrors() should be used.
func (d *Data) Get(key string) interface{} {
	i, _ := d.GetWithErrors(key)
	return i
}

// GetWithErrors get a value and returns an error. Because of the
// api there is a Get() and a GetWithErrors() function.
func (d *Data) GetWithErrors(key string) (interface{}, error) {
	_, ok := d.Values[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	if !d.isExportedData(key) {
		return d.unexportedData[key], nil
	}
	return d.Data[key], nil
}

// GetType returns the DataType of a key. If there is no value availiable
// with the key an error is returned.
func (d *Data) GetType(key string) (DataType, error) {
	var err error
	dt, ok := d.Values[key]
	if !ok {
		err = ErrKeyNotFound
	}
	return dt, err
}

// GetKeys returns all the keys a slice of strings
func (d *Data) GetKeys() []string {
	var keys []string
	for k := range d.Values {
		keys = append(keys, k)
	}
	return keys
}

// Remove just deletes the key from the storage. True is returned
// when something could be removed. If the key not exists inside
// the storage nothing could be deleted, so false is returned.
func (d *Data) Remove(key string) bool {
	_, ok := d.Values[key]
	if !ok {
		return false
	}
	delete(d.Values, key)
	delete(d.Data, key)
	return true
}

// Unmarshal a slice of bytes into the data
func (d *Data) Unmarshal(b []byte) error {
	return nil
}

// SetData takes a storage and sets all the values. This also works for
// all the pointers. That function is used, when the storage is
// unmarshaled input.
// Just exported values are set new.
func (d *Data) SetData(data *Data) error {
	var err error
	for k, dType := range d.Values {
		switch dType {
		default:
			d.Data[k] = data.Data[k]
		case STRINGPOINTER:
			sp := d.cache[k].(*string)
			*sp = data.Data[k].(string)
			d.Data[k] = sp
		case INT:
			var v float64
			v, err = interfaceToFloat(data.Data[k])
			d.Data[k] = int(v)
		case INTPOINTER:
			var v float64
			v, err = interfaceToFloat(data.Data[k])
			ip := d.cache[k].(*int)
			*ip = int(v)
			d.Data[k] = ip
		case FUNCPOINTER:
			fp := d.cache[k].(*func())
			*fp = data.unexportedData[k].(func())
			d.Data[k] = fp
		case FLOAT64:
			var v float64
			v, err = interfaceToFloat(data.Data[k])
			d.Data[k] = v
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func interfaceToFloat(i interface{}) (float64, error) {
	switch v := i.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("interfaceToFloat: %T not expected Type", i)
	}
}

func (d *Data) isExportedData(key string) bool {
	return isExportedType(d.Values[key])
}

// Marshal the storage into json format
func (d *Data) Marshal() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}
