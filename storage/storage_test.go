package storage

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var str = "myString"
var myint = 123456
var myfunc = func() {}
var testCases = []struct {
	key   string
	dType DataType
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
		"stringTable",
		STRINGTABLE,
		[][]string{
			{"abc", "def", "hij"},
			{"abc", "def", "hij"},
		},
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
		"pointer to func",
		FUNCPOINTER,
		&myfunc,
	},
	{
		"element options",
		OPTION,
		[]*Option{
			&Option{
				"class",
				[]string{"active", "blue"},
			},
			&Option{
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

func TestGetType(t *testing.T) {
	data := New()
	for _, tc := range testCases {
		data.Set(tc.key, tc.value)
		dt, _ := data.GetType(tc.key)
		if dt != tc.dType {
			t.Errorf("GetType() returns wrong type!\nInput: %s\nGot: %v Exp: %v",
				tc.key,
				data.Values[tc.key],
				tc.dType,
			)
		}
	}
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
	d := data
	err = d.Unmarshal(b)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(data.Data, d.Data) {
		t.Errorf("Unmarshal gets another result\nExp: %#v\nGot: %#v", data, d)
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

func TestInterfaceToFloat(t *testing.T) {
	tests := []struct {
		invalue  interface{}
		expvalue float64
		expError bool
	}{
		{
			123,
			float64(123),
			false,
		},
		{
			"123",
			float64(123),
			false,
		},
		{
			float64(123),
			float64(123),
			false,
		},
		{
			"string",
			float64(0),
			true,
		},
	}
	for _, test := range tests {
		got, err := interfaceToFloat(test.invalue)
		if test.expError && err == nil {
			t.Error("Testcase: ", test.invalue)
			t.Error("Expected an error!")
		}
		if test.expvalue != got {
			t.Errorf("Exp: %v\nGot: %v", test.expvalue, got)
		}
	}
}
