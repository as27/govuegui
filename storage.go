package govuegui

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type dataType string

const (
	STRING        dataType = "STRING"
	STRINGSLICE            = "STRINGARRAY"
	STRINGPOINTER          = "STRINGPOINTER"
	INT                    = "INT"
	INTPOINTER             = "INTPOINTER"
	FLOAT64                = "FLOAT64"
	OPTION                 = "OPTION"
	FUNCPOINTER            = "FUNCPOINTER"
)

func isExportedType(d dataType) bool {
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
	// Values is a map which know the internal datatype every key value
	// pair is exported here.
	Values map[string]dataType `json:"values"`
	// Data conatains all the data, which can be exported
	Data map[string]interface{} `json:"data"`
	// uneportedData contains all the types, which can not be marshaled
	// into json format.
	unexportedData map[string]interface{}
	// cache is used, when modifiing a value of pointer
	cache map[string]interface{}
}

// NewStorage returns a pointer to a new empty storage
func NewStorage() *Data {
	return &Data{
		Values:         make(map[string]dataType),
		Data:           make(map[string]interface{}),
		unexportedData: make(map[string]interface{}),
		cache:          make(map[string]interface{}),
	}
}

// Set is used to set the data of a type.
func (d *Data) Set(key string, i interface{}) error {
	switch i.(type) {
	case string:
		d.Values[key] = STRING
	case []string:
		d.Values[key] = STRINGSLICE
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

// SetData a slice of bytes into the data
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
