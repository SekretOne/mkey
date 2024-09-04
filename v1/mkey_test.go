package mkey

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"testing"
)

func testSymetry[T any](t *testing.T, name, tag string, input T, want *types.AttributeValueMemberS, wantErr bool) {
	t.Run(name, func(t *testing.T) {
		got, err := MarshalFields(input, tag)
		if (err != nil) != wantErr {
			t.Errorf("MarshalFields() error = %v, wantErr %v", err, wantErr)
			return
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("MarshalFields() got = %v, want = %v", got, want)
		}

		var back T

		err = UnmarshalFields(&back, tag, got)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(back, input) {
			t.Errorf("UnmarshalFields() got = %v, want %v", back, input)
		}
	})
}

func TestMarshalMultiFieldKeyWithTag(t *testing.T) {
	type args struct {
		input any
		tag   string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.AttributeValueMemberS
		wantErr bool
	}{
		{
			name: "The simplest example",
			args: args{
				input: struct {
					A string
				}{
					A: "value",
				},
				tag: "",
			},
			want: &types.AttributeValueMemberS{
				Value: "A=value",
			},
			wantErr: false,
		},
		{
			name: "Non exported fields are skipped",
			args: args{
				input: struct {
					a string
					B string
					c string
					D string
				}{
					a: "value-1",
					B: "value-2",
					c: "value-3",
					D: "value-4",
				},
				tag: "",
			},
			want: &types.AttributeValueMemberS{
				Value: "B=value-2|D=value-4",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalFields(tt.args.input, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalFields() got = %v, want = %v", got, tt.want)
			}

			//back := reflect.New(reflect.TypeOf(tt.args.input)).Elem().Interface()
			// back := reflect.New(reflect.TypeOf(tt.args.input)).Elem().Interface()
			back := tt.args.input
			bv := reflect.ValueOf(&back).Elem()
			bv.Set(reflect.Zero(bv.Type()))

			err = UnmarshalFields(&back, "", got)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(back, tt.args.input) {
				t.Errorf("UnmarshalFields() got = %v, original %v", back, tt.args.input)
			}
		})
	}
}
