package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type listItem struct {
	Name   string
	Suffix string
}

var list = []listItem{
	{Name: "Buffersink"},
	{Name: "Buffersrc"},
	{Name: "CodecContext"},
	{Name: "CodecContext", Suffix: "2"},
	{Name: "Dictionary"},
	{Name: "FilterCommand"},
	{Name: "FormatContextCtx"},
	{Name: "FormatContext"},
	{Name: "FormatEvent"},
	{Name: "IOContext"},
	{Name: "IOFormat"},
	{Name: "Packet"},
	{Name: "Seek"},
	{Name: "StreamEvent"},
}

var tmpl = `// Code generated by astiav. DO NOT EDIT.
package astiav
{{ range $val := . }}
type {{ $val.Name }}Flags{{ $val.Suffix }} flags

func New{{ $val.Name }}Flags{{ $val.Suffix }}(fs ...{{ $val.Name }}Flag{{ $val.Suffix }}) {{ $val.Name }}Flags{{ $val.Suffix }} {
	o := {{ $val.Name }}Flags{{ $val.Suffix }}(0)
	for _, f := range fs {
		o = o.Add(f)
	}
	return o
}

func (fs {{ $val.Name }}Flags{{ $val.Suffix }}) Add(f {{ $val.Name }}Flag{{ $val.Suffix }}) {{ $val.Name }}Flags{{ $val.Suffix }} {
	return {{ $val.Name }}Flags{{ $val.Suffix }}(flags(fs).add(int(f)))
}

func (fs {{ $val.Name }}Flags{{ $val.Suffix }}) Del(f {{ $val.Name }}Flag{{ $val.Suffix }}) {{ $val.Name }}Flags{{ $val.Suffix }} {
	return {{ $val.Name }}Flags{{ $val.Suffix }}(flags(fs).del(int(f)))
}

func (fs {{ $val.Name }}Flags{{ $val.Suffix }}) Has(f {{ $val.Name }}Flag{{ $val.Suffix }}) bool { return flags(fs).has(int(f)) }
{{ end }}`

var tmplTest = `// Code generated by astiav. DO NOT EDIT.
package astiav_test
import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)
{{ range $val := . }}
func Test{{ $val.Name }}Flags{{ $val.Suffix }}(t *testing.T) {
	fs := astiav.New{{ $val.Name }}Flags{{ $val.Suffix }}(astiav.{{ $val.Name }}Flag{{ $val.Suffix }}(1))
	require.True(t, fs.Has(astiav.{{ $val.Name }}Flag{{ $val.Suffix }}(1)))
	fs = fs.Add(astiav.{{ $val.Name }}Flag{{ $val.Suffix }}(2))
	require.True(t, fs.Has(astiav.{{ $val.Name }}Flag{{ $val.Suffix }}(2)))
	fs = fs.Del(astiav.{{ $val.Name }}Flag{{ $val.Suffix }}(2))
	require.False(t, fs.Has(astiav.{{ $val.Name }}Flag{{ $val.Suffix }}(2)))
}
{{ end }}`

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(fmt.Errorf("main: getting working directory failed: %w", err))
	}

	f, err := os.Create(filepath.Join(dir, "flags.go"))
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating file failed: %w", err))
	}
	defer f.Close()

	if err = template.Must(template.New("tmpl").Parse(tmpl)).Execute(f, list); err != nil {
		log.Fatal(fmt.Errorf("main: executing template failed: %w", err))
	}

	ft, err := os.Create(filepath.Join(dir, "flags_test.go"))
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating test file failed: %w", err))
	}
	defer ft.Close()

	if err = template.Must(template.New("tmpl").Parse(tmplTest)).Execute(ft, list); err != nil {
		log.Fatal(fmt.Errorf("main: executing template failed: %w", err))
	}
}
