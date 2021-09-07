package tests

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInsertData(t *testing.T) {
	//InsertData()
	//InsertMetadata()
	//InsertResource()
	//InsertTask()
	//RevokeIdentity()
	GetData()

	type testSlice []int
	as := make(testSlice, 1)
	as = append(as, 2)
	fmt.Println("test slice", len(as))
	fmt.Println("testSlice,type is:", reflect.TypeOf(testSlice{}))
}
