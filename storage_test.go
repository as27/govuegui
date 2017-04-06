package govuegui

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var str = "myString"
var myint = 123456
var testCases = []struct {
	key   string
	dType dataType
	value interface{}
}{
	{
		"int1",
		INT,
		123,
	},
	{
		"int2",
		INT,
		456,
	},
	{
		"string",
		STRING,
		"myString",
	},
	{
		"float",
		FLOAT64,
		float64(1.1234),
	},
	{
		"sliceOfStrings",
		STRINGSLICE,
		[]string{"abc", "def", "hij"},
	},
	{
		"string pointer",
		STRINGPOINTER,
		&str,
	},
	{
		"pointer to int",
		INTPOINTER,
		&myint,
	},
	{
		"element options",
		OPTION,
		[]Option{
			Option{
				"class",
				[]string{"active", "blue"},
			},
			Option{
				"title",
				[]string{"myTitle"},
			},
		},
	},
}

func ExampleData_Set() {
	store := New()
	store.Set("myString", "this is my string")
	store.Set("myInt", 1234)
	i := store.Get("myInt").(int)
	fmt.Println(i * 2)
	// Output: 2468
}
func TestSet(t *testing.T) {
	data := New()
	for _, tc := range testCases {
		data.Set(tc.key, tc.value)
		if data.Values[tc.key] != tc.dType {
			t.Errorf("Set() wrong type!\nInput: %s\nGot: %v Exp: %v",
				tc.key,
				data.Values[tc.key],
				tc.dType,
			)
		}
	}

}

func TestGet(t *testing.T) {
	data := New()
	for _, tc := range testCases {
		data.Set(tc.key, tc.value)
		val := data.Get(tc.key)
		if !reflect.DeepEqual(val, tc.value) {
			t.Errorf("Get() wrong value\nInput: %s\nGot: %v Exp: %v",
				tc.key,
				val,
				tc.value,
			)
		}
	}
}

func TestErrors(t *testing.T) {
	data := New()
	type notSupported int
	a := notSupported(123)
	err := data.Set("notSupported", a)
	if err != ErrTypeNotSupported {
		t.Error("Should return that type is not supported")
	}
	err = data.Set("myString", "a string")
	if err != nil {
		t.Error("Should return no error")
	}
	_, err = data.GetWithErrors("NotExist")
	if err != ErrKeyNotFound {
		t.Error("Should return error that key is not found")
	}
}

func TestMarshal(t *testing.T) {
	data := New()
	for _, tc := range testCases {
		data.Set(tc.key, tc.value)
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		//t.Error("Can not Marshal the data")
		t.Error(err)
	}
	d := New()
	err = json.Unmarshal(b, d)
	if err != nil {
		t.Error(err)
	}
}

func TestRemove(t *testing.T) {
	data := New()
	for _, tc := range testCases {
		data.Set(tc.key, tc.value)
		val, err := data.GetWithErrors(tc.key)
		if err != nil || !reflect.DeepEqual(val, tc.value) {
			t.Errorf("Get should return %s", tc.value)
		}
		ok := data.Remove(tc.key)
		if !ok {
			t.Errorf("Should return true, because key: %s exists.", tc.key)
		}
		_, err = data.GetWithErrors(tc.key)
		if err != ErrKeyNotFound {
			t.Errorf("Should not find the key %s", tc.key)
		}

	}

}
