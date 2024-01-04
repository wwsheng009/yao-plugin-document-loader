package main

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/yaoapp/kun/grpc"
)

func TestDocumentLoader_Exec(t *testing.T) {
	type fields struct {
		Plugin grpc.Plugin
	}
	type args struct {
		method string
		args   []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *grpc.Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				method: "text",
				args:   []interface{}{"./yaoapp/data/test.ziw"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &DocumentLoader{
				Plugin: tt.fields.Plugin,
			}
			var output io.Writer = os.Stdout
			doc.SetLogger(output, grpc.Trace)
			got, err := doc.Exec(tt.args.method, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DocumentLoader.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DocumentLoader.Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}
