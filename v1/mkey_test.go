package mkey

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"testing"
)

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
			got, err := MarshalMultiFieldKeyWithTag(tt.args.input, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalMultiFieldKeyWithTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalMultiFieldKeyWithTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}
