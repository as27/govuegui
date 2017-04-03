package storage

import (
	"reflect"
	"testing"
)

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
}

func TestSet(t *testing.T) {
	data := New()
	for _, tc := range testCases {
		data.Set(tc.key, tc.value)
		if data.values[tc.key] != tc.dType {
			t.Errorf("Set() wrong type!\nInput: %s\nGot: %v Exp: %v",
				tc.key,
				data.values[tc.key],
				tc.dType,
			)
		}
	}

}

func TestGet(t *testing.T) {
	data := New()
	for _, tc := range testCases {
		data.Set(tc.key, tc.value)
		val, _ := data.Get(tc.key)
		if !reflect.DeepEqual(val, tc.value) {
			t.Errorf("Get() wrong value\nInput: %s\nGot: %v Exp: %v",
				tc.key,
				val,
				tc.value,
			)
		}
	}
}
