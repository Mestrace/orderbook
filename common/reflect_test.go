package common

import (
	"reflect"
	"testing"
)

func TestUnit_SetFieldByName(t *testing.T) {
	type testType struct {
		Str   string
		Int64 int64
		Int32 int32
		Int   int
	}
	type args struct {
		S     interface{}
		field string
		value string
	}

	tests := []struct {
		name      string
		args      args
		want      bool
		wantErr   bool
		wantValue interface{}
	}{
		{
			name: "valid string",
			args: args{
				S:     &testType{Str: "unchanged"},
				field: "Str",
				value: "changed",
			},
			want:      true,
			wantErr:   false,
			wantValue: &testType{Str: "changed"},
		},
		{
			name: "valid int",
			args: args{
				S:     &testType{Int32: -1},
				field: "Int32",
				value: "10",
			},
			want:      true,
			wantErr:   false,
			wantValue: &testType{Int32: 10},
		},
		{
			name: "unknown field",
			args: args{
				S:     &testType{},
				field: "Unknown",
				value: "unknown",
			},
			want:      false,
			wantErr:   false,
			wantValue: &testType{},
		},
		{
			name: "int not a number",
			args: args{
				S:     &testType{Int32: -1},
				field: "Int32",
				value: "clearly_not_a_number",
			},
			want:      false,
			wantErr:   true,
			wantValue: &testType{Int32: -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetFieldByName(tt.args.S, tt.args.field, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetFieldByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetFieldByName() = %v, want %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(tt.args.S, tt.wantValue) {
				t.Errorf("SetFieldByName() sets %v, want %v", tt.args.S, tt.wantValue)
			}
		})
	}
}
