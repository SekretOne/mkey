package mkey

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"reflect"
	"testing"
)

// AssertSymmetricalInDynamodb verifies an input with marshal to dynamodb and restore back to the original value.
//
// First marshals to a dynamodb attribute value, then creates a new zero value instance of the same type to unmarshal
// the result, and finally compares with the original input.
func AssertSymmetricalInDynamodb[T any](t *testing.T, input T) {
	t.Helper()

	m, err := attributevalue.Marshal(input)
	if err != nil {
		t.Fatalf(`%T failed marshal to atttibute value 

input: 
	%v

error: 
	%v`, input, yellow(fmt.Sprintf("%+v", input)), red(err.Error()))
	}

	var output T

	err = attributevalue.Unmarshal(m, &output)
	if err != nil {
		t.Fatalf("failed unmarshal from atttibute value using: %+v got err: %v", m, err)
	}

	if !reflect.DeepEqual(input, output) {
		t.Errorf("%T is not symmetrical got: %+v want: %+v", output, output, input)
	}
}

func red(s string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", s)
}
func yellow(s string) string {
	return fmt.Sprintf("\033[33m%s\033[0m", s)
}
