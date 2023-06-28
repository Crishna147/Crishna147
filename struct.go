package main

import (
	"fmt"
	"reflect"
)

type employee struct {
	Ename string `json:"ename" db:"full_name"`
	EID   int    `json:"eid"  db:"eid"`
	Esal  int    `json:"esal"  db:"esal"`
}

func main() {

	e := employee{
		Ename: "Shinchan",
		EID:   147,
		Esal:  147125,
	}
	t := reflect.TypeOf(e)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagEName := field.Tag.Get("json")
		fmt.Printf("Field: %s, JSON Tag: %s\n", field.Name, tagEName)
	}

}
------------------------------------------------------------------------
PS D:\go\sample1> go run struct1.go     
Field: Ename, JSON Tag: ename
Field: EID, JSON Tag: eid
Field: Esal, JSON Tag: esal








