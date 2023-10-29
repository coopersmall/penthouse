package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

type object struct {
	Package string
	Imports []string
	Name    string
	Methods []method
}

type method struct {
	Name    string
	Params  []param
	Returns []ret
}

type param struct {
	Name string
	Type string
}

type ret struct {
	Type string
}

func generateMockObjects(filename string, generator *template.Template) {
	// Parse the source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing source file: %v\n", err)
		return
	}

	// get the package name
	packageName := node.Name.Name

	// get imports
	imports := make([]string, 0)
	for _, imp := range node.Imports {
		imports = append(imports, imp.Path.Value)
	}

	// Create a map to store interfaces and structs
	objects := make([]object, 0)

	// Inspect type declarations and identify interfaces and structs
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, isInterface := typeSpec.Type.(*ast.InterfaceType); isInterface {
						methods := make([]method, 0)
						for _, field := range typeSpec.Type.(*ast.InterfaceType).Methods.List {
							var (
								methodName = field.Names[0].Name
								params     = make([]param, 0)
								returns    = make([]ret, 0)
							)

							for _, p := range field.Type.(*ast.FuncType).Params.List {
								for _, paramName := range p.Names {
									params = append(params, param{
										Name: paramName.Name,
										Type: p.Type.(*ast.Ident).Name,
									})
								}
							}

							for _, r := range field.Type.(*ast.FuncType).Results.List {
								if r == nil {
									continue
								}

								returns = append(returns, ret{
									Type: r.Type.(*ast.Ident).Name,
								})
							}

							methods = append(methods, method{
								Name:    methodName,
								Params:  params,
								Returns: returns,
							})
						}

						obj := object{
							Package: packageName,
							Imports: imports,
							Name:    typeSpec.Name.Name,
							Methods: methods,
						}

						objects = append(objects, obj)

					}
				}
			}
		}
	}

	for _, obj := range objects {
		var s strings.Builder
		if err := generator.Execute(&s, obj); err != nil {
			fmt.Printf("Error generating mock code for %s: %v\n", obj.Name, err)
			continue
		}

		mockFileName := strings.ToLower(obj.Name) + "_mock.go"
		if err := os.WriteFile(mockFileName, []byte(s.String()), 0644); err != nil {
			fmt.Printf("Error writing mock file for %s: %v\n", obj.Name, err)
		} else {
			fmt.Printf("Generated mock file: %s\n", mockFileName)
		}
	}
}

const mockTemplate = `package {{ .Package }}

import (
    "github.com/coopersmall/penthouse/mock"
    {{- range .Imports }}{{ . }}{{- end }}
)

type {{ .Name }}Mock struct {
    Mock mock.Mock
}

func New{{ .Name }}Mock() *{{ .Name }}Mock {
    return &{{ .Name }}Mock{
        Mock: mock.NewMock(),
    }
}

{{ range .Methods }}
func (m *{{ $.Name }}Mock) {{ .Name }}({{ $numParams := len .Params}}{{ range $idx, $param := .Params }}{{ if and (ne $idx $numParams) (ne $idx 0) }}, {{end}}{{ .Name }} {{ .Type }}{{ end }}) {{ if ne 0 (len .Returns) }}({{$numRets := len .Returns}}{{ range $idx, $ret := .Returns }}{{if and (ne $idx $numRets) (ne $idx 0)}}, {{end}}{{ .Type }}{{ end }}){{ end }} {
    args := m.Mock.CallMethod("{{.Name}}", {{ range $idx, $params := .Params }}{{if and (ne $idx $numParams) (ne $idx 0)}}, {{end}}{{ $params.Name }}{{ end }})
    return {{ $numReturns := len .Returns }}{{ range $idx, $ret := .Returns }}{{ if and (ne $idx $numReturns) (ne $idx 0 ) }}, {{end}}{{if eq "error" $ret.Type}}mock.Error(args[{{ $idx }}]){{ else }}args[{{ $idx }}].({{ $ret.Type }}){{ end }}{{ end }}
}
{{ end }}
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_mock <filename.go>")
		os.Exit(1)
	}

	generator := template.Must(template.New("mock").Parse(mockTemplate))

	filename := os.Args[1]
	generateMockObjects(filename, generator)
}
