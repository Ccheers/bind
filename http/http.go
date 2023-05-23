package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Ccheers/bind/internal/encoding"
	"github.com/Ccheers/bind/internal/encoding/json"

	_ "github.com/Ccheers/bind/internal/encoding/form"
	_ "github.com/Ccheers/bind/internal/encoding/json"
	_ "github.com/Ccheers/bind/internal/encoding/proto"
	_ "github.com/Ccheers/bind/internal/encoding/xml"
	_ "github.com/Ccheers/bind/internal/encoding/yaml"
)

const (
	baseContentType = "application"
)

// BindRequestVars decodes the request vars to object.
func BindRequestVars(r *http.Request, raws map[string]string, v interface{}) error {
	vars := make(url.Values, len(raws))
	for k, v := range raws {
		vars[k] = []string{v}
	}
	return BindURLValues(vars, v)
}

// BindRequestQuery decodes the request vars to object.
func BindRequestQuery(r *http.Request, v interface{}) error {
	return BindURLValues(r.URL.Query(), v)
}

type OptionFunc func(*Options)

type Options struct {
	defaultCodec string
}

// WithDefaultCodec set default codec.
func WithDefaultCodec(codec string) OptionFunc {
	return func(o *Options) {
		o.defaultCodec = codec
	}
}

func defaultOption() Options {
	return Options{
		defaultCodec: json.Name,
	}
}

// BindRequestBody decodes the request body to object.
func BindRequestBody(r *http.Request, v interface{}, opts ...OptionFunc) error {
	options := defaultOption()
	for _, opt := range opts {
		opt(&options)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrBind, err.Error())
	}

	if len(data) == 0 {
		return nil
	}

	codec, _ := CodecForRequest(r, "Content-Type", options)
	if err = codec.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%w: body unmarshal %s", ErrBind, err.Error())
	}
	return nil
}

// CodecForRequest get encoding.Codec via http.Request
func CodecForRequest(r *http.Request, name string, opt Options) (encoding.Codec, bool) {
	for _, accept := range r.Header[name] {
		codec := encoding.GetCodec(ContentSubtype(accept))
		if codec != nil {
			return codec, true
		}
	}
	return encoding.GetCodec(opt.defaultCodec), false
}

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}

// ContentSubtype returns the content-subtype for the given content-type.  The
// given content-type must be a valid content-type that starts with
// but no content-subtype will be returned.
// according rfc7231.
// contentType is assumed to be lowercase already.
func ContentSubtype(contentType string) string {
	left := strings.Index(contentType, "/")
	if left == -1 {
		return ""
	}
	right := strings.Index(contentType, ";")
	if right == -1 {
		right = len(contentType)
	}
	if right < left {
		return ""
	}
	return contentType[left+1 : right]
}
