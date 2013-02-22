package main

import(
	"reflect"
)

func isMap(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Map
}
