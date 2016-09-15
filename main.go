package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"gopkg.in/alecthomas/kingpin.v2"
)

var reactTemplate = "rx.got"
var tplFile = filepath.Join(os.Getenv("GOPATH"), "src", "rxgo", reactTemplate)

var (
	commonTypes = strings.Fields("int string float32 float64")
	baseTypes   = strings.Fields("bool rune byte string uint int uint8 int8 uint16 int16 uint32 int32 uint64 int64 float32 float64 complex64 complex128 time.Time time.Duration []byte")
	fullTypes   []string
	mapTypes    []string
)

func containsFullType(name string) bool {
	for _, n := range fullTypes {
		if n == name {
			return true
		}
	}
	return false
}

func appendFullTypes(names ...string) {
	for _, name := range names {
		if !containsFullType(name) {
			fullTypes = append(fullTypes, name)
		}
	}
}

var (
	packageArg      = kingpin.Arg("package", "Go package.").Required().String()
	typesArg        = kingpin.Arg("types", "List of types to provide Reactive types for.").Strings()
	commonTypesFlag = kingpin.Flag("common-types", "Define extensions for common Go types (bool rune byte string uint int uint8 int8 uint16 int16 uint32 int32 uint64 int64 float32 float64 complex64 complex128 time.Time time.Duration []byte).").Action(func(*kingpin.ParseContext) error {
		appendFullTypes(commonTypes...)
		return nil
	}).Bool()
	baseTypesFlag = kingpin.Flag("base-types", "Define extensions for all base Go types (int string float32 float64).").Action(func(*kingpin.ParseContext) error {
		appendFullTypes(baseTypes...)
		return nil
	}).Bool()
	importsFlag   = kingpin.Flag("import", "Extra imports.").PlaceHolder("PKG...").Strings()
	outputFlag    = kingpin.Flag("output", "File to write to.").Short('o').String()
	debugFlag     = kingpin.Flag("debug", "Debug mode.").Bool()
	maxReplayFlag = kingpin.Flag("max-replay", "Maximum size of replayed data.").Default("16384").Int()
	mapTypesFlag  = kingpin.Flag("map-types", "Limit Map and FlatMap target types to this list").PlaceHolder("<type>,...").String()
)

type Context struct {
	Package       string
	Generate      []string
	Types         []string
	MapTypes      []string
	Imports       []string
	MaxReplaySize int
}

func typeName(t ast.Expr) []string {
	switch n := t.(type) {
	case *ast.StarExpr:
		return typeName(n.X)
	case *ast.SelectorExpr:
		// return append(typeName(n.X), typeName(n.Sel)...)
		return typeName(n.Sel)
	case *ast.MapType:
		keys := typeName(n.Key)
		keys = append(keys, typeName(n.Value)...)
		keys = append(keys, "Map")
		return keys
	case *ast.ArrayType:
		return append(typeName(n.Elt), "Slice")
	case *ast.Ident:
		return []string{strings.Title(n.Name)}
	default:
		panic(fmt.Sprintf("unknown expression node %s %s\n", reflect.TypeOf(t), t))
	}
}

func main() {
	kingpin.Parse()

	t := &template.Template{}
	t = template.Must(t.Funcs(template.FuncMap{
		"TypeName": func(s string) string {
			st, err := parser.ParseExpr(s)
			kingpin.FatalIfError(err, "invalid type %q", s)
			return strings.Join(typeName(st), "")
		},
		"IsNumeric": func(v string) bool {
			switch v {
			case "byte", "uint", "int", "uint8", "int8", "uint16", "int16", "uint32",
				"int32", "uint64", "int64", "float32", "float64":
				return true
			}
			return false
		},
		"Join": func(a []string) string { return strings.Join(a, " ") },
	}).ParseFiles(tplFile))

	cmd := exec.Command("goimports")
	if *debugFlag {
		cmd = exec.Command("cat")
	}

	outputFile := os.Stdout
	if *outputFlag != "" {
		dir := filepath.Dir(*outputFlag)
		var err error
		outputFile, err = ioutil.TempFile(dir, ".rxgo-")
		kingpin.FatalIfError(err, "")
		defer os.Remove(outputFile.Name())
	}
	cmd.Stdout = outputFile

	// Add typesArg to fullTypes
	appendFullTypes(*typesArg...)

	// Allow explicit map and flatmap target types.
	if len(*mapTypesFlag) > 0 {
		for _, v := range strings.Split(*mapTypesFlag, ",") {
			mapTypes = append(mapTypes, v)
			// A map target type is also a full type
			appendFullTypes(v)
		}
	} else {
		mapTypes = fullTypes
	}

	stdin, err := cmd.StdinPipe()
	kingpin.FatalIfError(err, "")
	kingpin.FatalIfError(cmd.Start(), "")
	err = t.ExecuteTemplate(stdin, reactTemplate, &Context{
		Generate:      append([]string{"$GOPATH/bin/rxgo"}, os.Args[1:]...),
		Package:       *packageArg,
		Types:         fullTypes,
		MapTypes:      mapTypes,
		Imports:       *importsFlag,
		MaxReplaySize: *maxReplayFlag,
	})
	kingpin.FatalIfError(err, "")
	stdin.Close()
	kingpin.FatalIfError(cmd.Wait(), "")
	outputFile.Close()
	if *outputFlag != "" {
		kingpin.FatalIfError(os.Rename(outputFile.Name(), *outputFlag), "")
	}
}
