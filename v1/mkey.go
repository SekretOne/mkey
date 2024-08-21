package mkey

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"strings"
)

const tagDefault = "multi-key"

const (
	_ = iota
	MarshalErr
	UnmarshalErr
)

type Substringer interface {
	MarshallSubstring() (string, error)
	UnmarshallSubstring(input string) error
}

type MultiKey[T any] struct {
	Val T
}

func (m MultiKey[T]) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MultiKey[T]) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	panic("implement me")

	//sv, ok := value.(*types.AttributeValueMemberS)
	//if !ok {
	//	return fmt.Errorf("expected value to be type *types.AttributeValueMemberS; got %T", value)
	//}
	//
	//s := sv.Value
}

func MarshalMultiFieldKeyWithTag(input any, tag string) (*types.AttributeValueMemberS, error) {
	var parts []string

	t := reflect.TypeOf(input)

	for _, f := range reflect.VisibleFields(t) {
		if !f.IsExported() {
			continue
		}

		fn := f.Name
		fv := reflect.ValueOf(input).FieldByIndex(f.Index)

		if name, ok := f.Tag.Lookup(tag); ok {
			fn = name
		}

		var part string

		switch v := fv.Interface().(type) {
		case string:
			part = fmt.Sprintf("%s=%s", fn, v)
		default:
			panic(fmt.Sprintf("expected string or struct tag; got %T", fv.Interface()))
		}

		parts = append(parts, part)
	}

	memberString := strings.Join(parts, "|")
	return &types.AttributeValueMemberS{Value: memberString}, nil
}
