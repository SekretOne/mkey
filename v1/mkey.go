package mkey

import (
	"encoding"
	"errors"
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

func MarshalFields(input any, tag string) (*types.AttributeValueMemberS, error) {
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
		case encoding.TextMarshaler:
			if bytes, err := v.MarshalText(); err != nil {
				return nil, err
			} else {
				part = fmt.Sprintf("%s=%s", fn, (string)(bytes))
			}
		default:
			part = fmt.Sprintf("%s=%s", fn, v)
		}

		parts = append(parts, part)
	}

	memberString := strings.Join(parts, "|")
	return &types.AttributeValueMemberS{Value: memberString}, nil
}

func UnmarshalFields(output any, tag string, value types.AttributeValue) error {
	avs, ok := value.(*types.AttributeValueMemberS)
	if !ok {
		return errors.New("must use string member av") //todo error must use string member
	}

	val := reflect.ValueOf(output).Elem()
	t := val.Type()
	s := avs.Value
	// original := s // todo use in error messages

	for i := range t.NumField() {
		f := t.Field(i)

		if !f.IsExported() {
			continue
		}

		fn := f.Name + "="
		separator := "|"
		if name, ok := f.Tag.Lookup(tag); ok {
			fn = name
		}

		if s, ok = strings.CutPrefix(s, fn); !ok {
			return errors.New("cut prefix") // todo
		}

		ind := strings.Index(s, separator)
		if ind == -1 {
			return errors.New("no separator") //todo not found separator
		}

		v := s[:ind]
		s = s[ind+len(separator):]

		switch fv := val.Field(i).Interface().(type) {
		case string:
			val.Field(i).SetString(v)
		case encoding.TextUnmarshaler:
			if err := fv.UnmarshalText(([]byte)(v)); err != nil {
				return err //todo better error
			}
		}
	}

	return nil
}
