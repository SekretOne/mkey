package mkey

import (
	"encoding"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	TagDefaultName    = "mkey"
	defaultTerminator = "#"
)

type MultiKey[T any] struct {
	Val T
}

func (m MultiKey[T]) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return MarshalDynamoDBAttributeValue(m.Val)
}

func (m *MultiKey[T]) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	return UnmarshalDynamoDBAttributeValue(&m.Val, value)
}

// MarshalDynamoDBAttributeValue conforms to the aws sdk v2 attributevalue.Mashaler interface.
//
// Always marshals to types.AttributeValueMemberS. Fields of the struct are marshalled in their declared order,
// without meta prefix, delineated by the '#' character. No hanging terminator.
func MarshalDynamoDBAttributeValue(input any) (types.AttributeValue, error) {
	s, err := MarshalFields(input, TagDefaultName)
	if err != nil {
		return nil, err
	}

	return &types.AttributeValueMemberS{Value: s}, nil
}

// UnmarshalDynamoDBAttributeValue conforms to the aws sdk 2 attributevalue.Unmarshaler interface.
//
// Errors if the attribute value is not of string type, or cannot parse the subfields.
//
// Fields of the struct are unmarshalled in their declared order,  without meta prefix, delineated by the '#' character.
// No hanging terminator.
func UnmarshalDynamoDBAttributeValue(output any, value types.AttributeValue) error {
	avs, ok := value.(*types.AttributeValueMemberS)
	if !ok {
		return newUnmarshalError(
			fmt.Sprintf("expected value to be type *types.AttributeValueMemberS; got %T", value),
			nil,
		)
	}

	if err := UnmarshalFields(output, avs.Value, TagDefaultName); err != nil {
		return err
	}

	return nil
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

// parses the tag instructions under the given tag name.
//
// returns the meta, terminator and true when instructions found.
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

func MarshalFields(input any, tag string) (string, error) {
	var output string

	t := reflect.TypeOf(input)

	// list exported fields
	exported := exportedFields(t)

	for i, f := range exported {
		fv := reflect.ValueOf(input).FieldByIndex(f.Index)

		var meta, s, terminator string

		// if not last, set default terminator
		if i < len(exported)-1 {
			terminator = defaultTerminator
		}

		if pre, term, ok := lookupTagInstructions(f, tag); ok {
			meta = pre
			terminator = term
		}

		switch v := fv.Interface().(type) {
		case encoding.TextMarshaler:
			bytes, err := v.MarshalText()
			if err != nil {
				return "", newMarshalError(
					fmt.Sprintf("mkey: fail to marshal text on %T.%s: %s", input, f.Name, err.Error()),
					err,
				)
			}

			s = (string)(bytes)
		case string:
			s = v
		case uint, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			s = fmt.Sprint(v)
		default:
			return "", newMarshalError(
				fmt.Sprintf("mkey: no marshal option on field %T.%s", input, f.Name),
				nil)
		}

		output += meta + s + terminator
	}

	return output, nil
}

func UnmarshalFields(output any, s, tag string) error {
	val := reflect.ValueOf(output).Elem()
	t := val.Type()
	//original := s // todo use in error messages

	exported := exportedFields(t)

	for i, field := range exported {
		f := t.FieldByIndex(field.Index)

		var meta, terminator string

		// if not last, set default terminator
		if i < len(exported)-1 {
			terminator = defaultTerminator
		}

		if pre, term, ok := lookupTagInstructions(f, tag); ok {
			meta = pre
			terminator = term
		}

		if remaining, ok := strings.CutPrefix(s, meta); !ok {
			return fmt.Errorf("expected prefix of %q not found in %q for into %s", meta, s, f.Type.Name())
		} else {
			s = remaining
		}

		ind := strings.Index(s, terminator)
		if ind == -1 {
			return fmt.Errorf("expected terminator %q not found in string %q for type %T", terminator, s, f)
		}

		v := s

		// if we have a separator, grab the content up to the terminator
		if len(terminator) > 0 {
			v = s[:ind]
			s = s[ind+len(terminator):]
		}

		fv := val.FieldByIndex(f.Index)

		switch tv := fv.Interface().(type) {
		case encoding.TextUnmarshaler:
			if err := tv.UnmarshalText(([]byte)(v)); err != nil {
				return fmt.Errorf("%w: cannot unmarshal text %v into %T", err, v, fv)
			}
		case string:
			fv.SetString(v)
		case int, int8, int16, int32, int64:
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			fv.SetInt(n)
		case uint, uint8, uint16, uint32, uint64:
			n, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			fv.SetUint(n)
		case float64:
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			fv.SetFloat(n)
		case float32:
			n, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return err
			}
			fv.SetFloat(n)
		default:
			return fmt.Errorf("no unmarshal option for field %v of type %T", f.Name, fv)
		}
	}

	return nil
}
