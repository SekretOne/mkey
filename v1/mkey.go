package mkey

import (
	"encoding"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const TagDefaultName = "mkey"

const (
	_ = iota
	MarshalErr
	UnmarshalErr
)

type MultiKey[T any] struct {
	Val      T
	TagValue string
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

func exportedFields(t reflect.Type) []reflect.StructField {
	var exported []reflect.StructField
	for i := range t.NumField() {
		f := t.Field(i)

		if f.IsExported() {
			exported = append(exported, f)
		}
	}

	return exported
}

func lookupTagInstructions(f reflect.StructField, tag string) (string, string, bool) {
	var meta, terminator string

	if tv, ok := f.Tag.Lookup(tag); ok {
		splits := strings.SplitN(tv, " ", 2)

		if len(splits) > 2 {
			log.Panicf("bad tag format on type %T with tag value %q", f.Type, tv)
		}

		if len(splits) == 2 {
			terminator = splits[1]
		}

		meta = splits[0]

		return meta, terminator, true
	}

	return "", "", false
}

func MarshalFields(input any, tag string) (*types.AttributeValueMemberS, error) {
	var s string

	t := reflect.TypeOf(input)

	// list exported fields
	exported := exportedFields(t)

	for i, f := range exported {
		fv := reflect.ValueOf(input).FieldByIndex(f.Index)

		var meta, terminator string

		// if not last, set default terminator
		if i < len(exported)-1 {
			terminator = "#"
		}

		if pre, term, ok := lookupTagInstructions(f, tag); ok {
			meta = pre
			terminator = term
		}

		switch v := fv.Interface().(type) {
		case encoding.TextMarshaler:
			if bytes, err := v.MarshalText(); err != nil {
				return nil, err
			} else {
				s += meta + (string)(bytes) + terminator
			}
		default:
			s += meta + fmt.Sprint(v) + terminator
		}
	}

	return &types.AttributeValueMemberS{Value: s}, nil
}

func UnmarshalFields(output any, tag string, value types.AttributeValue) error {
	avs, ok := value.(*types.AttributeValueMemberS)
	if !ok {
		return errors.New("must use string member av") //todo error must use string member
	}

	val := reflect.ValueOf(output).Elem()
	t := val.Type()
	s := avs.Value
	//original := s // todo use in error messages

	exported := exportedFields(t)

	for i, field := range exported {
		f := t.FieldByIndex(field.Index)

		var meta, terminator string

		// if not last, set default terminator
		if i < len(exported)-1 {
			terminator = "#"
		}

		if pre, term, ok := lookupTagInstructions(f, tag); ok {
			meta = pre
			terminator = term
		}

		if s, ok = strings.CutPrefix(s, meta); !ok {
			return fmt.Errorf("expected prefix of %q not found in %q for into %s", s, meta, f.Type.Name())
		}

		ind := strings.Index(s, terminator)
		if ind == -1 {
			return errors.New("expected terminator of %q not found in %q for type %T")
		}

		v := s

		// if we have a separator, grab the content up to the terminator
		if len(terminator) > 0 {
			v = s[:ind]
			s = s[ind+len(terminator):]
		}

		switch fv := val.Field(i).Interface().(type) {
		case string:
			val.Field(i).SetString(v)
		case encoding.TextUnmarshaler:
			if err := fv.UnmarshalText(([]byte)(v)); err != nil {
				return fmt.Errorf("%w: cannot unmarshal text %v into %T", err, v, fv)
			}
		case int, int8, int16, int32, int64:
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			val.Field(i).SetInt(n)
		case uint, uint8, uint16, uint32, uint64:
			n, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			val.Field(i).SetUint(n)
		case float64:
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			val.Field(i).SetFloat(n)
		case float32:
			n, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return err
			}
			val.Field(i).SetFloat(n)
		default:
			return fmt.Errorf("no unmarshal option for field %v of type %T", f.Name, fv)
		}
	}

	return nil
}
