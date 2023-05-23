package http

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestBindRequestQuery(t *testing.T) {

	req, err := http.NewRequest(
		http.MethodPost, "http://localhost:8080?age=1",
		strings.NewReader("<xml></xml>"))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/xml")
	type dst struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var res dst
	if err := BindRequestQuery(req, &res); err != nil {
		t.Errorf("BindRequestQuery() error = %v", err)
	}
	if err := BindRequestBody(req, &res); err != nil {
		t.Errorf("BindRequestBody() error = %v", err)
	}

	log.Println(res)
}

func mustNewReq(method, contentType, url string, reader io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", contentType)
	return req
}

func TestTryMyBestBind1(t *testing.T) {
	type dst struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var res dst

	formBody := &bytes.Buffer{}
	fw := multipart.NewWriter(formBody)
	fw.WriteField("name", "haijun")
	fw.WriteField("age", "12")
	fw.Close()

	formBody2 := &bytes.Buffer{}
	fw = multipart.NewWriter(formBody2)
	wr, _ := fw.CreateFormFile("file", ".gitignore")
	rd, _ := os.Open("/Users/eric/GoProject/bind/.gitignore")
	io.Copy(wr, rd)
	rd.Close()
	fw.WriteField("name", "haijun")
	fw.WriteField("age", "12")
	fw.Close()

	type args struct {
		r    *http.Request
		v    interface{}
		opts []OptionFunc
	}
	tests := []struct {
		name string
		args args
		want *dst
	}{
		{
			name: "1",
			args: args{
				r: mustNewReq(http.MethodPost, "application/json", "http://localhost:8080?age=1", strings.NewReader("{\"name\":\"haijun\",\"age\":2}")),
				v: &res,
			},
			want: &dst{
				Name: "haijun",
				Age:  2,
			},
		},
		{
			name: "2",
			args: args{
				r: mustNewReq(http.MethodPost, fw.FormDataContentType(), "http://localhost:8080?age=1", formBody),
				v: &res,
			},
			want: &dst{
				Name: "haijun",
				Age:  1,
			},
		},
		{
			name: "3",
			args: args{
				r: mustNewReq(http.MethodPost, fw.FormDataContentType(), "http://localhost:8080", formBody2),
				v: &res,
			},
			want: &dst{
				Name: "haijun",
				Age:  12,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TryMyBestBind(tt.args.r, tt.args.v, tt.args.opts...)
			if !reflect.DeepEqual(tt.args.v, tt.want) {
				t.Errorf("got=%+v, want=%+v", tt.args.v, tt.want)
			}
		})
	}
}
