package gir

import (
	"encoding/xml"
	"fmt"

	"github.com/dave/jennifer/jen"
)

type Record struct {
	XMLName xml.Name `xml:"http://www.gtk.org/introspection/core/1.0 record"`
	Name    string   `xml:"name,attr"`

	CType         string `xml:"http://www.gtk.org/introspection/c/1.0 type,attr"`
	CSymbolPrefix string `xml:"http://www.gtk.org/introspection/c/1.0 symbol-prefix,attr"`

	GLibTypeName string `xml:"http://www.gtk.org/introspection/glib/1.0 type-name,attr"`
	GLibGetType  string `xml:"http://www.gtk.org/introspection/glib/1.0 get-type,attr"`

	Methods   []Method   `xml:"http://www.gtk.org/introspection/core/1.0 method"`
	Functions []Function `xml:"http://www.gtk.org/introspection/core/1.0 function"`
}

func (r Record) IsIgnored() bool {
	// unsure why some records don't have GetType
	return r.GLibGetType == ""
}

func (r Record) GenerateAll() *jen.Statement {
	f := new(jen.Statement)
	f.Add(r.GenType())
	f.Line()
	f.Add(r.GenMarshaler())
	f.Line()
	f.Add(r.GenNative())
	f.Line()
	f.Add(r.GenMethods())
	return f
}

func (r Record) GenType() *jen.Statement {
	return jen.Type().Id(r.Name).Qual("C", r.CType)
}

func (r Record) GenMarshalerItem() *jen.Statement {
	return GenMarshalerItem(r.GLibGetType, r.GoName())
}

func (r Record) GenMarshaler() *jen.Statement {
	var goName = r.GoName()

	return GenMarshalerFn(
		goName,
		jen.Return(
			jen.Parens(jen.Op("*").Qual("C", r.CType)).Call(
				jen.Qual("unsafe", "Pointer").Call(
					jen.Qual("C", "g_value_get_boxed").Call(
						jen.Parens(jen.Op("*").Qual("C", "GValue")).Call(
							jen.Qual("unsafe", "Pointer").Call(jen.Id("p")),
						),
					),
				),
			),
			jen.Nil(),
		),
	)
}

func (r Record) GenNative() *jen.Statement {
	i := firstChar(r.Name)
	p := jen.Id(i).Op("*").Id(r.GoName())

	f := jen.Add(GenCommentReflowLines("native", fmt.Sprintf(
		"turns the current *%s into the native C pointer type.",
		r.GoName(),
	)))

	f.Func().Params(p).Id("native").Params().Id("*" + r.CGoType())
	f.Block(
		jen.Return(
			jen.Parens(jen.Op("*").Qual("C", r.CType)).Call(
				jen.Qual("unsafe", "Pointer").Call(jen.Id(i)),
			),
		),
	)

	f.Line()
	f.Line()

	f.Func().Params(p).Id("Native").Params().Uintptr()
	f.Block(
		jen.Return(
			jen.Uintptr().Call(
				jen.Qual("unsafe", "Pointer").Call(jen.Id(i).Dot("native").Call()),
			),
		),
	)

	f.Line()
	return f
}

func (r Record) GenMethods() *jen.Statement {
	var stmt = make(jen.Statement, 0, len(r.Methods)*3)
	for _, method := range r.Methods {
		if method.IsIgnored() {
			continue
		}

		stmt.Add(method.GenFunc(r.Name))
		stmt.Line()
	}

	return &stmt
}

func (r Record) CGoType() string {
	return CGoType(r.CType)
}

func (r Record) GoName() string {
	return snakeToGo(true, r.Name)
}
