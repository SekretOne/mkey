package mkey

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"testing"
)

func TestMarshalMultiFieldKeyWithTag(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		withTag string // can be blank, use default
		want    *types.AttributeValueMemberS
		wantErr bool
	}{
		{
			name: "The simplest example",
			input: struct {
				A string
				B string
			}{
				A: "first",
				B: "second",
			},
			want: &types.AttributeValueMemberS{
				Value: "first#second",
			},
			wantErr: false,
		},
		{
			name: "A series of strings",
			input: struct {
				First  string `mkey:"First= |"`
				Second string `mkey:"Second= |"`
				Third  string `mkey:"Third="`
			}{
				First:  "1st",
				Second: "2nd",
				Third:  "3rd",
			},
			withTag: "mkey",
			want: &types.AttributeValueMemberS{
				Value: "First=1st|Second=2nd|Third=3rd",
			},
		},
		{
			name: "overriding the tag",
			input: struct {
				A string `mkey:"a-term: ::"`
				B string `mkey:"b-term: ::"`
				C string `mkey:"c-term:"`
			}{
				A: "alpha",
				B: "beta",
				C: "charlie",
			},
			want: &types.AttributeValueMemberS{
				Value: "a-term:alpha::b-term:beta::c-term:charlie",
			},
			withTag: TagDefaultName,
		},
		{
			name: "integers",
			input: struct {
				A int
				B int8
				C uint64
			}{
				A: -20,
				B: 0,
				C: 30,
			},
			want: &types.AttributeValueMemberS{
				Value: "-20#0#30",
			},
		},
		{
			name: "Non exported fields are skipped",
			input: struct {
				privateA string
				B        string
				C        string
				privateD string
			}{
				B: "value-b",
				C: "value-c",
			},
			want: &types.AttributeValueMemberS{
				Value: "value-b#value-c",
			},
		},
		{
			name: "use an explicit terminator on the last sub field",
			input: struct {
				One   string
				Two   string
				Three string `mkey:" #"`
			}{
				One:   "foo",
				Two:   "biz",
				Three: "baz",
			},
			want: &types.AttributeValueMemberS{
				Value: "foo#biz#baz#",
			},
			withTag: TagDefaultName,
		},
		{
			name: "aws example with the original ghost busters",
			input: struct {
				Country      string
				Region       string
				State        string
				County       string
				City         string
				Neighborhood string
			}{
				Country:      "USA",
				Region:       "East",
				State:        "NY",
				County:       "Queens",
				City:         "New York",
				Neighborhood: "Something strange",
			},
			want: &types.AttributeValueMemberS{
				Value: "foo#biz#baz#",
			},
			withTag: TagDefaultName,
		},
		{
			name: "floats",
			input: struct {
				F32 float32
				F64 float64
			}{
				F32: 0.123,
				F64: 0.4621,
			},
			withTag: TagDefaultName,
			want:    &types.AttributeValueMemberS{Value: "0.123#0.4621"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalFields(tt.input, tt.withTag)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalFields() got = %q, want = %q", got.Value, tt.want.Value)
				return
			}

			back := reflect.New(reflect.TypeOf(tt.input)).Interface()
			err = UnmarshalFields(back, TagDefaultName, got)
			if err != nil {
				t.Error(err)
				return
			}

			bv := reflect.ValueOf(back).Elem().Interface()

			if !reflect.DeepEqual(bv, tt.input) {
				t.Errorf("UnmarshalFields() got = %+v, original %+v", bv, tt.input)
			}
		})
	}
}
