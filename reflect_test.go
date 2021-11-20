package digpro

import (
	"reflect"
	"testing"
)

func Test_makeParameterObjectType_copyFromParameterObject(t *testing.T) {

	var fooNilPtr *Foo

	type args struct {
		structOrStructPtr interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error nil",
			args: args{
				structOrStructPtr: nil,
			},
			wantErr: true,
		},
		{
			name: "error nil struct",
			args: args{
				structOrStructPtr: fooNilPtr,
			},
			wantErr: true,
		},
		{
			name: "error string",
			args: args{
				structOrStructPtr: "string",
			},
			wantErr: true,
		},
		{
			name: "error int ptr",
			args: args{
				structOrStructPtr: new(int),
			},
			wantErr: true,
		},
		{
			name: "error field conflict",
			args: args{
				structOrStructPtr: struct {
					A string
					a string
				}{},
			},
			wantErr: true,
		},
		{
			name: "success struct ptr",
			args: args{
				structOrStructPtr: new(Foo),
			},
			wantErr: false,
		},
		{
			name: "success struct value",
			args: args{
				structOrStructPtr: Foo{},
			},
			wantErr: false,
		},
		{
			name: "success ignore tag",
			args: args{
				structOrStructPtr: struct {
					A int
					a int `digpro:"ignore"`
				}{
					A: 0,
					a: int(mockDataMapping[reflect.Int].Int()), // not conflict, and equal int value
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parameterObjectType, fieldMapping, err := makeParameterObjectType(tt.args.structOrStructPtr, false)
			if (err != nil) != tt.wantErr {
				// not want err, but got error , return
				if !tt.wantErr {
					t.Errorf("makeParameterObjectType() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// want err, but got not error, continue...
			}
			if err != nil {
				return
			}
			injectedValue, err := copyFromParameterObject(tt.args.structOrStructPtr, mockParameterObject(parameterObjectType), fieldMapping)
			if (err != nil) != tt.wantErr {
				t.Errorf("copyFromParameterObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			assertStructOrStructPtr(t, injectedValue)
			if reflect.TypeOf(tt.args.structOrStructPtr) != reflect.TypeOf(injectedValue) {
				t.Errorf("reflect.TypeOf(injectedValue) got = %s, want %s", reflect.TypeOf(injectedValue), reflect.TypeOf(tt.args.structOrStructPtr))
				return
			}
		})
	}
}

func Test_copyFromParameterObject(t *testing.T) {
	type args struct {
		structOrStructPtr    interface{}
		parameterObjectValue reflect.Value
		fieldMapping         map[string]int
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "error nil",
			args: args{
				structOrStructPtr: nil,
			},
			wantErr: true,
		},
		{
			name: "error nil struct",
			args: args{
				structOrStructPtr: fooNilPtr,
			},
			wantErr: true,
		},
		{
			name: "error string",
			args: args{
				structOrStructPtr: "string",
			},
			wantErr: true,
		},
		{
			name: "error int ptr",
			args: args{
				structOrStructPtr: new(int),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := copyFromParameterObject(tt.args.structOrStructPtr, tt.args.parameterObjectValue, tt.args.fieldMapping)
			if (err != nil) != tt.wantErr {
				t.Errorf("copyFromParameterObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("copyFromParameterObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
