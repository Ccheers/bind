package http

import (
	"log"
	"net/http"
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
