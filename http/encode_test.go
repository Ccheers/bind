package http

import (
	"testing"

	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/Ccheers/bind/internal/testdata/binding"
)

func TestEncodeURL(t *testing.T) {
	tests := []struct {
		pathTemplate string
		request      *binding.HelloRequest
		needQuery    bool
		want         string
	}{
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "2233!!!!"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/test/sub/2233!!!!",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
			request:      nil,
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{}/sub/{sub.naming}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "hello"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/{}/sub/hello",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{}/sub/{sub.naming}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "hello"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/{}/sub/hello",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{}/sub/{sub.name.cc}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "hello"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/{}/sub/",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{}/sub/{test_repeated}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "hello"}, TestRepeated: []string{"123", "456"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/{}/sub/123",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "5566!!!"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/test/sub/5566!!!",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/sub",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "2233!!!"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/sub",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.name}",
			request:      &binding.HelloRequest{Name: "test"},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/test/sub/",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.name}",
			request:      &binding.HelloRequest{Name: "test"},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/test/sub/",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub",
			request:      &binding.HelloRequest{Name: "go", Sub: &binding.Sub{Name: "cc-bind"}},
			needQuery:    true,
			want:         "http://helloworld.Greeter/helloworld/go/sub?sub.naming=cc-bind",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/sub/{sub.naming}",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "cc-bind"}, UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"name", "sub.naming"}}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/sub/cc-bind?updateMask=name,sub.naming",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/sub/[{sub.naming}]",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "cc-bind"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/sub/[cc-bind]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/[{name}]/sub/[{sub.naming}]",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "cc-bind"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/[test]/sub/[cc-bind]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/[{}]/sub/[{sub.naming}]",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "cc-bind"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/[{}]/sub/[cc-bind]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/[{}]/sub/[{sub.naming}]/{[]}",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "cc-bind"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/[{}]/sub/[cc-bind]/{[]}",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{[sub]}/[{sub.naming}]",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "cc-bind"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/{[sub]}/[cc-bind]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{[name]}/[{sub.naming}]",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "cc-bind"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/{[name]}/[cc-bind]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{}/[]/[{sub.naming}]",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "cc-bind"}},
			needQuery:    false,
			want:         "http://helloworld.Greeter/helloworld/{}/[]/[cc-bind]",
		},
	}

	for _, test := range tests {
		if EncodeURL(test.pathTemplate, test.request, test.needQuery) != test.want {
			t.Fatalf("want: %s, got: %s", test.want, EncodeURL(test.pathTemplate, test.request, test.needQuery))
		}
	}
}