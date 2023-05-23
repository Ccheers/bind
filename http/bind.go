package http

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"

	"github.com/Ccheers/bind/internal/encoding"
	"github.com/Ccheers/bind/internal/encoding/form"
)

var (
	ErrBind = fmt.Errorf("bind error")
)

const (
	MegaByte = 1024 * 1024
)

// TryMyBestBind decodes the request to object.
// 不支持解析 multipart/form-data ， 需要解析大文件数据请使用 BindForm
// 只有 MegaByte 以下的内容才能被 form-data 解析
func TryMyBestBind(r *http.Request, v interface{}, opts ...OptionFunc) {
	_ = BindRequestQuery(r, v)

	// 大参数的数据不能支持
	if r.ContentLength > MegaByte {
		return
	}
	// reset body.
	data, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	_ = BindForm(r, v)
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	_ = BindRequestBody(r, v)
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	return
}

// BindURLValues bind vars parameters to target.
func BindURLValues(vars url.Values, target interface{}) error {
	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(vars.Encode()), target); err != nil {
		return fmt.Errorf("%w: %v, target=%T, module=BindURLValues", ErrBind, err, target)
	}
	return nil
}

// BindForm bind form parameters to target.
func BindForm(req *http.Request, target interface{}) error {
	if err := req.ParseMultipartForm(math.MaxInt32); err != nil {
		return err
	}
	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(req.Form.Encode()), target); err != nil {
		return fmt.Errorf("%w: %v, target=%T, module=BindURLValues", ErrBind, err, target)
	}
	return nil
}
