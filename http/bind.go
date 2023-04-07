package http

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Ccheers/bind/internal/encoding"
	"github.com/Ccheers/bind/internal/encoding/form"
)

var (
	ErrBind = fmt.Errorf("bind error")
)

// BindQuery bind vars parameters to target.
func BindQuery(vars url.Values, target interface{}) error {
	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(vars.Encode()), target); err != nil {
		return fmt.Errorf("%w: %v, target=%T, module=BindQuery", ErrBind, err, target)
	}
	return nil
}

// BindForm bind form parameters to target.
func BindForm(req *http.Request, target interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(req.Form.Encode()), target); err != nil {
		return fmt.Errorf("%w: %v, target=%T, module=BindQuery", ErrBind, err, target)
	}
	return nil
}
