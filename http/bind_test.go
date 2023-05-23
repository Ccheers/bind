package http

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type (
	TestBind struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	TestBind2 struct {
		Age int `json:"age"`
	}
)

func TestBindQuery(t *testing.T) {
	type args struct {
		vars   url.Values
		target interface{}
	}

	tests := []struct {
		name string
		args args
		err  error
		want interface{}
	}{
		{
			name: "test",
			args: args{
				vars:   map[string][]string{"name": {"cc-bind"}, "url": {"https://go-cc-bind.dev/"}},
				target: &TestBind{},
			},
			err:  nil,
			want: &TestBind{"cc-bind", "https://go-cc-bind.dev/"},
		},
		{
			name: "test1",
			args: args{
				vars:   map[string][]string{"age": {"cc-bind"}, "url": {"https://go-cc-bind.dev/"}},
				target: &TestBind2{},
			},
			err: ErrBind,
		},
		{
			name: "test2",
			args: args{
				vars:   map[string][]string{"age": {"1"}, "url": {"https://go-cc-bind.dev/"}},
				target: &TestBind2{},
			},
			err:  nil,
			want: &TestBind2{Age: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BindURLValues(tt.args.vars, tt.args.target)
			if !errors.Is(err, tt.err) {
				t.Fatalf("BindURLValues() error = %v, err %v", err, tt.err)
			}
			if err == nil && !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("BindURLValues() target = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}

func TestBindForm(t *testing.T) {
	type args struct {
		req    *http.Request
		target interface{}
	}

	tests := []struct {
		name string
		args args
		err  error
		want *TestBind
	}{
		{
			name: "error not nil",
			args: args{
				req:    &http.Request{Method: http.MethodPost},
				target: &TestBind{},
			},
			err:  errors.New("missing form body"),
			want: nil,
		},
		{
			name: "error is nil",
			args: args{
				req: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{"Content-Type": {"application/x-www-form-urlencoded; param=value"}},
					Body:   io.NopCloser(strings.NewReader("name=cc-bind&url=https://go-cc-bind.dev/")),
				},
				target: &TestBind{},
			},
			err:  nil,
			want: &TestBind{"cc-bind", "https://go-cc-bind.dev/"},
		},
		{
			name: "error BadRequest",
			args: args{
				req: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{"Content-Type": {"application/x-www-form-urlencoded; param=value"}},
					Body:   io.NopCloser(strings.NewReader("age=a")),
				},
				target: &TestBind2{},
			},
			err:  ErrBind,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BindForm(tt.args.req, tt.args.target)
			if err == nil && !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("BindForm() target = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}
